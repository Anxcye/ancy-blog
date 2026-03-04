/**
 * File: auth.ts
 * Purpose: Store access token and user identity for admin authentication state.
 * Module: frontend-admin-react/store/auth, client state layer.
 * Related: login page, axios interceptor, and route auth guard.
 */
import { create } from 'zustand';

interface AuthState {
  accessToken: string;
  username: string;
  setAuth: (token: string, username: string) => void;
  clearAuth: () => void;
}

const TOKEN_KEY = 'ancy_admin_access_token';
const USER_KEY = 'ancy_admin_username';

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: localStorage.getItem(TOKEN_KEY) || '',
  username: localStorage.getItem(USER_KEY) || '',
  setAuth: (token, username) => {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, username);
    set({ accessToken: token, username });
  },
  clearAuth: () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    set({ accessToken: '', username: '' });
  },
}));
