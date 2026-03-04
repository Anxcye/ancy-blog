// File: api/modules/auth.ts
// Purpose: Encapsulate admin authentication HTTP requests.
// Module: frontend-admin/api/auth, domain gateway layer.
// Related: login view, app store session token, backend auth handler.
import { httpClient } from '@/api/http';

interface LoginResponse {
  data: {
    token: string;
  };
}

export async function login(username: string, password: string): Promise<string> {
  const response = await httpClient.post<LoginResponse>('/auth/login', { username, password });
  return response.data.data.token;
}
