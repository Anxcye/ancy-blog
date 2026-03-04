/**
 * File: LoginPage.tsx
 * Purpose: Provide admin sign-in form and bootstrap authenticated session.
 * Module: frontend-admin-react/pages/auth, presentation layer.
 * Related: auth store, http client, and backend auth login endpoint.
 */
import { useState } from 'react';
import type { ReactElement } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Card, Form, Input, Typography, message } from 'antd';

import { httpClient } from '../lib/http';
import { useAuthStore } from '../store/auth';

type LoginResponse = {
  code: string;
  message: string;
  data: {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    user: {
      id: string;
      username: string;
      displayName: string;
    };
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
      setAuth(payload.accessToken, payload.user.username);
      messageApi.success('登录成功');
      navigate('/');
    } catch {
      messageApi.error('登录失败，请检查账号密码');
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="login-page">
      {contextHolder}
      <Card className="login-card" bordered={false}>
        <Typography.Title level={3}>Ancy Blog 管理后台</Typography.Title>
        <Typography.Paragraph type="secondary">使用管理员账号登录并管理站点内容。</Typography.Paragraph>
        <Form layout="vertical" onFinish={onFinish} autoComplete="off">
          <Form.Item label="用户名" name="username" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input placeholder="admin" />
          </Form.Item>
          <Form.Item label="密码" name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password placeholder="••••••••" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            登录
          </Button>
        </Form>
      </Card>
    </main>
  );
}
