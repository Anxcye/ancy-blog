/**
 * File: api.ts
 * Purpose: Define shared API response envelope and pagination types.
 * Module: frontend-admin-react/types, domain types layer.
 * Related: all API modules and react-query hooks.
 */

export interface ApiResponse<T> {
  code: string;
  message: string;
  data: T;
}

export interface PaginatedData<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

export type PaginatedResponse<T> = ApiResponse<PaginatedData<T>>;
