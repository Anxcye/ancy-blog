/**
 * Purpose: Display paginated gallery photo list with filters, upload, batch operations, and drawer editor.
 * Module: frontend-admin-react/pages/gallery, presentation layer.
 * Related: gallery API module, gallery types, AdminLayout.
 */

import {
  CameraOutlined,
  DeleteOutlined,
  EditOutlined,
  PlusOutlined,
  ReloadOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Col,
  Drawer,
  Form,
  Image,
  Input,
  Popconfirm,
  Row,
  Select,
  Space,
  Switch,
  Table,
  Tag,
  Typography,
  Upload,
  message,
} from 'antd';
import type { UploadChangeParam } from 'antd/es/upload';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useState } from 'react';

import {
  batchUpdatePhotoStatus,
  createGalleryTag,
  deleteGalleryPhoto,
  getGalleryPhoto,
  listGalleryPhotos,
  listGalleryTags,
  updateGalleryPhoto,
  uploadGalleryPhoto,
} from '../../api/gallery';
import { useAuthStore } from '../../store/auth';
import type {
  GalleryPhoto,
  GalleryPhotoFormValues,
  GalleryPhotoListParams,
  GalleryTag,
  PhotoStatus,
  ProcessingStatus,
} from '../../types/gallery';

const STATUS_COLOR: Record<string, string> = {
  draft: 'default',
  published: 'success',
  hidden: 'warning',
};

const STATUS_LABEL: Record<string, string> = {
  draft: '草稿',
  published: '已发布',
  hidden: '已隐藏',
};

const PROC_COLOR: Record<string, string> = {
  pending: 'default',
  processing: 'processing',
  completed: 'success',
  failed: 'error',
};

const PROC_LABEL: Record<string, string> = {
  pending: '等待中',
  processing: '处理中',
  completed: '已完成',
  failed: '失败',
};

function fmtDate(iso?: string): string {
  if (!iso) return '—';
  const d = new Date(iso);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function GalleryPage(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const accessToken = useAuthStore((s) => s.accessToken);

  // List state
  const [params, setParams] = useState<GalleryPhotoListParams>({
    page: 1,
    pageSize: 20,
    status: '',
    tag: '',
    keyword: '',
  });
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);

  // Drawer state
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingPhoto, setEditingPhoto] = useState<GalleryPhoto | null>(null);
  const [form] = Form.useForm<GalleryPhotoFormValues>();

  // Upload state
  const [uploading, setUploading] = useState(false);

  // Queries
  const { data, isLoading } = useQuery({
    queryKey: ['gallery-photos', params],
    queryFn: () => listGalleryPhotos(params),
  });

  const { data: tags = [] } = useQuery({
    queryKey: ['gallery-tags'],
    queryFn: listGalleryTags,
  });

  // Mutations
  const deleteMut = useMutation({
    mutationFn: deleteGalleryPhoto,
    onSuccess: () => {
      messageApi.success('照片已删除');
      queryClient.invalidateQueries({ queryKey: ['gallery-photos'] });
    },
    onError: () => messageApi.error('删除失败'),
  });

  const batchStatusMut = useMutation({
    mutationFn: ({ ids, status }: { ids: string[]; status: string }) =>
      batchUpdatePhotoStatus(ids, status),
    onSuccess: (r) => {
      messageApi.success(`已更新 ${r.count} 张`);
      setSelectedKeys([]);
      queryClient.invalidateQueries({ queryKey: ['gallery-photos'] });
    },
    onError: () => messageApi.error('状态更新失败'),
  });

  const updateMut = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: GalleryPhotoFormValues }) =>
      updateGalleryPhoto(id, payload),
    onSuccess: () => {
      messageApi.success('照片已更新');
      setDrawerOpen(false);
      setEditingPhoto(null);
      queryClient.invalidateQueries({ queryKey: ['gallery-photos'] });
    },
    onError: () => messageApi.error('更新失败'),
  });

  const createTagMut = useMutation({
    mutationFn: createGalleryTag,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['gallery-tags'] });
    },
  });

  // Drawer helpers
  function openEdit(photo: GalleryPhoto): void {
    setEditingPhoto(photo);
    form.setFieldsValue({
      title: photo.title,
      slug: photo.slug,
      description: photo.description,
      status: photo.status,
      locationName: photo.locationName,
      locationCity: photo.locationCity,
      locationState: photo.locationState,
      locationCountry: photo.locationCountry,
      takenAt: photo.takenAt,
      cameraMake: photo.cameraMake,
      cameraModel: photo.cameraModel,
      lensModel: photo.lensModel,
      focalLength: photo.focalLength,
      aperture: photo.aperture,
      shutterSpeed: photo.shutterSpeed,
      iso: photo.iso,
      takenAtDisplay: photo.takenAtDisplay,
      cameraDisplay: photo.cameraDisplay,
      locationDisplay: photo.locationDisplay,
      exifDisplay: photo.exifDisplay,
      tagsDisplay: photo.tagsDisplay,
      sortOrder: photo.sortOrder,
      tagSlugs: photo.tagSlugs,
    });
    setDrawerOpen(true);
  }

  function closeDrawer(): void {
    setDrawerOpen(false);
    setEditingPhoto(null);
    form.resetFields();
  }

  function handleSave(): void {
    form.validateFields().then((values) => {
      if (editingPhoto) {
        updateMut.mutate({ id: editingPhoto.id, payload: values });
      }
    });
  }

  // Fetch full photo detail when opening editor (to get latest data)
  useEffect(() => {
    if (editingPhoto?.id) {
      getGalleryPhoto(editingPhoto.id).then((photo) => {
        setEditingPhoto(photo);
        form.setFieldsValue({
          title: photo.title,
          slug: photo.slug,
          description: photo.description,
          status: photo.status,
          locationName: photo.locationName,
          locationCity: photo.locationCity,
          locationState: photo.locationState,
          locationCountry: photo.locationCountry,
          takenAt: photo.takenAt,
          cameraMake: photo.cameraMake,
          cameraModel: photo.cameraModel,
          lensModel: photo.lensModel,
          focalLength: photo.focalLength,
          aperture: photo.aperture,
          shutterSpeed: photo.shutterSpeed,
          iso: photo.iso,
          takenAtDisplay: photo.takenAtDisplay,
          cameraDisplay: photo.cameraDisplay,
          locationDisplay: photo.locationDisplay,
          exifDisplay: photo.exifDisplay,
          tagsDisplay: photo.tagsDisplay,
          sortOrder: photo.sortOrder,
          tagSlugs: photo.tagSlugs,
        });
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editingPhoto?.id]);

  // Upload handler
  const handleUploadChange = useCallback(
    (info: UploadChangeParam) => {
      if (info.file.status === 'uploading') {
        setUploading(true);
      }
      if (info.file.status === 'done') {
        setUploading(false);
        messageApi.success(`${info.file.name} 上传成功`);
        queryClient.invalidateQueries({ queryKey: ['gallery-photos'] });
      }
      if (info.file.status === 'error') {
        setUploading(false);
        messageApi.error(`${info.file.name} 上传失败`);
      }
    },
    [messageApi, queryClient],
  );

  function patchFilter(patch: Partial<GalleryPhotoListParams>): void {
    setParams((prev) => ({ ...prev, ...patch, page: 1 }));
    setSelectedKeys([]);
  }

  const columns = [
    {
      title: '缩略图',
      dataIndex: 'displayUrl',
      key: 'thumbnail',
      width: 80,
      render: (url: string, record: GalleryPhoto) => (
        <Image
          src={url}
          alt={record.title}
          width={60}
          height={60}
          style={{ objectFit: 'cover', borderRadius: 4 }}
          fallback="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='60' height='60'%3E%3Crect width='60' height='60' fill='%23f0f0f0'/%3E%3Ctext x='50%25' y='50%25' dominant-baseline='middle' text-anchor='middle' fill='%23999' font-size='12'%3E无图%3C/text%3E%3C/svg%3E"
        />
      ),
    },
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      render: (title: string, record: GalleryPhoto) => (
        <Button
          type="link"
          style={{ padding: 0, height: 'auto', textAlign: 'left', whiteSpace: 'normal' }}
          onClick={() => openEdit(record)}
        >
          {title || record.slug || '未命名'}
        </Button>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (s: PhotoStatus) => (
        <Tag color={STATUS_COLOR[s] ?? 'default'}>{STATUS_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '处理状态',
      dataIndex: 'processingStatus',
      key: 'processingStatus',
      width: 90,
      render: (s: ProcessingStatus) => (
        <Tag color={PROC_COLOR[s] ?? 'default'}>{PROC_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '相机',
      key: 'camera',
      width: 160,
      render: (_: unknown, record: GalleryPhoto) => {
        const parts = [record.cameraMake, record.cameraModel].filter(Boolean);
        return parts.length > 0 ? parts.join(' ') : '—';
      },
    },
    {
      title: '标签',
      dataIndex: 'tagSlugs',
      key: 'tags',
      width: 160,
      render: (slugs: string[]) =>
        slugs?.length > 0 ? (
          <Space size={2} wrap>
            {slugs.map((s) => (
              <Tag key={s} color="blue">
                {tags.find((t) => t.slug === s)?.name ?? s}
              </Tag>
            ))}
          </Space>
        ) : (
          '—'
        ),
    },
    {
      title: '尺寸',
      key: 'dimensions',
      width: 100,
      render: (_: unknown, record: GalleryPhoto) =>
        record.width && record.height ? `${record.width}×${record.height}` : '—',
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
      render: (_: unknown, record: GalleryPhoto) => (
        <Space>
          <Button
            type="text"
            icon={<EditOutlined />}
            size="small"
            onClick={() => openEdit(record)}
          />
          <Popconfirm
            title="确定删除这张照片？"
            description={
              record.articleRefCount > 0
                ? `该照片被 ${record.articleRefCount} 篇文章引用，删除后文章中将无法显示`
                : undefined
            }
            onConfirm={() => deleteMut.mutate(record.id)}
          >
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

  const batchPending = batchStatusMut.isPending;

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
          flexWrap: 'wrap',
          gap: 8,
        }}
      >
        <Typography.Title level={3} style={{ margin: 0 }}>
          <CameraOutlined style={{ marginRight: 8 }} />
          画廊管理
        </Typography.Title>
        <Upload
          name="file"
          accept="image/*"
          showUploadList={false}
          action="/api/v1/admin/gallery/photos/upload"
          headers={{ Authorization: `Bearer ${accessToken}` }}
          onChange={handleUploadChange}
          multiple
        >
          <Button type="primary" icon={<UploadOutlined />} loading={uploading}>
            上传照片
          </Button>
        </Upload>
      </div>

      {/* Filter toolbar */}
      <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
        <Col xs={24} sm={10} md={8}>
          <Input.Search
            placeholder="搜索标题或关键词"
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
              { value: 'hidden', label: '已隐藏' },
            ]}
            onChange={(v) => patchFilter({ status: (v as PhotoStatus) || '' })}
          />
        </Col>
        {tags.length > 0 && (
          <Col xs={12} sm={7} md={4}>
            <Select
              style={{ width: '100%' }}
              placeholder="标签"
              allowClear
              options={tags.map((t: GalleryTag) => ({ value: t.slug, label: t.name }))}
              onChange={(v) => patchFilter({ tag: (v as string) || '' })}
            />
          </Col>
        )}
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
            <Typography.Text type="secondary">已选 {selectedKeys.length} 张</Typography.Text>
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
            <Button
              size="small"
              loading={batchPending}
              onClick={() => batchStatusMut.mutate({ ids: selectedKeys, status: 'hidden' })}
            >
              批量隐藏
            </Button>
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
        dataSource={data?.rows ?? []}
        loading={isLoading}
        pagination={{
          current: params.page,
          pageSize: params.pageSize,
          total: data?.total ?? 0,
          showSizeChanger: true,
          showTotal: (total) => `共 ${total} 张`,
          onChange: (page, pageSize) => setParams((prev) => ({ ...prev, page, pageSize })),
        }}
        scroll={{ x: 1000 }}
      />

      {/* Edit drawer */}
      <Drawer
        title={editingPhoto ? `编辑照片 — ${editingPhoto.title || editingPhoto.slug}` : '编辑照片'}
        width={640}
        open={drawerOpen}
        onClose={closeDrawer}
        extra={
          <Space>
            <Button onClick={closeDrawer}>取消</Button>
            <Button type="primary" loading={updateMut.isPending} onClick={handleSave}>
              保存
            </Button>
          </Space>
        }
      >
        {editingPhoto && (
          <>
            {/* Preview */}
            <div style={{ textAlign: 'center', marginBottom: 24 }}>
              <Image
                src={editingPhoto.displayUrl}
                alt={editingPhoto.title}
                style={{ maxWidth: '100%', maxHeight: 300, objectFit: 'contain', borderRadius: 8 }}
              />
            </div>

            <Form form={form} layout="vertical">
              {/* Basic info */}
              <Typography.Title level={5}>基本信息</Typography.Title>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="slug" label="Slug" rules={[{ required: true, message: '请输入slug' }]}>
                    <Input />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item name="description" label="描述">
                <Input.TextArea rows={3} />
              </Form.Item>
              <Form.Item name="status" label="状态">
                <Select
                  options={[
                    { value: 'draft', label: '草稿' },
                    { value: 'published', label: '已发布' },
                    { value: 'hidden', label: '已隐藏' },
                  ]}
                />
              </Form.Item>
              <Form.Item name="sortOrder" label="排序">
                <Input type="number" />
              </Form.Item>

              {/* Tags */}
              <Typography.Title level={5}>标签</Typography.Title>
              <Form.Item name="tagSlugs" label="标签">
                <Select
                  mode="multiple"
                  allowClear
                  placeholder="选择标签"
                  options={tags.map((t: GalleryTag) => ({ value: t.slug, label: t.name }))}
                  dropdownRender={(menu) => (
                    <>
                      {menu}
                      <div style={{ padding: '8px', borderTop: '1px solid #f0f0f0' }}>
                        <NewTagInline
                          onCreated={(tag) => {
                            createTagMut.mutate({ name: tag.name, slug: tag.slug });
                          }}
                        />
                      </div>
                    </>
                  )}
                />
              </Form.Item>

              {/* Location */}
              <Typography.Title level={5}>位置信息</Typography.Title>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="locationName" label="地点名称">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="locationCity" label="城市">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="locationState" label="省/州">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="locationCountry" label="国家">
                    <Input />
                  </Form.Item>
                </Col>
              </Row>

              {/* EXIF */}
              <Typography.Title level={5}>拍摄信息</Typography.Title>
              <Form.Item name="takenAt" label="拍摄时间">
                <Input placeholder="2024-01-15T10:30:00Z" />
              </Form.Item>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="cameraMake" label="相机品牌">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="cameraModel" label="相机型号">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="lensModel" label="镜头">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="focalLength" label="焦距">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="aperture" label="光圈">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="shutterSpeed" label="快门速度">
                    <Input />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="iso" label="ISO">
                    <Input />
                  </Form.Item>
                </Col>
              </Row>

              {/* Display switches */}
              <Typography.Title level={5}>展示开关</Typography.Title>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item name="takenAtDisplay" label="显示拍摄时间" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="cameraDisplay" label="显示相机信息" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="locationDisplay" label="显示位置信息" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="exifDisplay" label="显示EXIF信息" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item name="tagsDisplay" label="显示标签" valuePropName="checked">
                    <Switch />
                  </Form.Item>
                </Col>
              </Row>

              {/* Read-only info */}
              <Typography.Title level={5}>其他信息</Typography.Title>
              <Typography.Text type="secondary">
                处理状态：
                <Tag color={PROC_COLOR[editingPhoto.processingStatus] ?? 'default'}>
                  {PROC_LABEL[editingPhoto.processingStatus] ?? editingPhoto.processingStatus}
                </Tag>
                {editingPhoto.processingError && (
                  <span style={{ color: '#ff4d4f' }}> — {editingPhoto.processingError}</span>
                )}
              </Typography.Text>
              <br />
              <Typography.Text type="secondary">
                尺寸：{editingPhoto.width}×{editingPhoto.height} · 文章引用：{editingPhoto.articleRefCount}
              </Typography.Text>
              <br />
              <Typography.Text type="secondary">
                创建：{fmtDate(editingPhoto.createdAt)} · 更新：{fmtDate(editingPhoto.updatedAt)}
              </Typography.Text>
            </Form>
          </>
        )}
      </Drawer>
    </section>
  );
}

// Inline tag creator inside the Select dropdown
function NewTagInline({ onCreated }: { onCreated: (tag: { name: string; slug: string }) => void }): ReactElement {
  const [name, setName] = useState('');

  function handleAdd(): void {
    const trimmed = name.trim();
    if (!trimmed) return;
    const slug = trimmed
      .toLowerCase()
      .replace(/[^a-z0-9\u4e00-\u9fff]+/g, '-')
      .replace(/^-|-$/g, '');
    onCreated({ name: trimmed, slug });
    setName('');
  }

  return (
    <Space>
      <Input
        size="small"
        placeholder="新标签名称"
        value={name}
        onChange={(e) => setName(e.target.value)}
        onPressEnter={handleAdd}
        style={{ width: 140 }}
      />
      <Button size="small" type="text" icon={<PlusOutlined />} onClick={handleAdd}>
        添加
      </Button>
    </Space>
  );
}
