/**
 * File: InteractionPage.tsx
 * Purpose: Provide unified comment moderation and friend-link review interface.
 * Module: frontend-admin-react/pages/interaction, presentation layer.
 * Related: comments/links API modules, comment/link types, and AdminLayout.
 */

import {
  CheckOutlined,
  CloseOutlined,
  LinkOutlined,
  MessageOutlined,
  StopOutlined,
} from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Col,
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

import { listComments, updateComment } from '../../api/comments';
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
  const [params, setParams] = useState<CommentListParams>({
    page: 1,
    pageSize: 20,
    status: 'pending',
  });

  const { data, isLoading } = useQuery({
    queryKey: ['comments', params],
    queryFn: () => listComments(params),
  });

  const updateMut = useMutation({
    mutationFn: ({ id, status }: { id: string; status: CommentStatus }) =>
      updateComment(id, { status }),
    onSuccess: (_r, vars) => {
      messageApi.success(`评论已${COMMENT_STATUS_LABEL[vars.status] ?? vars.status}`);
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
    onError: () => messageApi.error('操作失败'),
  });

  const columns = [
    {
      title: '评论者',
      key: 'author',
      width: 130,
      render: (_: unknown, r: Comment) => (
        <div>
          <Typography.Text strong style={{ fontSize: 13 }}>{r.nickname}</Typography.Text>
          <br />
          <Typography.Text type="secondary" style={{ fontSize: 11 }}>{r.ip}</Typography.Text>
        </div>
      ),
    },
    {
      title: '内容',
      dataIndex: 'content',
      key: 'content',
      render: (content: string) => (
        <Typography.Paragraph
          ellipsis={{ rows: 2 }}
          style={{ margin: 0, fontSize: 13, maxWidth: 380 }}
        >
          {content}
        </Typography.Paragraph>
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
      width: 112,
      render: (_: unknown, r: Comment) => (
        <Space size={2}>
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
