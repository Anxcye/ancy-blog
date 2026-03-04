/**
 * File: AdminLayout.tsx
 * Purpose: Provide admin shell with light sider, frosted header, and mobile bottom navigation.
 * Module: frontend-admin-react/layouts/admin, presentation frame layer.
 * Related: all module pages and auth store.
 */

import {
  AppstoreOutlined,
  DashboardOutlined,
  FileTextOutlined,
  GlobalOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  MessageOutlined,
  ThunderboltOutlined, // still used in desktop submenu item
} from '@ant-design/icons';
import { Button, Layout, Menu, Space, Typography } from 'antd';
import type { ReactElement } from 'react';
import { useMemo, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

import { useAuthStore } from '../store/auth';

const { Header, Sider, Content } = Layout;

// Resolve sidebar active key (granular — articles/moments are separate)
function resolveKey(pathname: string): string {
  if (pathname.startsWith('/content/articles')) return 'articles';
  if (pathname.startsWith('/content/moments')) return 'moments';
  if (pathname.startsWith('/interaction')) return 'interaction';
  if (pathname.startsWith('/site')) return 'site';
  if (pathname.startsWith('/system')) return 'system';
  return 'dashboard';
}

// Resolve mobile bottom nav active key (content section collapses to one tab)
function resolveMobileKey(pathname: string): string {
  if (pathname.startsWith('/content')) return 'content';
  if (pathname.startsWith('/interaction')) return 'interaction';
  if (pathname.startsWith('/site')) return 'site';
  if (pathname.startsWith('/system')) return 'system';
  return 'dashboard';
}

// 5 tabs mirroring the 5 desktop top-level sections
const MOBILE_NAV = [
  { key: 'dashboard', path: '/', icon: <DashboardOutlined />, label: '工作台' },
  { key: 'content', path: '/content/articles', icon: <FileTextOutlined />, label: '内容' },
  { key: 'interaction', path: '/interaction', icon: <MessageOutlined />, label: '互动' },
  { key: 'site', path: '/site', icon: <GlobalOutlined />, label: '站点' },
  { key: 'system', path: '/system', icon: <AppstoreOutlined />, label: '系统' },
];

export function AdminLayout(): ReactElement {
  const [collapsed, setCollapsed] = useState(false);
  const username = useAuthStore((s) => s.username);
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const navigate = useNavigate();
  const location = useLocation();

  const activeKey = useMemo(() => resolveKey(location.pathname), [location.pathname]);
  const mobileActiveKey = useMemo(() => resolveMobileKey(location.pathname), [location.pathname]);
  const openKeys = useMemo(
    () => (location.pathname.startsWith('/content') ? ['content'] : []),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [],
  );

  function handleLogout(): void {
    clearAuth();
    navigate('/login');
  }

  const menuItems = [
    {
      key: 'dashboard',
      icon: <DashboardOutlined />,
      label: '工作台',
      onClick: () => navigate('/'),
    },
    {
      key: 'content',
      icon: <FileTextOutlined />,
      label: '内容管理',
      children: [
        {
          key: 'articles',
          icon: <FileTextOutlined />,
          label: '文章',
          onClick: () => navigate('/content/articles'),
        },
        {
          key: 'moments',
          icon: <ThunderboltOutlined />,
          label: '瞬间',
          onClick: () => navigate('/content/moments'),
        },
      ],
    },
    {
      key: 'interaction',
      icon: <MessageOutlined />,
      label: '互动中心',
      onClick: () => navigate('/interaction'),
    },
    {
      key: 'site',
      icon: <GlobalOutlined />,
      label: '站点设置',
      onClick: () => navigate('/site'),
    },
    {
      key: 'system',
      icon: <AppstoreOutlined />,
      label: '系统设置',
      onClick: () => navigate('/system'),
    },
  ];

  // Derive initials for the avatar
  const initials = (username?.[0] ?? 'A').toUpperCase();

  return (
    <Layout className="admin-layout">
      {/* ——— Desktop sider ——— */}
      <Sider
        className="admin-sider"
        trigger={null}
        collapsible
        collapsed={collapsed}
        collapsedWidth={64}
        theme="light"
        width={220}
      >
        {/* Brand logo */}
        <div className="sider-brand" onClick={() => navigate('/')}>
          <div className="sider-brand-icon">A</div>
          {!collapsed && <span className="sider-brand-name">Ancy</span>}
        </div>

        <Menu
          mode="inline"
          selectedKeys={[activeKey]}
          defaultOpenKeys={openKeys}
          items={menuItems}
          style={{ border: 'none', flex: 1 }}
        />

        {/* Collapse toggle at the bottom of the sider */}
        <div className="sider-footer">
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed((v) => !v)}
            className="sider-collapse-btn"
          />
        </div>
      </Sider>

      <Layout>
        {/* ——— Top header ——— */}
        <Header className="admin-header">
          {/* Mobile: show brand name in header */}
          <span className="admin-header-brand">Ancy</span>

          <Space align="center" size={4}>
            <div className="admin-avatar" title={username || 'admin'}>
              {initials}
            </div>
            <Typography.Text className="admin-username">{username || 'admin'}</Typography.Text>
            <Button
              type="text"
              icon={<LogoutOutlined />}
              onClick={handleLogout}
              className="admin-logout-btn"
              size="small"
            />
          </Space>
        </Header>

        {/* ——— Page content ——— */}
        <Content className="admin-content">
          <Outlet />
        </Content>
      </Layout>

      {/* ——— Mobile bottom navigation ——— */}
      <nav className="mobile-bottom-nav" aria-label="底部导航">
        {MOBILE_NAV.map((item) => {
          const isActive = mobileActiveKey === item.key;
          return (
            <button
              key={item.key}
              className={`mnav-item${isActive ? ' is-active' : ''}`}
              onClick={() => navigate(item.path)}
              aria-current={isActive ? 'page' : undefined}
            >
              <span className="mnav-icon">{item.icon}</span>
              <span className="mnav-label">{item.label}</span>
            </button>
          );
        })}
      </nav>
    </Layout>
  );
}
