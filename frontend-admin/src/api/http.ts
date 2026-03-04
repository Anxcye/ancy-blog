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
  (error) => Promise.reject(error),
);
