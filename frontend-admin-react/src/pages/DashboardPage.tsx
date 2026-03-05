/**
 * File: DashboardPage.tsx
 * Purpose: Admin dashboard showing real-time stats, recent content, and quick actions.
 * Module: frontend-admin-react/pages/dashboard, presentation layer.
 * Related: dashboard API module, AdminLayout, articles/moments/comments APIs.
 */

import {
  CheckCircleOutlined,
  CommentOutlined,
  EditOutlined,
  FileTextOutlined,
  MessageOutlined,
  PlusOutlined,
  RightOutlined,
  ThunderboltOutlined,
  ClockCircleOutlined,
  FireOutlined,
} from '@ant-design/icons';
import { useQuery } from '@tanstack/react-query';
import { Badge, Button, Card, Col, Row, Skeleton, Statistic, Tag, Typography } from 'antd';
import type { ReactElement } from 'react';
import { useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';

import {
  fetchDashboardStats,
  fetchRecentArticles,
  fetchRecentComments,
} from '../api/dashboard';
import type { ArticleListItem } from '../types/article';
import type { Comment } from '../types/comment';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

// ─────────────────────────────────────────────────────────────
// Sub-components
// ─────────────────────────────────────────────────────────────

interface StatCardProps {
  title: string;
  value: number | undefined;
  icon: ReactElement;
  accentColor: string;
  accentBg: string;
  subtitle?: string;
  loading?: boolean;
  onClick?: () => void;
}

function StatCard({
  title,
  value,
  icon,
  accentColor,
  accentBg,
  subtitle,
  loading,
  onClick,
}: StatCardProps): ReactElement {
  return (
    <Card
      className="dash-stat-card"
      style={{ borderRadius: 16, cursor: onClick ? 'pointer' : 'default' }}
      styles={{ body: { padding: '20px 24px' } }}
      onClick={onClick}
    >
      <div className="dash-stat-inner">
        {/* Icon badge */}
        <div
          className="dash-stat-icon"
          style={{ background: accentBg, color: accentColor }}
        >
          {icon}
        </div>

        {/* Numbers */}
        <div className="dash-stat-text">
          {loading ? (
            <Skeleton.Input active size="small" style={{ width: 60, marginBottom: 4 }} />
          ) : (
            <Statistic
              value={value ?? 0}
              styles={{
                content: {
                  fontSize: 28,
                  fontWeight: 800,
                  color: '#1a1d23',
                  lineHeight: '1.15',
                  letterSpacing: '-0.5px',
                },
              }}
            />
          )}
          <span className="dash-stat-title">{title}</span>
          {subtitle && <span className="dash-stat-subtitle">{subtitle}</span>}
        </div>

        {/* Arrow hint if clickable */}
        {onClick && (
          <RightOutlined
            style={{ color: '#cbd5e1', fontSize: 12, marginLeft: 'auto', alignSelf: 'flex-start' }}
          />
        )}
      </div>
    </Card>
  );
}

// ─── Article Status Tag ───────────────────────────────────────

const STATUS_MAP: Record<string, { label: string; color: string }> = {
  published: { label: '已发布', color: 'success' },
  draft: { label: '草稿', color: 'default' },
  scheduled: { label: '定时', color: 'processing' },
  archived: { label: '归档', color: 'warning' },
};

function ArticleStatusTag({ status }: { status: string }): ReactElement {
  const cfg = STATUS_MAP[status] ?? { label: status, color: 'default' };
  return <Tag color={cfg.color} style={{ borderRadius: 6, fontSize: 11 }}>{cfg.label}</Tag>;
}

// ─── Recent Article Row ────────────────────────────────────

function ArticleRow({ item, navigate }: { item: ArticleListItem; navigate: (path: string) => void }): ReactElement {
  return (
    <div
      className="dash-row"
      onClick={() => navigate(`/content/articles/${item.id}/edit`)}
      role="button"
      tabIndex={0}
    >
      <div className="dash-row-main">
        <span className="dash-row-title">{item.title}</span>
        <span className="dash-row-time">
          <ClockCircleOutlined style={{ marginRight: 4, fontSize: 11 }} />
          {dayjs(item.updatedAt).fromNow()}
        </span>
      </div>
      <ArticleStatusTag status={item.status} />
    </div>
  );
}

// ─── Comment Status Tag ────────────────────────────────────

const COMMENT_STATUS_MAP: Record<string, { label: string; color: string }> = {
  pending: { label: '待审', color: 'warning' },
  approved: { label: '已通过', color: 'success' },
  rejected: { label: '已拒绝', color: 'error' },
  spam: { label: '垃圾', color: 'error' },
};

function CommentRow({ item }: { item: Comment }): ReactElement {
  const cfg = COMMENT_STATUS_MAP[item.status] ?? { label: item.status, color: 'default' };
  return (
    <div className="dash-row dash-row--comment">
      <div className="dash-row-avatar">
        {(item.nickname?.[0] ?? '?').toUpperCase()}
      </div>
      <div className="dash-row-body">
        <div className="dash-row-comment-header">
          <span className="dash-row-nickname">{item.nickname}</span>
          <Tag color={cfg.color} style={{ borderRadius: 6, fontSize: 11, marginLeft: 6 }}>{cfg.label}</Tag>
          <span className="dash-row-time" style={{ marginLeft: 'auto' }}>
            {dayjs(item.createdAt).fromNow()}
          </span>
        </div>
        <p className="dash-row-comment-content">{item.content}</p>
      </div>
    </div>
  );
}

// ─────────────────────────────────────────────────────────────
// Main Dashboard Page
// ─────────────────────────────────────────────────────────────

export function DashboardPage(): ReactElement {
  const navigate = useNavigate();

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard', 'stats'],
    queryFn: fetchDashboardStats,
    staleTime: 60_000,
  });

  const { data: recentArticles, isLoading: articlesLoading } = useQuery({
    queryKey: ['dashboard', 'recent-articles'],
    queryFn: fetchRecentArticles,
    staleTime: 60_000,
  });

  const { data: recentComments, isLoading: commentsLoading } = useQuery({
    queryKey: ['dashboard', 'recent-comments'],
    queryFn: fetchRecentComments,
    staleTime: 60_000,
  });

  // Greeting based on time of day
  const hour = new Date().getHours();
  const greeting =
    hour < 6 ? '夜深了' : hour < 12 ? '早上好' : hour < 18 ? '下午好' : '晚上好';

  return (
    <div className="dash-page">

      {/* ═══ Hero greeting ═══ */}
      <div className="dash-hero">
        <div className="dash-hero-left">
          <Typography.Title level={3} className="dash-hero-title">
            {greeting}，<span className="dash-hero-accent">欢迎回来</span> 👋
          </Typography.Title>
          <Typography.Text className="dash-hero-sub">
            今天是 {dayjs().format('YYYY年MM月DD日 dddd')}，博客一切正常。
          </Typography.Text>
        </div>
        <div className="dash-hero-actions">
          <Button
            type="primary"
            icon={<PlusOutlined />}
            size="middle"
            onClick={() => navigate('/content/articles/new')}
            className="dash-cta-btn"
          >
            写文章
          </Button>
          <Button
            icon={<ThunderboltOutlined />}
            size="middle"
            onClick={() => navigate('/content/moments')}
            className="dash-cta-btn-sec"
          >
            发瞬间
          </Button>
        </div>
      </div>

      {/* ═══ Primary stat row ═══ */}
      <Row gutter={[16, 16]} className="dash-stats-row">
        <Col xs={24} sm={8}>
          <StatCard
            title="文章总数"
            value={stats?.totalArticles}
            icon={<FileTextOutlined />}
            accentColor="#1f8f8a"
            accentBg="rgba(31,143,138,0.1)"
            subtitle={`${stats?.publishedArticles ?? '—'} 已发布`}
            loading={statsLoading}
            onClick={() => navigate('/content/articles')}
          />
        </Col>
        <Col xs={24} sm={8}>
          <StatCard
            title="瞬间"
            value={stats?.totalMoments}
            icon={<ThunderboltOutlined />}
            accentColor="#7c3aed"
            accentBg="rgba(124,58,237,0.09)"
            loading={statsLoading}
            onClick={() => navigate('/content/moments')}
          />
        </Col>
        <Col xs={24} sm={8}>
          <StatCard
            title="待审评论"
            value={stats?.pendingComments}
            icon={<MessageOutlined />}
            accentColor="#f59e0b"
            accentBg="rgba(245,158,11,0.1)"
            subtitle={`共 ${stats?.totalComments ?? '—'} 条`}
            loading={statsLoading}
            onClick={() => navigate('/interaction')}
          />
        </Col>
      </Row>

      {/* ═══ Secondary info row ═══ */}
      <Row gutter={[16, 16]} style={{ marginTop: 12 }}>
        {/* Secondary stat: drafts */}
        <Col xs={12} sm={6}>
          <Card
            className="dash-mini-card"
            styles={{ body: { padding: '16px 18px' } }}
            style={{ borderRadius: 14 }}
          >
            <div className="dash-mini-inner">
              <EditOutlined className="dash-mini-icon" style={{ color: '#64748b' }} />
              <div>
                <div className="dash-mini-value">
                  {statsLoading ? <Skeleton.Input active size="small" style={{ width: 36 }} /> : (stats?.draftArticles ?? 0)}
                </div>
                <div className="dash-mini-label">草稿</div>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={12} sm={6}>
          <Card
            className="dash-mini-card"
            styles={{ body: { padding: '16px 18px' } }}
            style={{ borderRadius: 14 }}
          >
            <div className="dash-mini-inner">
              <CheckCircleOutlined className="dash-mini-icon" style={{ color: '#22c55e' }} />
              <div>
                <div className="dash-mini-value">
                  {statsLoading ? <Skeleton.Input active size="small" style={{ width: 36 }} /> : (stats?.publishedArticles ?? 0)}
                </div>
                <div className="dash-mini-label">已发布</div>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={12} sm={6}>
          <Card
            className="dash-mini-card"
            styles={{ body: { padding: '16px 18px' } }}
            style={{ borderRadius: 14 }}
          >
            <div className="dash-mini-inner">
              <CommentOutlined className="dash-mini-icon" style={{ color: '#3b82f6' }} />
              <div>
                <div className="dash-mini-value">
                  {statsLoading ? <Skeleton.Input active size="small" style={{ width: 36 }} /> : (stats?.totalComments ?? 0)}
                </div>
                <div className="dash-mini-label">全部评论</div>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={12} sm={6}>
          <Card
            className="dash-mini-card"
            styles={{ body: { padding: '16px 18px' } }}
            style={{ borderRadius: 14 }}
          >
            <div className="dash-mini-inner">
              <FireOutlined className="dash-mini-icon" style={{ color: '#f97316' }} />
              <div>
                <div className="dash-mini-value">
                  {statsLoading ? <Skeleton.Input active size="small" style={{ width: 36 }} /> : (stats?.totalArticles ?? 0)}
                </div>
                <div className="dash-mini-label">全部内容</div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      {/* ═══ Content panels row ═══ */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>

        {/* Recent Articles */}
        <Col xs={24} lg={14}>
          <Card
            className="dash-panel"
            style={{ borderRadius: 16 }}
            styles={{ body: { padding: 0 } }}
            title={
              <div className="dash-panel-header">
                <span className="dash-panel-title">
                  <FileTextOutlined style={{ color: '#1f8f8a', marginRight: 8 }} />
                  最近文章
                </span>
                <Button
                  type="link"
                  size="small"
                  onClick={() => navigate('/content/articles')}
                  style={{ color: '#64748b', padding: 0 }}
                >
                  查看全部 <RightOutlined style={{ fontSize: 10 }} />
                </Button>
              </div>
            }
          >
            {articlesLoading ? (
              <div style={{ padding: '12px 24px' }}>
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={i} active avatar={false} paragraph={{ rows: 1 }} style={{ marginBottom: 8 }} />
                ))}
              </div>
            ) : recentArticles && recentArticles.length > 0 ? (
              <div className="dash-list">
                {recentArticles.map((item) => (
                  <ArticleRow key={item.id} item={item} navigate={navigate} />
                ))}
              </div>
            ) : (
              <div className="dash-empty">
                <FileTextOutlined style={{ fontSize: 32, color: '#cbd5e1' }} />
                <span>暂无文章，去创作第一篇吧</span>
                <Button
                  type="primary"
                  size="small"
                  icon={<PlusOutlined />}
                  onClick={() => navigate('/content/articles/new')}
                >
                  写文章
                </Button>
              </div>
            )}
          </Card>
        </Col>

        {/* Recent Comments */}
        <Col xs={24} lg={10}>
          <Card
            className="dash-panel"
            style={{ borderRadius: 16, height: '100%' }}
            styles={{ body: { padding: 0 } }}
            title={
              <div className="dash-panel-header">
                <span className="dash-panel-title">
                  <MessageOutlined style={{ color: '#f59e0b', marginRight: 8 }} />
                  最近评论
                  {(stats?.pendingComments ?? 0) > 0 && (
                    <Badge
                      count={stats?.pendingComments}
                      size="small"
                      style={{ marginLeft: 8, backgroundColor: '#f59e0b' }}
                    />
                  )}
                </span>
                <Button
                  type="link"
                  size="small"
                  onClick={() => navigate('/interaction')}
                  style={{ color: '#64748b', padding: 0 }}
                >
                  审核 <RightOutlined style={{ fontSize: 10 }} />
                </Button>
              </div>
            }
          >
            {commentsLoading ? (
              <div style={{ padding: '12px 24px' }}>
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={i} active avatar paragraph={{ rows: 1 }} style={{ marginBottom: 8 }} />
                ))}
              </div>
            ) : recentComments && recentComments.length > 0 ? (
              <div className="dash-list">
                {recentComments.map((item) => (
                  <CommentRow key={item.id} item={item} />
                ))}
              </div>
            ) : (
              <div className="dash-empty">
                <CommentOutlined style={{ fontSize: 32, color: '#cbd5e1' }} />
                <span>暂无评论</span>
              </div>
            )}
          </Card>
        </Col>
      </Row>

      {/* ═══ Quick action shortcuts ═══ */}
      <div className="dash-shortcuts">
        <Typography.Text className="dash-shortcuts-label">快捷操作</Typography.Text>
        <div className="dash-shortcuts-grid">
          {[
            { icon: <PlusOutlined />, label: '新建文章', path: '/content/articles/new', color: '#1f8f8a' },
            { icon: <ThunderboltOutlined />, label: '发布瞬间', path: '/content/moments', color: '#7c3aed' },
            { icon: <MessageOutlined />, label: '处理评论', path: '/interaction', color: '#f59e0b' },
            { icon: <FileTextOutlined />, label: '内容管理', path: '/content/articles', color: '#3b82f6' },
          ].map((s) => (
            <button
              key={s.path}
              className="dash-shortcut-btn"
              onClick={() => navigate(s.path)}
            >
              <span className="dash-shortcut-icon" style={{ color: s.color, background: `${s.color}14` }}>
                {s.icon}
              </span>
              <span className="dash-shortcut-label">{s.label}</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
