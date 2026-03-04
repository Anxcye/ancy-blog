// File: auth.go
// Purpose: Handle authentication, token issuance, and token refresh workflows.
// Module: backend/internal/service, authentication service layer.
// Related: auth middleware, auth handlers, and credential store.
package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/anxcye/ancy-blog/backend/internal/domain"
	"github.com/anxcye/ancy-blog/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type authSession struct {
	User         domain.User
	RefreshToken string
	ExpiresAt    time.Time
}

type AuthService struct {
	mu sync.RWMutex

	adminUsername string
	adminPassword string // plain-text fallback (in-memory / test mode)

	credentials repository.CredentialStore

	accessTTL  time.Duration
	refreshTTL time.Duration

	accessTokens  map[string]authSession
	refreshTokens map[string]authSession
}

type AuthResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

func NewAuthService(adminUsername, adminPassword string, accessTTL, refreshTTL time.Duration) *AuthService {
	if adminUsername == "" {
		adminUsername = "admin"
	}
	if adminPassword == "" {
		adminPassword = "123456"
	}
	if accessTTL <= 0 {
		accessTTL = 1 * time.Hour
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &AuthService{
		adminUsername: adminUsername,
		adminPassword: adminPassword,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		accessTokens:  make(map[string]authSession),
		refreshTokens: make(map[string]authSession),
	}
}

// WithCredentialStore wires a persistent credential backend. Call after construction
// to enable bcrypt password verification and persistence across restarts.
// Bootstraps the stored hash from the plain-text env password if none exists yet.
func (s *AuthService) WithCredentialStore(store repository.CredentialStore) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.credentials = store
	if _, ok := store.GetAdminPasswordHash(); !ok {
		if hash, err := bcrypt.GenerateFromPassword([]byte(s.adminPassword), bcrypt.DefaultCost); err == nil {
			_ = store.SetAdminPasswordHash(string(hash))
		}
	}
}

func (s *AuthService) Login(username, password string) (AuthResult, error) {
	if username != s.adminUsername {
		return AuthResult{}, errors.New("invalid credentials")
	}
	s.mu.RLock()
	store := s.credentials
	s.mu.RUnlock()

	if store != nil {
		hash, ok := store.GetAdminPasswordHash()
		if !ok || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
			return AuthResult{}, errors.New("invalid credentials")
		}
	} else {
		if password != s.adminPassword {
			return AuthResult{}, errors.New("invalid credentials")
		}
	}
	user := domain.User{ID: "admin-1", Username: s.adminUsername, DisplayName: "Administrator", IsAdmin: true}
	return s.issueTokens(user)
}

func (s *AuthService) Refresh(refreshToken string) (AuthResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanExpiredLocked()

	session, ok := s.refreshTokens[refreshToken]
	if !ok {
		return AuthResult{}, errors.New("invalid refresh token")
	}

	delete(s.refreshTokens, refreshToken)
	if session.RefreshToken != "" {
		delete(s.accessTokens, session.RefreshToken)
	}

	accessToken, err := newToken(32)
	if err != nil {
		return AuthResult{}, err
	}
	newRefreshToken, err := newToken(32)
	if err != nil {
		return AuthResult{}, err
	}

	accessSession := authSession{User: session.User, RefreshToken: newRefreshToken, ExpiresAt: time.Now().UTC().Add(s.accessTTL)}
	refreshSession := authSession{User: session.User, RefreshToken: accessToken, ExpiresAt: time.Now().UTC().Add(s.refreshTTL)}

	s.accessTokens[accessToken] = accessSession
	s.refreshTokens[newRefreshToken] = refreshSession

	return AuthResult{AccessToken: accessToken, RefreshToken: newRefreshToken, ExpiresIn: int64(s.accessTTL.Seconds())}, nil
}

func (s *AuthService) ResolveUser(accessToken string) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanExpiredLocked()
	session, ok := s.accessTokens[accessToken]
	if !ok {
		return domain.User{}, errors.New("unauthorized")
	}
	return session.User, nil
}

func (s *AuthService) ChangePassword(oldPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("新密码长度至少 6 位")
	}

	s.mu.RLock()
	store := s.credentials
	s.mu.RUnlock()

	if store != nil {
		hash, ok := store.GetAdminPasswordHash()
		if !ok || bcrypt.CompareHashAndPassword([]byte(hash), []byte(oldPassword)) != nil {
			return errors.New("旧密码不正确")
		}
		newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("密码加密失败")
		}
		if err := store.SetAdminPasswordHash(string(newHash)); err != nil {
			return errors.New("密码保存失败")
		}
	} else {
		s.mu.Lock()
		if oldPassword != s.adminPassword {
			s.mu.Unlock()
			return errors.New("旧密码不正确")
		}
		s.adminPassword = newPassword
		s.mu.Unlock()
	}

	// Invalidate all existing sessions so re-login is required
	s.mu.Lock()
	s.accessTokens = make(map[string]authSession)
	s.refreshTokens = make(map[string]authSession)
	s.mu.Unlock()
	return nil
}

func (s *AuthService) issueTokens(user domain.User) (AuthResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	accessToken, err := newToken(32)
	if err != nil {
		return AuthResult{}, err
	}
	refreshToken, err := newToken(32)
	if err != nil {
		return AuthResult{}, err
	}

	accessSession := authSession{User: user, RefreshToken: refreshToken, ExpiresAt: time.Now().UTC().Add(s.accessTTL)}
	refreshSession := authSession{User: user, RefreshToken: accessToken, ExpiresAt: time.Now().UTC().Add(s.refreshTTL)}

	s.accessTokens[accessToken] = accessSession
	s.refreshTokens[refreshToken] = refreshSession

	return AuthResult{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: int64(s.accessTTL.Seconds())}, nil
}

func (s *AuthService) cleanExpiredLocked() {
	now := time.Now().UTC()
	for token, session := range s.accessTokens {
		if now.After(session.ExpiresAt) {
			delete(s.accessTokens, token)
		}
	}
	for token, session := range s.refreshTokens {
		if now.After(session.ExpiresAt) {
			delete(s.refreshTokens, token)
		}
	}
}

func newToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
