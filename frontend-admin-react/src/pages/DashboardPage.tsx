/**
 * File: DashboardPage.tsx
 * Purpose: Provide admin dashboard landing page for rewrite baseline.
 * Module: frontend-admin-react/pages/dashboard, presentation layer.
 * Related: AdminLayout and future metrics/widgets modules.
 */
import { Card, Col, Row, Statistic, Typography } from 'antd';
import type { ReactElement } from 'react';

export function DashboardPage(): ReactElement {
  return (
    <section>
      <Typography.Title level={3}>工作台</Typography.Title>
      <Typography.Paragraph type="secondary">React 管理端重写基座已就绪，下一步迁移内容、互动、系统模块。</Typography.Paragraph>
      <Row gutter={[16, 16]}>
        <Col xs={24} md={8}>
          <Card>
            <Statistic title="文章总数" value={0} />
          </Card>
        </Col>
        <Col xs={24} md={8}>
          <Card>
            <Statistic title="草稿文章" value={0} />
          </Card>
        </Col>
        <Col xs={24} md={8}>
          <Card>
            <Statistic title="待审核评论" value={0} />
          </Card>
        </Col>
      </Row>
    </section>
  );
}
