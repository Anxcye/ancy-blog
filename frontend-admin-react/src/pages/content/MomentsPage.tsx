/**
 * File: MomentsPage.tsx
 * Purpose: Display paginated moments list with inline Drawer for create and edit operations.
 * Module: frontend-admin-react/pages/content, presentation layer.
 * Related: moments API module, moment types, and AdminLayout.
 */

import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Col,
  Drawer,
  Form,
  Input,
  Popconfirm,
  Row,
  Select,
  Space,
  Table,
  Tag,
  Typography,
  message,
} from 'antd';
import type { ReactElement } from 'react';
import { useEffect, useState } from 'react';

import {
  batchDeleteMoments,
  batchStatusMoments,
  createMoment,
  deleteMoment,
  listMoments,
  updateMoment,
} from '../../api/moments';
import type { Moment, MomentFormValues, MomentListParams, MomentStatus } from '../../types/moment';

const STATUS_COLOR: Record<string, string> = {
  draft: 'default',
  published: 'success',
  archived: 'warning',
};

const STATUS_LABEL: Record<string, string> = {
  draft: '草稿',
  published: '已发布',
  archived: '已归档',
};

function fmtDate(iso?: string): string {
  if (!iso) return '—';
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function MomentsPage(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<MomentFormValues>();

  const [params, setParams] = useState<MomentListParams>({ page: 1, pageSize: 20, status: '' });
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingMoment, setEditingMoment] = useState<Moment | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['moments', params],
    queryFn: () => listMoments(params),
  });

  const saveMut = useMutation({
    mutationFn: (values: MomentFormValues) =>
      editingMoment ? updateMoment(editingMoment.id, values) : createMoment(values),
    onSuccess: () => {
      messageApi.success(editingMoment ? '瞬间已保存' : '瞬间已发布');
      setDrawerOpen(false);
      queryClient.invalidateQueries({ queryKey: ['moments'] });
    },
    onError: () => messageApi.error('保存失败，请重试'),
  });

  const deleteMut = useMutation({
    mutationFn: deleteMoment,
    onSuccess: () => {
      messageApi.success('已删除');
      queryClient.invalidateQueries({ queryKey: ['moments'] });
    },
    onError: () => messageApi.error('删除失败'),
  });

  const batchDelMut = useMutation({
    mutationFn: batchDeleteMoments,
    onSuccess: (r) => {
      messageApi.success(`已删除 ${r.count} 条`);
      setSelectedKeys([]);
      queryClient.invalidateQueries({ queryKey: ['moments'] });
    },
    onError: () => messageApi.error('批量删除失败'),
  });

  const batchStatusMut = useMutation({
    mutationFn: ({ ids, status }: { ids: string[]; status: string }) =>
      batchStatusMoments(ids, status),
    onSuccess: (r) => {
      messageApi.success(`已更新 ${r.count} 条`);
      setSelectedKeys([]);
      queryClient.invalidateQueries({ queryKey: ['moments'] });
    },
    onError: () => messageApi.error('状态更新失败'),
  });

  // Populate form when opening the drawer in edit mode
  useEffect(() => {
    if (drawerOpen) {
      if (editingMoment) {
        form.setFieldsValue({ content: editingMoment.content, status: editingMoment.status });
      } else {
        form.resetFields();
        form.setFieldValue('status', 'published');
      }
    }
  }, [drawerOpen, editingMoment, form]);

  function openCreate(): void {
    setEditingMoment(null);
    setDrawerOpen(true);
  }

  function openEdit(moment: Moment): void {
    setEditingMoment(moment);
    setDrawerOpen(true);
  }

  const batchPending = batchDelMut.isPending || batchStatusMut.isPending;

  const columns = [
    {
      title: '内容',
      dataIndex: 'content',
      key: 'content',
      render: (content: string, record: Moment) => (
        <Button
          type="link"
          style={{ padding: 0, height: 'auto', textAlign: 'left', whiteSpace: 'pre-wrap', maxWidth: 480 }}
          onClick={() => openEdit(record)}
        >
          <Typography.Text ellipsis style={{ maxWidth: 440 }}>
            {content}
          </Typography.Text>
        </Button>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (s: MomentStatus) => (
        <Tag color={STATUS_COLOR[s] ?? 'default'}>{STATUS_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '发布时间',
      dataIndex: 'publishedAt',
      key: 'publishedAt',
      width: 152,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 152,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 100,
      render: (_: unknown, record: Moment) => (
        <Space>
          <Button
            type="text"
            icon={<EditOutlined />}
            size="small"
            onClick={() => openEdit(record)}
          />
          <Popconfirm title="确定删除这条瞬间？" onConfirm={() => deleteMut.mutate(record.id)}>
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
          瞬间管理
        </Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          发布瞬间
        </Button>
      </div>

      {/* Filter toolbar */}
      <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
        <Col xs={12} sm={6} md={4}>
          <Select
            style={{ width: '100%' }}
            placeholder="状态"
            allowClear
            options={[
              { value: 'draft', label: '草稿' },
              { value: 'published', label: '已发布' },
              { value: 'archived', label: '已归档' },
            ]}
            onChange={(v) => {
              setParams((p) => ({ ...p, status: (v as MomentStatus) || '', page: 1 }));
              setSelectedKeys([]);
            }}
          />
        </Col>
      </Row>

      {/* Batch action bar */}
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
            <Typography.Text type="secondary">已选 {selectedKeys.length} 条</Typography.Text>
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
              title={`确定删除选中的 ${selectedKeys.length} 条瞬间？`}
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
          showTotal: (total) => `共 ${total} 条`,
          onChange: (page, pageSize) => setParams((prev) => ({ ...prev, page, pageSize })),
        }}
        scroll={{ x: 640 }}
      />

      {/* Create / Edit drawer */}
      <Drawer
        title={editingMoment ? '编辑瞬间' : '发布瞬间'}
        width={480}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={saveMut.isPending} onClick={() => form.submit()}>
              保存
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={(values) => saveMut.mutate(values)}
        >
          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请填写内容' }]}
          >
            <Input.TextArea
              rows={8}
              placeholder="分享一个想法、一句话、一段文字..."
              showCount
              maxLength={2000}
            />
          </Form.Item>

          <Form.Item name="status" label="状态">
            <Select
              options={[
                { value: 'published', label: '立即发布' },
                { value: 'draft', label: '存为草稿' },
              ]}
            />
          </Form.Item>
        </Form>
      </Drawer>
    </section>
  );
}
