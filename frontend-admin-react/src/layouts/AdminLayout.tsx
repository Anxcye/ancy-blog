/**
 * File: AdminLayout.tsx
 * Purpose: Provide shared admin shell layout with navigation and logout action.
 * Module: frontend-admin-react/layouts/admin, presentation frame layer.
 * Related: dashboard route, auth store, and future module pages.
 */
import { LogoutOutlined, MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons';
import { Button, Layout, Menu, Space, Typography } from 'antd';
import type { ReactElement } from 'react';
import { useMemo, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

import { useAuthStore } from '../store/auth';

const { Header, Sider, Content } = Layout;

export function AdminLayout(): ReactElement {
  const [collapsed, setCollapsed] = useState(false);
  const username = useAuthStore((state) => state.username);
  const clearAuth = useAuthStore((state) => state.clearAuth);
  const navigate = useNavigate();
  const location = useLocation();

  const selectedKeys = useMemo(() => {
    if (location.pathname.startsWith('/content')) {
      return ['content'];
    }
    if (location.pathname.startsWith('/site')) {
      return ['site'];
    }
    return ['dashboard'];
  }, [location.pathname]);

  function handleLogout(): void {
    clearAuth();
    navigate('/login');
  }

  return (
    <Layout className="admin-layout">
      <Sider trigger={null} collapsible collapsed={collapsed} breakpoint="lg" collapsedWidth={72}>
        <div className="sider-logo">Ancy</div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={selectedKeys}
          items={[
            { key: 'dashboard', label: '工作台', onClick: () => navigate('/') },
            { key: 'content', label: '内容管理', onClick: () => navigate('/content/articles') },
            { key: 'site', label: '站点设置', onClick: () => navigate('/site') },
          ]}
        />
      </Sider>
      <Layout>
        <Header className="admin-header">
          <Button type="text" icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />} onClick={() => setCollapsed((v) => !v)} />
          <Space>
            <Typography.Text type="secondary">{username || 'admin'}</Typography.Text>
            <Button type="text" icon={<LogoutOutlined />} onClick={handleLogout}>
              退出
            </Button>
          </Space>
        </Header>
        <Content className="admin-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
