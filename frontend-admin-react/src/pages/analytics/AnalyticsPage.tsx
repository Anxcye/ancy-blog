/**
 * File: AnalyticsPage.tsx
 * Purpose: Present visitor analytics overview, top page paths, and raw visit records.
 * Module: frontend-admin-react/pages/analytics, presentation layer.
 * Related: analytics API module, analytics types, and AdminLayout.
 */

import { BarChartOutlined, ClusterOutlined, EyeOutlined, GlobalOutlined, MonitorOutlined } from '@ant-design/icons';
import { useQuery } from '@tanstack/react-query';
import { Button, Card, Col, Descriptions, Drawer, Input, Row, Select, Space, Statistic, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { ReactElement } from 'react';
import { useState } from 'react';

import { getAnalyticsOverview, listAnalyticsPages, listAnalyticsVisits } from '../../api/analytics';
import type { AnalyticsPathStat, AnalyticsVisit, AnalyticsVisitListParams } from '../../types/analytics';

function fmtDate(iso?: string): string {
  if (!iso) return '—';
  const d = new Date(iso);
  const p = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`;
}

function fmtDuration(seconds?: number): string {
  if (!seconds || seconds <= 0) return '0s';
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainSeconds = seconds % 60;
  if (hours > 0) return `${hours}h ${minutes}m ${remainSeconds}s`;
  if (minutes > 0) return `${minutes}m ${remainSeconds}s`;
  return `${remainSeconds}s`;
}

function hashString(input: string): number {
  let hash = 0;
  for (let i = 0; i < input.length; i += 1) {
    hash = ((hash << 5) - hash) + input.charCodeAt(i);
    hash |= 0;
  }
  return Math.abs(hash);
}

function tokenColors(value: string): { background: string; border: string; text: string } {
  const hash = hashString(value);
  const hue = hash % 360;
  return {
    background: `hsl(${hue} 75% 96%)`,
    border: `hsl(${hue} 60% 78%)`,
    text: `hsl(${hue} 65% 28%)`,
  };
}

function renderIdToken(value: string, secondary = false): ReactElement {
  const colors = tokenColors(value);
  return (
    <Space size={6} style={{ alignItems: 'center' }}>
      <span
        aria-hidden="true"
        style={{
          width: secondary ? 7 : 9,
          height: secondary ? 7 : 9,
          borderRadius: '50%',
          background: colors.text,
          boxShadow: `0 0 0 2px ${colors.background}`,
          flex: '0 0 auto',
        }}
      />
      <Typography.Text
        code
        style={{
          margin: 0,
          fontSize: secondary ? 12 : 13,
          lineHeight: 1.6,
          color: secondary ? 'rgba(0, 0, 0, 0.65)' : undefined,
        }}
      >
        {value}
      </Typography.Text>
    </Space>
  );
}

export function AnalyticsPage(): ReactElement {
  const [days, setDays] = useState(7);
  const [pathKeyword, setPathKeyword] = useState('');
  const [activeVisit, setActiveVisit] = useState<AnalyticsVisit | null>(null);
  const [visitFilters, setVisitFilters] = useState<AnalyticsVisitListParams>({
    page: 1,
    pageSize: 20,
  });
  const [draftFilters, setDraftFilters] = useState<AnalyticsVisitListParams>({
    page: 1,
    pageSize: 20,
  });

  const { data: overview, isLoading: overviewLoading } = useQuery({
    queryKey: ['analytics', 'overview', days],
    queryFn: () => getAnalyticsOverview(days),
  });

  const { data: pageData, isLoading: pagesLoading } = useQuery({
    queryKey: ['analytics', 'pages', days, pathKeyword],
    queryFn: () => listAnalyticsPages({ page: 1, pageSize: 10, days, path: pathKeyword || undefined }),
  });

  const { data: visitData, isLoading: visitsLoading } = useQuery({
    queryKey: ['analytics', 'visits', days, visitFilters],
    queryFn: () => listAnalyticsVisits({ ...visitFilters, days }),
  });

  const pageColumns: ColumnsType<AnalyticsPathStat> = [
    {
      title: '路径',
      dataIndex: 'path',
      key: 'path',
      render: (path: string, record) => (
        <Space direction="vertical" size={0}>
          <Typography.Text strong>{path}</Typography.Text>
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            {record.contentType || 'site'} {record.contentSlug ? `· ${record.contentSlug}` : ''}
          </Typography.Text>
        </Space>
      ),
    },
    { title: 'PV', dataIndex: 'pageViews', key: 'pageViews', width: 90 },
    { title: 'UV', dataIndex: 'uniqueVisitors', key: 'uniqueVisitors', width: 90 },
    { title: 'IP', dataIndex: 'uniqueIPs', key: 'uniqueIPs', width: 90 },
    {
      title: '累计停留',
      dataIndex: 'activeDurationSeconds',
      key: 'activeDurationSeconds',
      width: 120,
      render: (value: number) => fmtDuration(value),
    },
    {
      title: '最近访问',
      dataIndex: 'lastVisitedAt',
      key: 'lastVisitedAt',
      width: 160,
      render: (value: string) => fmtDate(value),
    },
  ];

  const visitColumns: ColumnsType<AnalyticsVisit> = [
    {
      title: '时间',
      dataIndex: 'occurredAt',
      key: 'occurredAt',
      width: 160,
      render: (value: string) => fmtDate(value),
    },
    {
      title: '路径',
      dataIndex: 'path',
      key: 'path',
      render: (path: string, record) => (
        <Space direction="vertical" size={0}>
          <Typography.Text>{path}</Typography.Text>
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            {record.eventType} · {record.contentType || 'site'}
          </Typography.Text>
        </Space>
      ),
    },
    {
      title: '停留时长',
      dataIndex: 'activeDurationSeconds',
      key: 'activeDurationSeconds',
      width: 120,
      render: (value: number) => fmtDuration(value),
    },
    {
      title: '访客',
      key: 'visitor',
      width: 220,
      render: (_value, record) => (
        <Space direction="vertical" size={0}>
          {renderIdToken(record.visitorId)}
          {renderIdToken(record.sessionId, true)}
        </Space>
      ),
    },
    {
      title: '来源',
      key: 'source',
      width: 220,
      render: (_value, record) => (
        <Space direction="vertical" size={0}>
          <Typography.Text>{record.ip}</Typography.Text>
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            {record.referrerHost || 'direct'}
          </Typography.Text>
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            {[record.countryName, record.regionName, record.cityName].filter(Boolean).join(' / ') || '—'}
          </Typography.Text>
        </Space>
      ),
    },
    {
      title: '设备',
      key: 'device',
      width: 150,
      render: (_value, record) => (
        <Space wrap>
          <Tag>{record.deviceType || 'unknown'}</Tag>
          <Tag>{record.browserName || 'other'}</Tag>
          <Tag>{record.osName || 'other'}</Tag>
          {record.isBot && <Tag color="warning">bot</Tag>}
        </Space>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      render: (_value, record) => (
        <Button
          type="text"
          size="small"
          icon={<EyeOutlined />}
          onClick={() => setActiveVisit(record)}
        >
          查看
        </Button>
      ),
    },
  ];

  return (
    <div>
      <Space direction="vertical" size={20} style={{ width: '100%' }}>
        <div>
          <Typography.Title level={3} style={{ marginBottom: 4 }}>访客统计</Typography.Title>
          <Typography.Text type="secondary">
            查看最近访问趋势、热门路径，以及每次页面访问的来源和停留时长。
          </Typography.Text>
        </div>

        <Space wrap>
          <Select
            value={days}
            onChange={setDays}
            options={[
              { value: 7, label: '最近 7 天' },
              { value: 30, label: '最近 30 天' },
              { value: 90, label: '最近 90 天' },
            ]}
            style={{ width: 140 }}
          />
          <Input.Search
            allowClear
            placeholder="按路径过滤，如 /articles"
            onSearch={setPathKeyword}
            style={{ width: 280 }}
          />
        </Space>

        <Row gutter={[16, 16]}>
          <Col xs={24} md={6}>
            <Card loading={overviewLoading}>
              <Statistic title="页面浏览 PV" value={overview?.pageViews ?? 0} prefix={<BarChartOutlined />} />
            </Card>
          </Col>
          <Col xs={24} md={6}>
            <Card loading={overviewLoading}>
              <Statistic title="独立访客 UV" value={overview?.uniqueVisitors ?? 0} prefix={<GlobalOutlined />} />
            </Card>
          </Col>
          <Col xs={24} md={6}>
            <Card loading={overviewLoading}>
              <Statistic title="独立 IP" value={overview?.uniqueIPs ?? 0} prefix={<ClusterOutlined />} />
            </Card>
          </Col>
          <Col xs={24} md={6}>
            <Card loading={overviewLoading}>
              <Statistic title="访问会话" value={overview?.uniqueSessions ?? 0} prefix={<MonitorOutlined />} />
            </Card>
          </Col>
        </Row>

        <Row gutter={[16, 16]}>
          <Col xs={24} xl={14}>
            <Card
              title="热门路径"
              extra={<Typography.Text type="secondary">Top 10</Typography.Text>}
            >
              <Table
                rowKey={(record) => `${record.path}:${record.contentType}:${record.contentId}:${record.contentSlug}`}
                columns={pageColumns}
                dataSource={pageData?.rows ?? overview?.topPaths ?? []}
                loading={pagesLoading || overviewLoading}
                pagination={false}
                scroll={{ x: 720 }}
              />
            </Card>
          </Col>
          <Col xs={24} xl={10}>
            <Card title="来源与设备">
              <Space direction="vertical" size={16} style={{ width: '100%' }}>
                <div>
                  <Typography.Text strong>来源主机</Typography.Text>
                  <div style={{ marginTop: 8 }}>
                    {(overview?.topReferrers ?? []).slice(0, 6).map((item) => (
                      <div key={item.referrerHost} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6 }}>
                        <Typography.Text>{item.referrerHost || 'direct'}</Typography.Text>
                        <Typography.Text type="secondary">{item.visits}</Typography.Text>
                      </div>
                    ))}
                  </div>
                </div>
                <div>
                  <Typography.Text strong>设备分布</Typography.Text>
                  <div style={{ marginTop: 8 }}>
                    {(overview?.deviceBreakdown ?? []).slice(0, 6).map((item) => (
                      <div key={item.deviceType} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6 }}>
                        <Typography.Text>{item.deviceType}</Typography.Text>
                        <Typography.Text type="secondary">{item.visits}</Typography.Text>
                      </div>
                    ))}
                  </div>
                </div>
              </Space>
            </Card>
          </Col>
        </Row>

        <Card title="最近访问事件">
          <Space direction="vertical" size={12} style={{ width: '100%', marginBottom: 16 }}>
            <Row gutter={[12, 12]}>
              <Col xs={24} md={8} xl={6}>
                <Input
                  allowClear
                  placeholder="路径"
                  value={draftFilters.path}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, path: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={6}>
                <Input
                  allowClear
                  placeholder="IP"
                  value={draftFilters.ip}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, ip: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={6}>
                <Input
                  allowClear
                  placeholder="Visitor ID"
                  value={draftFilters.visitorId}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, visitorId: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={6}>
                <Input
                  allowClear
                  placeholder="Session ID"
                  value={draftFilters.sessionId}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, sessionId: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Select
                  allowClear
                  style={{ width: '100%' }}
                  placeholder="设备"
                  value={draftFilters.deviceType}
                  onChange={(value) => setDraftFilters((prev) => ({ ...prev, deviceType: value }))}
                  options={[
                    { value: 'desktop', label: 'desktop' },
                    { value: 'mobile', label: 'mobile' },
                    { value: 'tablet', label: 'tablet' },
                    { value: 'bot', label: 'bot' },
                    { value: 'unknown', label: 'unknown' },
                  ]}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Select
                  allowClear
                  style={{ width: '100%' }}
                  placeholder="浏览器"
                  value={draftFilters.browserName}
                  onChange={(value) => setDraftFilters((prev) => ({ ...prev, browserName: value }))}
                  options={[
                    { value: 'Chrome', label: 'Chrome' },
                    { value: 'Safari', label: 'Safari' },
                    { value: 'Edge', label: 'Edge' },
                    { value: 'Firefox', label: 'Firefox' },
                    { value: 'Opera', label: 'Opera' },
                    { value: 'WeChat', label: 'WeChat' },
                    { value: 'Other', label: 'Other' },
                    { value: 'unknown', label: 'unknown' },
                  ]}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Select
                  allowClear
                  style={{ width: '100%' }}
                  placeholder="系统"
                  value={draftFilters.osName}
                  onChange={(value) => setDraftFilters((prev) => ({ ...prev, osName: value }))}
                  options={[
                    { value: 'Windows', label: 'Windows' },
                    { value: 'macOS', label: 'macOS' },
                    { value: 'iOS', label: 'iOS' },
                    { value: 'Android', label: 'Android' },
                    { value: 'Linux', label: 'Linux' },
                    { value: 'Other', label: 'Other' },
                    { value: 'unknown', label: 'unknown' },
                  ]}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Select
                  allowClear
                  style={{ width: '100%' }}
                  placeholder="事件"
                  value={draftFilters.eventType}
                  onChange={(value) => setDraftFilters((prev) => ({ ...prev, eventType: value }))}
                  options={[
                    { value: 'page_view', label: 'page_view' },
                    { value: 'page_ping', label: 'page_ping' },
                  ]}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Select
                  allowClear
                  style={{ width: '100%' }}
                  placeholder="Bot"
                  value={draftFilters.isBot}
                  onChange={(value) => setDraftFilters((prev) => ({ ...prev, isBot: value }))}
                  options={[
                    { value: 'false', label: '非 Bot' },
                    { value: 'true', label: 'Bot' },
                  ]}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Select
                  allowClear
                  style={{ width: '100%' }}
                  placeholder="内容类型"
                  value={draftFilters.contentType}
                  onChange={(value) => setDraftFilters((prev) => ({ ...prev, contentType: value }))}
                  options={[
                    { value: 'article', label: 'article' },
                    { value: 'moment', label: 'moment' },
                    { value: 'site', label: 'site' },
                  ]}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Input
                  allowClear
                  placeholder="国家"
                  value={draftFilters.countryName}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, countryName: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Input
                  allowClear
                  placeholder="省/州"
                  value={draftFilters.regionName}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, regionName: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Input
                  allowClear
                  placeholder="城市"
                  value={draftFilters.cityName}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, cityName: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={8} xl={4}>
                <Input
                  allowClear
                  placeholder="ISP"
                  value={draftFilters.isp}
                  onChange={(e) => setDraftFilters((prev) => ({ ...prev, isp: e.target.value || undefined }))}
                />
              </Col>
              <Col xs={24} md={24} xl={8}>
                <Space>
                  <Button
                    type="primary"
                    onClick={() => setVisitFilters({ ...draftFilters, page: 1, pageSize: visitFilters.pageSize ?? 20 })}
                  >
                    应用筛选
                  </Button>
                  <Button
                    onClick={() => {
                      const next = { page: 1, pageSize: visitFilters.pageSize ?? 20 };
                      setDraftFilters(next);
                      setVisitFilters(next);
                    }}
                  >
                    重置
                  </Button>
                </Space>
              </Col>
            </Row>
          </Space>
          <Table
            rowKey="eventId"
            columns={visitColumns}
            dataSource={visitData?.rows ?? []}
            loading={visitsLoading}
            pagination={{
              current: visitFilters.page ?? 1,
              pageSize: visitFilters.pageSize ?? 20,
              total: visitData?.total ?? 0,
              onChange: (page, pageSize) => setVisitFilters((prev) => ({ ...prev, page, pageSize })),
            }}
            scroll={{ x: 980 }}
          />
        </Card>

        <Drawer
          title="访问记录详情"
          width={720}
          open={!!activeVisit}
          onClose={() => setActiveVisit(null)}
        >
          {activeVisit && (
            <Space direction="vertical" size={16} style={{ width: '100%' }}>
              <Descriptions bordered size="small" column={1}>
                <Descriptions.Item label="发生时间">{fmtDate(activeVisit.occurredAt)}</Descriptions.Item>
                <Descriptions.Item label="接收时间">{fmtDate(activeVisit.receivedAt)}</Descriptions.Item>
                <Descriptions.Item label="最后活跃时间">{fmtDate(activeVisit.lastEngagedAt)}</Descriptions.Item>
                <Descriptions.Item label="活跃时长">{fmtDuration(activeVisit.activeDurationSeconds)}</Descriptions.Item>
                <Descriptions.Item label="路径">{activeVisit.path}</Descriptions.Item>
                <Descriptions.Item label="事件类型">{activeVisit.eventType}</Descriptions.Item>
                <Descriptions.Item label="内容类型">{activeVisit.contentType || 'site'}</Descriptions.Item>
                <Descriptions.Item label="内容标识">
                  {activeVisit.contentSlug || activeVisit.contentId || '—'}
                </Descriptions.Item>
                <Descriptions.Item label="页面标题">{activeVisit.pageTitle || '—'}</Descriptions.Item>
                <Descriptions.Item label="Route Name">{activeVisit.routeName || '—'}</Descriptions.Item>
                <Descriptions.Item label="Visitor ID">
                  {renderIdToken(activeVisit.visitorId)}
                </Descriptions.Item>
                <Descriptions.Item label="Session ID">
                  {renderIdToken(activeVisit.sessionId)}
                </Descriptions.Item>
                <Descriptions.Item label="IP">{activeVisit.ip || '—'}</Descriptions.Item>
                <Descriptions.Item label="归属地">
                  {[activeVisit.countryName, activeVisit.regionName, activeVisit.cityName].filter(Boolean).join(' / ') || '—'}
                </Descriptions.Item>
                <Descriptions.Item label="国家代码">{activeVisit.countryCode || '—'}</Descriptions.Item>
                <Descriptions.Item label="ISP">{activeVisit.isp || '—'}</Descriptions.Item>
                <Descriptions.Item label="来源页面">{activeVisit.referrer || '—'}</Descriptions.Item>
                <Descriptions.Item label="来源主机">{activeVisit.referrerHost || 'direct'}</Descriptions.Item>
                <Descriptions.Item label="设备画像">
                  <Space wrap>
                    <Tag>{activeVisit.deviceType || 'unknown'}</Tag>
                    <Tag>{activeVisit.browserName || 'other'}</Tag>
                    <Tag>{activeVisit.osName || 'other'}</Tag>
                    {activeVisit.isBot && <Tag color="warning">bot</Tag>}
                  </Space>
                </Descriptions.Item>
                <Descriptions.Item label="语言 / 时区">
                  {activeVisit.locale || '—'} / {activeVisit.timezone || '—'}
                </Descriptions.Item>
                <Descriptions.Item label="屏幕尺寸">
                  {activeVisit.screenWidth || 0} x {activeVisit.screenHeight || 0}
                </Descriptions.Item>
                <Descriptions.Item label="视口尺寸">
                  {activeVisit.viewportWidth || 0} x {activeVisit.viewportHeight || 0}
                </Descriptions.Item>
                <Descriptions.Item label="User-Agent">
                  <Typography.Paragraph style={{ marginBottom: 0, whiteSpace: 'pre-wrap', wordBreak: 'break-all' }}>
                    {activeVisit.userAgent || '—'}
                  </Typography.Paragraph>
                </Descriptions.Item>
              </Descriptions>
            </Space>
          )}
        </Drawer>
      </Space>
    </div>
  );
}
