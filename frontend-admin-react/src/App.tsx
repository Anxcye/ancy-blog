/**
 * File: App.tsx
 * Purpose: Compose route tree and provide browser router for admin application.
 * Module: frontend-admin-react/app, routing composition layer.
 * Related: AdminLayout, auth guard, and all module pages.
 */

import type { ReactElement } from 'react';
import { Navigate, RouterProvider, createBrowserRouter } from 'react-router-dom';

import { AdminLayout } from './layouts/AdminLayout';
import { ContentLayout } from './layouts/ContentLayout';
import { ArticleEditorPage } from './pages/content/ArticleEditorPage';
import { ArticlesPage } from './pages/content/ArticlesPage';
import { MomentsPage } from './pages/content/MomentsPage';
import { TaxonomyPage } from './pages/content/TaxonomyPage';
import { DashboardPage } from './pages/DashboardPage';
import { InteractionPage } from './pages/interaction/InteractionPage';
import { LoginPage } from './pages/LoginPage';
import { SitePage } from './pages/site/SitePage';
import { SystemPage } from './pages/system/SystemPage';
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
      { index: true, element: <DashboardPage /> },

      // Content module — ContentLayout provides the mobile tab switcher
      {
        path: 'content',
        element: <ContentLayout />,
        children: [
          { path: 'articles', element: <ArticlesPage /> },
          { path: 'articles/new', element: <ArticleEditorPage /> },
          { path: 'articles/:id/edit', element: <ArticleEditorPage /> },
          { path: 'moments', element: <MomentsPage /> },
          { path: 'taxonomy', element: <TaxonomyPage /> },
        ],
      },

      // Interaction module
      { path: 'interaction', element: <InteractionPage /> },

      // Site module
      { path: 'site', element: <SitePage /> },

      // System module
      { path: 'system', element: <SystemPage /> },
    ],
  },
]);

export default function App(): ReactElement {
  return <RouterProvider router={router} />;
}
