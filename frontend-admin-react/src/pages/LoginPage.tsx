/**
 * File: LoginPage.tsx
 * Purpose: Provide admin sign-in form and bootstrap authenticated session.
 * Module: frontend-admin-react/pages/auth, presentation layer.
 * Related: auth store, http client, and backend auth login endpoint.
 */

import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { Button, Card, Form, Input, Typography, message } from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { httpClient } from '../lib/http';
import { useAuthStore } from '../store/auth';

type LoginResponse = {
  code: string;
  message: string;
  data: {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
  };
};

export function LoginPage(): ReactElement {
  const [loading, setLoading] = useState(false);
  const [messageApi, contextHolder] = message.useMessage();
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);

  async function onFinish(values: { username: string; password: string }): Promise<void> {
    setLoading(true);
    try {
      const response = await httpClient.post<LoginResponse>('/auth/login', values);
      const payload = response.data.data;
      setAuth(payload.accessToken, values.username);
      messageApi.success('登录成功');
      navigate('/');
    } catch {
      messageApi.error('账号或密码有误，请重试');
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="login-page">
      {contextHolder}
      <Card className="login-card" bordered={false}>
        {/* Logo + title */}
        <div style={{ textAlign: 'center', marginBottom: 28 }}>
          <div
            style={{
              width: 48,
              height: 48,
              background: 'linear-gradient(135deg, #1f8f8a 0%, #2bb5b0 100%)',
              borderRadius: 14,
              display: 'grid',
              placeItems: 'center',
              color: '#fff',
              fontWeight: 800,
              fontSize: 22,
              margin: '0 auto 14px',
              boxShadow: '0 6px 18px rgba(31,143,138,0.30)',
            }}
          >
            A
          </div>
          <Typography.Title level={4} style={{ margin: 0, letterSpacing: '0.02em' }}>
            Ancy Blog
          </Typography.Title>
          <Typography.Text type="secondary" style={{ fontSize: 13 }}>
            管理后台 · 请使用管理员账号登录
          </Typography.Text>
        </div>

        <Form layout="vertical" onFinish={onFinish} autoComplete="off" size="large">
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined style={{ color: '#94a3b8' }} />} placeholder="用户名" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#94a3b8' }} />}
              placeholder="密码"
            />
          </Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            block
            loading={loading}
            style={{ height: 42, fontSize: 15, marginTop: 4 }}
          >
            登录
          </Button>
        </Form>
      </Card>
    </main>
  );
}
