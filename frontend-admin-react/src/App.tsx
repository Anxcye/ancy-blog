/**
 * File: App.tsx
 * Purpose: Compose route tree and provide browser router for admin application.
 * Module: frontend-admin-react/app, routing composition layer.
 * Related: AdminLayout, login page, dashboard page, and auth-guarded routes.
 */
import type { ReactElement } from 'react';
import { Navigate, RouterProvider, createBrowserRouter } from 'react-router-dom';

import { AdminLayout } from './layouts/AdminLayout';
import { DashboardPage } from './pages/DashboardPage';
import { LoginPage } from './pages/LoginPage';
import { useAuthStore } from './store/auth';

function RequireAuth({ children }: { children: ReactElement }): ReactElement {
  const token = useAuthStore((state) => state.accessToken);
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return children;
}

const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: (
      <RequireAuth>
        <AdminLayout />
      </RequireAuth>
    ),
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
    ],
  },
]);

export default function App(): ReactElement {
  return <RouterProvider router={router} />;
}
