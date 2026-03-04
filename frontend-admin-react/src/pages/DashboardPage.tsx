/**
 * File: DashboardPage.tsx
 * Purpose: Provide admin dashboard landing page with overview stats and quick links.
 * Module: frontend-admin-react/pages/dashboard, presentation layer.
 * Related: AdminLayout and future metrics/widgets modules.
 */

import {
  FileTextOutlined,
  MessageOutlined,
  ThunderboltOutlined,
} from '@ant-design/icons';
import { Card, Col, Row, Statistic, Typography } from 'antd';
import type { ReactElement } from 'react';

const STATS = [
  {
    title: '文章总数',
    value: 0,
    icon: <FileTextOutlined />,
    color: '#1f8f8a',
    bg: 'rgba(31,143,138,0.08)',
  },
  {
    title: '瞬间',
    value: 0,
    icon: <ThunderboltOutlined />,
    color: '#7c3aed',
    bg: 'rgba(124,58,237,0.08)',
  },
  {
    title: '待审评论',
    value: 0,
    icon: <MessageOutlined />,
    color: '#f59e0b',
    bg: 'rgba(245,158,11,0.08)',
  },
];

export function DashboardPage(): ReactElement {
  return (
    <section>
      <Typography.Title level={3} style={{ marginBottom: 4 }}>
        工作台
      </Typography.Title>
      <Typography.Paragraph type="secondary" style={{ marginBottom: 24 }}>
        欢迎回来，今天也是创作的好日子。
      </Typography.Paragraph>

      <Row gutter={[16, 16]}>
        {STATS.map((s) => (
          <Col xs={24} sm={8} key={s.title}>
            <Card
              className="card-hover"
              style={{ borderRadius: 14 }}
              styles={{ body: { display: 'flex', alignItems: 'center', gap: 16, padding: '20px 24px' } }}
            >
              <div
                style={{
                  width: 44,
                  height: 44,
                  borderRadius: 12,
                  background: s.bg,
                  color: s.color,
                  display: 'grid',
                  placeItems: 'center',
                  fontSize: 20,
                  flexShrink: 0,
                }}
              >
                {s.icon}
              </div>
              <Statistic
                title={
                  <span style={{ fontSize: 13, color: '#64748b', fontWeight: 500 }}>
                    {s.title}
                  </span>
                }
                value={s.value}
                styles={{ content: { fontSize: 26, fontWeight: 700, color: '#1a1d23', lineHeight: '1.2' } }}
              />
            </Card>
          </Col>
        ))}
      </Row>
    </section>
  );
}
