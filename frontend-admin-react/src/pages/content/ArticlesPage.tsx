/**
 * File: ArticlesPage.tsx
 * Purpose: Display paginated article list with status filters, keyword search, and batch operations.
 * Module: frontend-admin-react/pages/content, presentation layer.
 * Related: articles API module, article types, AdminLayout, and ArticleEditorPage.
 */

import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Col, Input, Popconfirm, Row, Select, Space, Table, Tag, Typography, message } from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { batchDeleteArticles, batchStatusArticles, deleteArticle, listArticles } from '../../api/articles';
import type { ArticleListItem, ArticleListParams, ArticleStatus, ContentKind } from '../../types/article';

const STATUS_COLOR: Record<string, string> = {
  draft: 'default',
  published: 'success',
  scheduled: 'processing',
  archived: 'warning',
};

const STATUS_LABEL: Record<string, string> = {
  draft: '草稿',
  published: '已发布',
  scheduled: '定时发布',
  archived: '已归档',
};

const KIND_LABEL: Record<string, string> = {
  post: '文章',
  page: '页面',
};

// Compact date formatter — avoids an extra dayjs import while antd is not yet wired
function fmtDate(iso?: string): string {
  if (!iso) return '—';
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function ArticlesPage(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  const [params, setParams] = useState<ArticleListParams>({
    page: 1,
    pageSize: 20,
    status: '',
    contentKind: '',
    keyword: '',
  });
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);

  const { data, isLoading } = useQuery({
    queryKey: ['articles', params],
    queryFn: () => listArticles(params),
  });

  const deleteMut = useMutation({
    mutationFn: deleteArticle,
    onSuccess: () => {
      messageApi.success('文章已删除');
      queryClient.invalidateQueries({ queryKey: ['articles'] });
    },
    onError: () => messageApi.error('删除失败'),
  });

  const batchDelMut = useMutation({
    mutationFn: batchDeleteArticles,
    onSuccess: (r) => {
      messageApi.success(`已删除 ${r.count} 篇`);
      setSelectedKeys([]);
      queryClient.invalidateQueries({ queryKey: ['articles'] });
    },
    onError: () => messageApi.error('批量删除失败'),
  });

  const batchStatusMut = useMutation({
    mutationFn: ({ ids, status }: { ids: string[]; status: string }) =>
      batchStatusArticles(ids, status),
    onSuccess: (r) => {
      messageApi.success(`已更新 ${r.count} 篇`);
      setSelectedKeys([]);
      queryClient.invalidateQueries({ queryKey: ['articles'] });
    },
    onError: () => messageApi.error('状态更新失败'),
  });

  function patchFilter(patch: Partial<ArticleListParams>): void {
    setParams((prev) => ({ ...prev, ...patch, page: 1 }));
    setSelectedKeys([]);
  }

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (title: string, record: ArticleListItem) => (
        <Button
          type="link"
          style={{ padding: 0, height: 'auto', textAlign: 'left', whiteSpace: 'normal' }}
          onClick={() => navigate(`/content/articles/${record.id}/edit`)}
        >
          {title}
        </Button>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 110,
      render: (s: ArticleStatus) => (
        <Tag color={STATUS_COLOR[s] ?? 'default'}>{STATUS_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '类型',
      dataIndex: 'contentKind',
      key: 'contentKind',
      width: 80,
      render: (k: ContentKind) => <Tag color="blue">{KIND_LABEL[k] ?? k}</Tag>,
    },
    {
      title: '发布时间',
      dataIndex: 'publishedAt',
      key: 'publishedAt',
      width: 152,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      width: 152,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      render: (_: unknown, record: ArticleListItem) => (
        <Space>
          <Button
            type="text"
            icon={<EditOutlined />}
            size="small"
            onClick={() => navigate(`/content/articles/${record.id}/edit`)}
          />
          <Popconfirm title="确定删除这篇文章？" onConfirm={() => deleteMut.mutate(record.id)}>
            <Button
              type="text"
              danger
              icon={<DeleteOutlined />}
              size="small"
              loading={deleteMut.isPending}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const batchPending = batchDelMut.isPending || batchStatusMut.isPending;

  return (
    <section>
      {ctx}

      {/* Page header */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 16,
        }}
      >
        <Typography.Title level={3} style={{ margin: 0 }}>
          文章管理
        </Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/content/articles/new')}>
          写文章
        </Button>
      </div>

      {/* Filter toolbar */}
      <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
        <Col xs={24} sm={10} md={8}>
          <Input.Search
            placeholder="搜索标题关键词"
            allowClear
            onSearch={(v) => patchFilter({ keyword: v })}
          />
        </Col>
        <Col xs={12} sm={7} md={4}>
          <Select
            style={{ width: '100%' }}
            placeholder="状态"
            allowClear
            options={[
              { value: 'draft', label: '草稿' },
              { value: 'published', label: '已发布' },
              { value: 'scheduled', label: '定时发布' },
              { value: 'archived', label: '已归档' },
            ]}
            onChange={(v) => patchFilter({ status: (v as ArticleStatus) || '' })}
          />
        </Col>
        <Col xs={12} sm={7} md={4}>
          <Select
            style={{ width: '100%' }}
            placeholder="类型"
            allowClear
            options={[
              { value: 'post', label: '文章' },
              { value: 'page', label: '页面' },
            ]}
            onChange={(v) => patchFilter({ contentKind: (v as ContentKind) || '' })}
          />
        </Col>
      </Row>

      {/* Batch action bar — visible only when rows are selected */}
      {selectedKeys.length > 0 && (
        <div
          style={{
            marginBottom: 12,
            padding: '8px 16px',
            background: '#f0f9ff',
            borderRadius: 8,
            border: '1px solid #bae6fd',
          }}
        >
          <Space wrap>
            <Typography.Text type="secondary">已选 {selectedKeys.length} 项</Typography.Text>
            <Button
              size="small"
              loading={batchPending}
              onClick={() => batchStatusMut.mutate({ ids: selectedKeys, status: 'published' })}
            >
              批量发布
            </Button>
            <Button
              size="small"
              loading={batchPending}
              onClick={() => batchStatusMut.mutate({ ids: selectedKeys, status: 'draft' })}
            >
              退回草稿
            </Button>
            <Popconfirm
              title={`确定删除选中的 ${selectedKeys.length} 篇文章？`}
              onConfirm={() => batchDelMut.mutate(selectedKeys)}
            >
              <Button size="small" danger loading={batchPending}>
                批量删除
              </Button>
            </Popconfirm>
            <Button size="small" type="text" onClick={() => setSelectedKeys([])}>
              取消选择
            </Button>
          </Space>
        </div>
      )}

      <Table
        rowKey="id"
        rowSelection={{
          selectedRowKeys: selectedKeys,
          onChange: (keys) => setSelectedKeys(keys as string[]),
        }}
        columns={columns}
        dataSource={data?.items ?? []}
        loading={isLoading}
        pagination={{
          current: params.page,
          pageSize: params.pageSize,
          total: data?.total ?? 0,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 篇`,
          onChange: (page, pageSize) => setParams((prev) => ({ ...prev, page, pageSize })),
        }}
        scroll={{ x: 720 }}
      />
    </section>
  );
}
