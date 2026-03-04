/**
 * File: main.tsx
 * Purpose: Bootstrap React admin application providers and mount root router entry.
 * Module: frontend-admin-react/app, client startup layer.
 * Related: router, auth store, query client, and Ant Design theme config.
 */
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ConfigProvider, theme } from 'antd';
import React from 'react';
import ReactDOM from 'react-dom/client';

import App from './App';
import './index.css';

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 30_000, retry: 1 } },
});

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#1f8f8a',
          colorPrimaryBg: 'rgba(31, 143, 138, 0.08)',
          colorBgContainer: '#ffffff',
          colorBgLayout: '#f4f6f9',
          colorBgElevated: '#ffffff',
          colorBorder: 'rgba(0,0,0,0.08)',
          colorBorderSecondary: 'rgba(0,0,0,0.05)',
          colorTextBase: '#1a1d23',
          colorTextSecondary: '#64748b',
          colorTextTertiary: '#94a3b8',
          borderRadius: 10,
          borderRadiusSM: 8,
          borderRadiusLG: 14,
          fontFamily:
            "'PingFang SC', 'Noto Sans SC', 'Microsoft YaHei', -apple-system, BlinkMacSystemFont, sans-serif",
          fontSize: 14,
          lineHeight: 1.6,
          boxShadow: '0 2px 12px rgba(0,0,0,0.06)',
          boxShadowSecondary: '0 1px 6px rgba(0,0,0,0.04)',
          motionDurationMid: '0.18s',
          motionEaseInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
        },
        components: {
          Layout: {
            siderBg: '#ffffff',
            headerBg: 'rgba(255,255,255,0.85)',
            bodyBg: '#f4f6f9',
          },
          Menu: {
            itemBg: 'transparent',
            subMenuItemBg: 'transparent',
            itemColor: '#64748b',
            itemHoverColor: '#1f8f8a',
            itemHoverBg: 'rgba(31,143,138,0.07)',
            itemSelectedColor: '#1f8f8a',
            itemSelectedBg: 'rgba(31,143,138,0.10)',
            itemActiveBg: 'rgba(31,143,138,0.12)',
            itemBorderRadius: 9,
            itemMarginInline: 10,
            itemMarginBlock: 2,
            itemPaddingInline: 12,
            groupTitleColor: '#94a3b8',
            iconSize: 16,
          },
          Card: {
            boxShadow: '0 1px 8px rgba(0,0,0,0.05)',
            boxShadowTertiary: 'none',
            paddingLG: 20,
          },
          Table: {
            headerBg: '#fafbfc',
            headerColor: '#64748b',
            headerBorderRadius: 0,
            rowHoverBg: 'rgba(31,143,138,0.03)',
            cellPaddingBlock: 12,
          },
          Button: {
            borderRadius: 8,
            controlHeight: 36,
            fontWeight: 500,
            primaryShadow: '0 2px 8px rgba(31,143,138,0.24)',
          },
          Input: {
            borderRadius: 8,
            controlHeight: 36,
          },
          Select: {
            borderRadius: 8,
            controlHeight: 36,
          },
          Drawer: {
            borderRadius: 16,
          },
          Statistic: {
            titleFontSize: 13,
          },
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </ConfigProvider>
  </React.StrictMode>,
);
