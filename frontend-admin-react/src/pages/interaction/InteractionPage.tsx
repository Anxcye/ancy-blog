/**
 * File: InteractionPage.tsx
 * Purpose: Provide unified comment moderation and friend-link review interface.
 * Module: frontend-admin-react/pages/interaction, presentation layer.
 * Related: comments/links API modules, comment/link types, and AdminLayout.
 */

import {
  CheckOutlined,
  CloseOutlined,
  EyeOutlined,
  LinkOutlined,
  MessageOutlined,
  StopOutlined,
} from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Col,
  Descriptions,
  Drawer,
  Form,
  Input,
  Modal,
  Popconfirm,
  Row,
  Select,
  Space,
  Table,
  Tabs,
  Tag,
  Tooltip,
  Typography,
  message,
} from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { listComments, replyComment, updateComment } from '../../api/comments';
import { listLinks, reviewLink } from '../../api/links';
import type { Comment, CommentListParams, CommentStatus } from '../../types/comment';
import type { Link, LinkListParams, LinkReviewPayload, ReviewStatus } from '../../types/link';

function fmtDate(iso?: string): string {
  if (!iso) return '—';
  const d = new Date(iso);
  const p = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`;
}

// ─────────────────────────────────────────────
// Comments tab
// ─────────────────────────────────────────────

const COMMENT_STATUS_COLOR: Record<string, string> = {
  pending: 'processing',
  approved: 'success',
  rejected: 'error',
  spam: 'warning',
  deleted: 'default',
};
const COMMENT_STATUS_LABEL: Record<string, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  spam: '垃圾',
  deleted: '已删除',
};

function CommentsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [replyForm] = Form.useForm<{ content: string }>();
  const [params, setParams] = useState<CommentListParams>({
    page: 1,
    pageSize: 20,
    status: 'pending',
  });
  const [activeComment, setActiveComment] = useState<Comment | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['comments', params],
    queryFn: () => listComments(params),
  });

  const updateMut = useMutation({
    mutationFn: ({ id, status }: { id: string; status: CommentStatus }) =>
      updateComment(id, { status }),
    onSuccess: (comment, vars) => {
      messageApi.success(`评论已${COMMENT_STATUS_LABEL[vars.status] ?? vars.status}`);
      if (activeComment?.id === comment.id) {
        setActiveComment(comment);
      }
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
    onError: () => messageApi.error('操作失败'),
  });

  const replyMut = useMutation({
    mutationFn: ({ id, content }: { id: string; content: string }) =>
      replyComment(id, { content }),
    onSuccess: () => {
      messageApi.success('管理员回复已发送');
      replyForm.resetFields();
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
    onError: () => messageApi.error('回复失败，请重试'),
  });

  function openCommentDetail(comment: Comment): void {
    setActiveComment(comment);
    replyForm.resetFields();
  }

  function closeCommentDetail(): void {
    setActiveComment(null);
    replyForm.resetFields();
  }

  function submitReply(): void {
    if (!activeComment) return;
    replyForm.validateFields().then((values) => {
      replyMut.mutate({ id: activeComment.id, content: values.content });
    });
  }

  function openContentTarget(comment: Comment): void {
    if (comment.contentType === 'article') {
      navigate(`/content/articles/${comment.contentId}/edit`);
      return;
    }
    navigate(`/content/moments?edit=${encodeURIComponent(comment.contentId)}`);
  }

  const contentTypeLabel = activeComment?.contentType === 'moment' ? '瞬间' : '文章';

  const columns = [
    {
      title: '评论者',
      key: 'author',
      width: 130,
      render: (_: unknown, r: Comment) => (
        <Button type="link" style={{ padding: 0, height: 'auto', textAlign: 'left' }} onClick={() => openCommentDetail(r)}>
          <div>
            <Typography.Text strong style={{ fontSize: 13 }}>{r.nickname}</Typography.Text>
            <br />
            <Typography.Text type="secondary" style={{ fontSize: 11 }}>{r.ip}</Typography.Text>
          </div>
        </Button>
      ),
    },
    {
      title: '内容',
      dataIndex: 'content',
      key: 'content',
      render: (content: string, record: Comment) => (
        <Button type="link" style={{ padding: 0, height: 'auto', textAlign: 'left' }} onClick={() => openCommentDetail(record)}>
          <Typography.Paragraph
            ellipsis={{ rows: 2 }}
            style={{ margin: 0, fontSize: 13, maxWidth: 380 }}
          >
            {content}
          </Typography.Paragraph>
        </Button>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (s: CommentStatus) => (
        <Tag color={COMMENT_STATUS_COLOR[s]}>{COMMENT_STATUS_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '风险',
      dataIndex: 'riskScore',
      key: 'riskScore',
      width: 60,
      render: (score: number) => (
        <Typography.Text type={score > 50 ? 'danger' : score > 20 ? 'warning' : 'secondary'}>
          {score}
        </Typography.Text>
      ),
    },
    {
      title: '时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 140,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 148,
      render: (_: unknown, r: Comment) => (
        <Space size={2}>
          <Tooltip title="查看详情">
            <Button
              type="text"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => openCommentDetail(r)}
            />
          </Tooltip>
          {r.status !== 'approved' && (
            <Tooltip title="通过">
              <Button
                type="text"
                size="small"
                icon={<CheckOutlined />}
                style={{ color: '#1f8f8a' }}
                loading={updateMut.isPending}
                onClick={() => updateMut.mutate({ id: r.id, status: 'approved' })}
              />
            </Tooltip>
          )}
          {r.status !== 'rejected' && (
            <Tooltip title="拒绝">
              <Button
                type="text"
                size="small"
                icon={<CloseOutlined />}
                danger
                loading={updateMut.isPending}
                onClick={() => updateMut.mutate({ id: r.id, status: 'rejected' })}
              />
            </Tooltip>
          )}
          {r.status !== 'spam' && (
            <Tooltip title="标为垃圾">
              <Button
                type="text"
                size="small"
                icon={<StopOutlined />}
                style={{ color: '#f59e0b' }}
                loading={updateMut.isPending}
                onClick={() => updateMut.mutate({ id: r.id, status: 'spam' })}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  return (
    <>
      {ctx}
      <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
        <Col xs={14} sm={8} md={5}>
          <Select
            style={{ width: '100%' }}
            defaultValue="pending"
            options={[
              { value: '', label: '全部状态' },
              { value: 'pending', label: '待审核' },
              { value: 'approved', label: '已通过' },
              { value: 'rejected', label: '已拒绝' },
              { value: 'spam', label: '垃圾' },
            ]}
            onChange={(v) =>
              setParams((p) => ({ ...p, status: (v as CommentStatus) || '', page: 1 }))
            }
          />
        </Col>
      </Row>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data?.rows ?? []}
        loading={isLoading}
        pagination={{
          current: params.page,
          pageSize: params.pageSize,
          total: data?.total ?? 0,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (page, pageSize) => setParams((p) => ({ ...p, page, pageSize })),
        }}
        scroll={{ x: 680 }}
      />

      <Drawer
        title={activeComment ? `评论详情 · ${activeComment.nickname}` : '评论详情'}
        width={760}
        open={activeComment !== null}
        onClose={closeCommentDetail}
        destroyOnClose
        extra={
          activeComment ? (
            <Space size={4}>
              {activeComment.status !== 'approved' && (
                <Button
                  size="small"
                  icon={<CheckOutlined />}
                  loading={updateMut.isPending}
                  onClick={() => updateMut.mutate({ id: activeComment.id, status: 'approved' })}
                >
                  通过
                </Button>
              )}
              {activeComment.status !== 'rejected' && (
                <Button
                  size="small"
                  danger
                  icon={<CloseOutlined />}
                  loading={updateMut.isPending}
                  onClick={() => updateMut.mutate({ id: activeComment.id, status: 'rejected' })}
                >
                  拒绝
                </Button>
              )}
              {activeComment.status !== 'spam' && (
                <Button
                  size="small"
                  icon={<StopOutlined />}
                  loading={updateMut.isPending}
                  onClick={() => updateMut.mutate({ id: activeComment.id, status: 'spam' })}
                >
                  标垃圾
                </Button>
              )}
            </Space>
          ) : null
        }
      >
        {activeComment && (
          <Space direction="vertical" size={20} style={{ width: '100%' }}>
            <div>
              <Space wrap size={[8, 8]}>
                <Tag color={COMMENT_STATUS_COLOR[activeComment.status]}>{COMMENT_STATUS_LABEL[activeComment.status] ?? activeComment.status}</Tag>
                <Tag>{contentTypeLabel}</Tag>
                <Tag>{activeComment.source || 'web'}</Tag>
                {activeComment.toCommentNickname ? <Tag color="blue">回复 {activeComment.toCommentNickname}</Tag> : null}
              </Space>
            </div>

            <div>
              <Typography.Text type="secondary">评论内容</Typography.Text>
              <div
                style={{
                  marginTop: 8,
                  padding: 16,
                  borderRadius: 12,
                  border: '1px solid #f0f0f0',
                  background: '#fafafa',
                  whiteSpace: 'pre-wrap',
                  lineHeight: 1.8,
                }}
              >
                {activeComment.content || '—'}
              </div>
            </div>

            <Descriptions
              title="访客与审核信息"
              size="small"
              bordered
              column={1}
              items={[
                { key: 'comment-id', label: '评论 ID', children: <Typography.Text copyable>{activeComment.id}</Typography.Text> },
                {
                  key: 'content-id',
                  label: '内容目标',
                  children: (
                    <Space>
                      <Typography.Text copyable>{`${contentTypeLabel} · ${activeComment.contentId}`}</Typography.Text>
                      <Button type="link" style={{ paddingInline: 0 }} onClick={() => openContentTarget(activeComment)}>
                        打开目标
                      </Button>
                    </Space>
                  ),
                },
                { key: 'parent-id', label: '父评论 ID', children: activeComment.parentId ? <Typography.Text copyable>{activeComment.parentId}</Typography.Text> : '—' },
                { key: 'root-id', label: '根评论 ID', children: activeComment.rootId ? <Typography.Text copyable>{activeComment.rootId}</Typography.Text> : '—' },
                { key: 'nickname', label: '昵称', children: activeComment.nickname || '—' },
                { key: 'email', label: '邮箱', children: activeComment.email || '—' },
                { key: 'website', label: '网站', children: activeComment.website ? <Typography.Link href={activeComment.website} target="_blank">{activeComment.website}</Typography.Link> : '—' },
                { key: 'ip', label: 'IP', children: activeComment.ip || '—' },
                { key: 'ua', label: 'User-Agent', children: activeComment.userAgent ? <Typography.Paragraph copyable style={{ marginBottom: 0, whiteSpace: 'pre-wrap' }}>{activeComment.userAgent}</Typography.Paragraph> : '—' },
                { key: 'risk', label: '风险分', children: String(activeComment.riskScore ?? 0) },
                { key: 'reply-count', label: '回复数', children: String(activeComment.replyCount ?? 0) },
                { key: 'created-at', label: '提交时间', children: fmtDate(activeComment.createdAt) },
                { key: 'approved-at', label: '审核时间', children: activeComment.approvedAt ? fmtDate(activeComment.approvedAt) : '—' },
              ]}
            />

            <Form form={replyForm} layout="vertical">
              <Form.Item
                name="content"
                label="管理员回复"
                rules={[{ required: true, message: '请填写回复内容' }]}
              >
                <Input.TextArea
                  rows={5}
                  maxLength={2000}
                  showCount
                  placeholder="以管理员身份回复这条评论，提交后会直接公开显示。"
                />
              </Form.Item>
              <Space>
                <Button type="primary" loading={replyMut.isPending} onClick={submitReply}>
                  发送回复
                </Button>
                <Typography.Text type="secondary">
                  管理员回复会作为线程中的子评论创建，来源标记为 `admin`。
                </Typography.Text>
              </Space>
            </Form>
          </Space>
        )}
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// Links tab
// ─────────────────────────────────────────────

const REVIEW_STATUS_COLOR: Record<string, string> = {
  pending: 'processing',
  approved: 'success',
  rejected: 'error',
};
const REVIEW_STATUS_LABEL: Record<string, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
};

function LinksTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<{ reviewNote: string }>();
  const [params, setParams] = useState<LinkListParams>({
    page: 1,
    pageSize: 20,
    reviewStatus: 'pending',
  });
  const [rejectModal, setRejectModal] = useState<{ open: boolean; id: string }>({
    open: false,
    id: '',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['links', params],
    queryFn: () => listLinks(params),
  });

  const reviewMut = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: LinkReviewPayload }) =>
      reviewLink(id, payload),
    onSuccess: (_r, vars) => {
      messageApi.success(
        vars.payload.reviewStatus === 'approved' ? '友链申请已通过' : '友链申请已拒绝',
      );
      setRejectModal({ open: false, id: '' });
      queryClient.invalidateQueries({ queryKey: ['links'] });
    },
    onError: () => messageApi.error('操作失败'),
  });

  function handleRejectSubmit(): void {
    form.validateFields().then((vals) => {
      reviewMut.mutate({
        id: rejectModal.id,
        payload: { reviewStatus: 'rejected', reviewNote: vals.reviewNote },
      });
    });
  }

  const columns = [
    {
      title: '站点',
      key: 'site',
      render: (_: unknown, r: Link) => (
        <div>
          <Typography.Text strong style={{ fontSize: 13 }}>{r.name}</Typography.Text>
          <br />
          <Typography.Link href={r.url} target="_blank" style={{ fontSize: 11 }}>
            {r.url}
          </Typography.Link>
        </div>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      render: (desc: string) => (
        <Typography.Text type="secondary" style={{ fontSize: 13 }}>{desc || '—'}</Typography.Text>
      ),
    },
    {
      title: '审核状态',
      dataIndex: 'reviewStatus',
      key: 'reviewStatus',
      width: 100,
      render: (s: ReviewStatus) => (
        <Tag color={REVIEW_STATUS_COLOR[s]}>{REVIEW_STATUS_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '拒绝备注',
      dataIndex: 'reviewNote',
      key: 'reviewNote',
      width: 140,
      render: (note: string) => (
        <Typography.Text type="secondary" style={{ fontSize: 12 }}>{note || '—'}</Typography.Text>
      ),
    },
    {
      title: '提交时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 140,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 90,
      render: (_: unknown, r: Link) => (
        <Space size={4}>
          {r.reviewStatus !== 'approved' && (
            <Popconfirm
              title="确认通过这条友链申请？"
              onConfirm={() =>
                reviewMut.mutate({ id: r.id, payload: { reviewStatus: 'approved' } })
              }
            >
              <Button
                type="text"
                size="small"
                icon={<CheckOutlined />}
                style={{ color: '#1f8f8a' }}
                loading={reviewMut.isPending}
              />
            </Popconfirm>
          )}
          {r.reviewStatus !== 'rejected' && (
            <Tooltip title="拒绝并填写原因">
              <Button
                type="text"
                size="small"
                icon={<CloseOutlined />}
                danger
                onClick={() => {
                  form.resetFields();
                  setRejectModal({ open: true, id: r.id });
                }}
              />
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  return (
    <>
      {ctx}
      <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
        <Col xs={14} sm={8} md={5}>
          <Select
            style={{ width: '100%' }}
            defaultValue="pending"
            options={[
              { value: '', label: '全部状态' },
              { value: 'pending', label: '待审核' },
              { value: 'approved', label: '已通过' },
              { value: 'rejected', label: '已拒绝' },
            ]}
            onChange={(v) =>
              setParams((p) => ({ ...p, reviewStatus: (v as ReviewStatus) || '', page: 1 }))
            }
          />
        </Col>
      </Row>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data?.rows ?? []}
        loading={isLoading}
        pagination={{
          current: params.page,
          pageSize: params.pageSize,
          total: data?.total ?? 0,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (page, pageSize) => setParams((p) => ({ ...p, page, pageSize })),
        }}
        scroll={{ x: 680 }}
      />

      <Modal
        title="拒绝友链申请"
        open={rejectModal.open}
        onCancel={() => setRejectModal({ open: false, id: '' })}
        onOk={handleRejectSubmit}
        okText="确认拒绝"
        okButtonProps={{ danger: true, loading: reviewMut.isPending }}
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 12 }}>
          <Form.Item
            name="reviewNote"
            label="拒绝原因（将向申请方展示）"
            rules={[{ required: true, message: '请填写拒绝原因' }]}
          >
            <Input.TextArea rows={3} placeholder="例如：站点内容与本站方向不符" />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}

// ─────────────────────────────────────────────
// Page entry
// ─────────────────────────────────────────────

export function InteractionPage(): ReactElement {
  return (
    <section>
      <Typography.Title level={3} style={{ marginBottom: 16 }}>
        互动中心
      </Typography.Title>
      <Tabs
        items={[
          {
            key: 'comments',
            label: <Space size={6}><MessageOutlined />评论审核</Space>,
            children: <CommentsTab />,
          },
          {
            key: 'links',
            label: <Space size={6}><LinkOutlined />友链管理</Space>,
            children: <LinksTab />,
          },
        ]}
      />
    </section>
  );
}
