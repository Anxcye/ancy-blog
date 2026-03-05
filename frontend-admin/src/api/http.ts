// File: api/http.ts
// Purpose: Create shared Axios HTTP client for admin backend APIs.
// Module: frontend-admin/api, infrastructure client layer.
// Related: auth module, content modules, app store token source.
import axios from 'axios';

import { useAppStore } from '@/stores/app';

export const httpClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1',
  timeout: 15000,
});

httpClient.interceptors.request.use((config) => {
  const appStore = useAppStore();

  if (appStore.token) {
    config.headers.Authorization = `Bearer ${appStore.token}`;
  }

  return config;
});

httpClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error?.response?.status as number | undefined;
    if (status === 401) {
      const appStore = useAppStore();
      appStore.clearToken();
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  },
);
