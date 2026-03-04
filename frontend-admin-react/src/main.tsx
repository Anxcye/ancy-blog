/**
 * File: main.tsx
 * Purpose: Bootstrap React admin application providers and mount root router entry.
 * Module: frontend-admin-react/app, client startup layer.
 * Related: router, auth store, query client, and Ant Design theme config.
 */
import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ConfigProvider, theme } from 'antd';

import App from './App';
import './index.css';

const queryClient = new QueryClient();

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1f8f8a',
          borderRadius: 10,
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </ConfigProvider>
  </React.StrictMode>,
);
