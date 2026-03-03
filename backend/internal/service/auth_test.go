// File: auth_test.go
// Purpose: Verify authentication service behavior for login, token refresh, and user resolution.
// Module: backend/internal/service, authentication unit test layer.
// Related: auth.go and middleware auth enforcement.
package service

import (
	"testing"
	"time"
)

func TestAuthServiceLoginAndResolveUser(t *testing.T) {
	svc := NewAuthService("admin", "secret", time.Hour, 24*time.Hour)

	result, err := svc.Login("admin", "secret")
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatalf("expected tokens to be issued, got: %#v", result)
	}
	if result.ExpiresIn != int64(time.Hour.Seconds()) {
		t.Fatalf("unexpected expiresIn: %d", result.ExpiresIn)
	}

	user, err := svc.ResolveUser(result.AccessToken)
	if err != nil {
		t.Fatalf("expected resolve user success, got error: %v", err)
	}
	if user.Username != "admin" || !user.IsAdmin {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestAuthServiceRefreshRotatesTokens(t *testing.T) {
	svc := NewAuthService("admin", "secret", time.Hour, 24*time.Hour)
	loginResult, err := svc.Login("admin", "secret")
	if err != nil {
		t.Fatalf("expected login success, got error: %v", err)
	}

	refreshResult, err := svc.Refresh(loginResult.RefreshToken)
	if err != nil {
		t.Fatalf("expected refresh success, got error: %v", err)
	}
	if refreshResult.AccessToken == loginResult.AccessToken {
		t.Fatalf("expected access token rotation")
	}
	if refreshResult.RefreshToken == loginResult.RefreshToken {
		t.Fatalf("expected refresh token rotation")
	}

	if _, err := svc.ResolveUser(loginResult.AccessToken); err == nil {
		t.Fatalf("expected old access token to be invalid after refresh")
	}
	if _, err := svc.ResolveUser(refreshResult.AccessToken); err != nil {
		t.Fatalf("expected new access token to be valid, got error: %v", err)
	}
}

func TestAuthServiceInvalidCredentials(t *testing.T) {
	svc := NewAuthService("admin", "secret", time.Hour, 24*time.Hour)
	if _, err := svc.Login("admin", "wrong"); err == nil {
		t.Fatalf("expected invalid credentials error")
	}
}
