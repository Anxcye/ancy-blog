/**
 * File: http.ts
 * Purpose: Provide shared axios client with auth header and unauthorized handling.
 * Module: frontend-admin-react/lib/http, network gateway layer.
 * Related: auth store and all API modules.
 */
import axios from 'axios';

import { useAuthStore } from '../store/auth';

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const httpClient = axios.create({
  baseURL,
  timeout: 15000,
});

httpClient.interceptors.request.use((config) => {
  const token = useAuthStore.getState().accessToken;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

httpClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status;
    if (status === 401) {
      useAuthStore.getState().clearAuth();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  },
);
