/**
 * File: SitePage.tsx
 * Purpose: Provide site settings, social links, footer items, and nav items management.
 * Module: frontend-admin-react/pages/site, presentation layer.
 * Related: site API module, site types, and AdminLayout.
 */

import {
  DeleteOutlined,
  EditOutlined,
  GlobalOutlined,
  PlusOutlined,
  SaveOutlined,
} from '@ant-design/icons';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Avatar,
  Button,
  Card,
  Col,
  Drawer,
  Form,
  Input,
  InputNumber,
  Popconfirm,
  Row,
  Select,
  Space,
  Switch,
  Table,
  Tabs,
  Tag,
  Typography,
  message,
} from 'antd';
import type { ReactElement } from 'react';
import { useEffect, useState } from 'react';

import {
  createFooterItem,
  createNavItem,
  createSocialLink,
  deleteFooterItem,
  deleteNavItem,
  deleteSocialLink,
  getSiteSettings,
  listFooterItems,
  listNavItems,
  listSocialLinks,
  updateFooterItem,
  updateNavItem,
  updateSiteSettings,
  updateSocialLink,
} from '../../api/site';
import type {
  FooterItem,
  FooterItemFormValues,
  NavItem,
  NavItemFormValues,
  SocialLink,
  SocialLinkFormValues,
} from '../../types/site';

// ─────────────────────────────────────────────
// 1. Site settings tab
// ─────────────────────────────────────────────

function SettingsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const [form] = Form.useForm();

  const { data, isLoading } = useQuery({
    queryKey: ['site-settings'],
    queryFn: getSiteSettings,
  });

  useEffect(() => {
    if (data) form.setFieldsValue(data);
  }, [data, form]);

  const saveMut = useMutation({
    mutationFn: updateSiteSettings,
    onSuccess: () => messageApi.success('站点设置已保存'),
    onError: () => messageApi.error('保存失败'),
  });

  return (
    <>
      {ctx}
      <Card style={{ maxWidth: 560 }}>
        <Form
          form={form}
          layout="vertical"
          disabled={isLoading}
          onFinish={(vals) => saveMut.mutate(vals)}
        >
          <Form.Item
            name="siteName"
            label="站点名称"
            rules={[{ required: true, message: '请填写站点名称' }]}
          >
            <Input placeholder="Ancy Blog" />
          </Form.Item>
          <Form.Item name="avatarUrl" label="头像 URL">
            <Input placeholder="https://cdn.example.com/avatar.png" />
          </Form.Item>
          <Form.Item
            name="heroIntroMd"
            label={
              <Space>
                首页介绍
                <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                  支持 Markdown
                </Typography.Text>
              </Space>
            }
          >
            <Input.TextArea rows={4} placeholder="Hi, I'm Ancy. I build things." />
          </Form.Item>
          <Form.Item name="defaultLocale" label="默认语言">
            <Select
              options={[
                { value: 'zh-CN', label: '简体中文 (zh-CN)' },
                { value: 'en-US', label: 'English (en-US)' },
                { value: 'ja-JP', label: '日本語 (ja-JP)' },
              ]}
            />
          </Form.Item>

          <Typography.Text strong style={{ display: 'block', marginBottom: 12 }}>SEO / 元信息</Typography.Text>

          <Form.Item name="siteDescription" label="站点描述">
            <Input.TextArea rows={2} placeholder="简短描述，用于搜索引擎和分享卡片" />
          </Form.Item>
          <Form.Item name="seoKeywords" label="SEO 关键词">
            <Input placeholder="用英文逗号分隔，如: blog, tech, golang" />
          </Form.Item>
          <Form.Item name="ogImageUrl" label="默认 OG 图片 URL">
            <Input placeholder="https://cdn.example.com/og-default.png" />
          </Form.Item>

          <Button
            type="primary"
            htmlType="submit"
            icon={<SaveOutlined />}
            loading={saveMut.isPending}
          >
            保存设置
          </Button>
        </Form>
      </Card>
    </>
  );
}


// ─────────────────────────────────────────────
// 2. Comment policy tab
// ─────────────────────────────────────────────

function CommentPolicyTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const [form] = Form.useForm();

  const { data, isLoading } = useQuery({
    queryKey: ['site-settings'],
    queryFn: getSiteSettings,
  });

  useEffect(() => {
    if (data) form.setFieldsValue(data);
  }, [data, form]);

  const saveMut = useMutation({
    mutationFn: updateSiteSettings,
    onSuccess: () => messageApi.success('评论设置已保存'),
    onError: () => messageApi.error('保存失败'),
  });

  return (
    <>
      {ctx}
      <div style={{ maxWidth: 560 }}>
        <Typography.Paragraph type="secondary" style={{ marginBottom: 20 }}>
          控制全站评论的开启状态与审核策略。更改后立即生效，无需重启服务。
        </Typography.Paragraph>

        <Form
          form={form}
          disabled={isLoading}
          onFinish={(vals) => saveMut.mutate(vals)}
        >
          <div className="comment-policy-section">
            {/* Row 1: enable comments */}
            <div className="comment-policy-row">
              <div className="comment-policy-info">
                <span className="comment-policy-label">开启评论功能</span>
                <span className="comment-policy-desc">
                  关闭后，全站所有文章与瞬间均不显示评论区，访客无法提交新评论
                </span>
              </div>
              <Form.Item name="commentEnabled" valuePropName="checked" noStyle>
                <Switch />
              </Form.Item>
            </div>

            {/* Row 2: require approval */}
            <div className="comment-policy-row comment-policy-row--last">
              <div className="comment-policy-info">
                <span className="comment-policy-label">评论需审核后展示</span>
                <span className="comment-policy-desc">
                  开启后，新评论提交后默认为「待审核」状态，不会立即公开显示。
                  需前往<Typography.Text strong style={{ fontSize: 12 }}>互动中心</Typography.Text>手动通过审核
                </span>
              </div>
              <Form.Item name="commentRequireApproval" valuePropName="checked" noStyle>
                <Switch />
              </Form.Item>
            </div>
          </div>

          <Button
            type="primary"
            htmlType="submit"
            icon={<SaveOutlined />}
            loading={saveMut.isPending}
          >
            保存评论设置
          </Button>
        </Form>
      </div>
    </>
  );
}


// ─────────────────────────────────────────────
// 3. Social links tab
// ─────────────────────────────────────────────

const PLATFORM_LABELS: Record<string, string> = {
  github: 'GitHub',
  mail: 'Email',
  x: 'X (Twitter)',
  linkedin: 'LinkedIn',
  custom: '自定义',
};

function SocialLinksTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<SocialLinkFormValues>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<SocialLink | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['social-links'],
    queryFn: listSocialLinks,
  });

  const saveMut = useMutation({
    mutationFn: async (vals: SocialLinkFormValues): Promise<void> => {
      if (editing) {
        await updateSocialLink(editing.id, vals);
      } else {
        await createSocialLink(vals);
      }
    },
    onSuccess: () => {
      messageApi.success(editing ? '已更新' : '已添加');
      setDrawerOpen(false);
      queryClient.invalidateQueries({ queryKey: ['social-links'] });
    },
    onError: () => messageApi.error('保存失败'),
  });

  const deleteMut = useMutation({
    mutationFn: deleteSocialLink,
    onSuccess: () => {
      messageApi.success('已删除');
      queryClient.invalidateQueries({ queryKey: ['social-links'] });
    },
    onError: () => messageApi.error('删除失败'),
  });

  function openCreate(): void {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({ orderNum: 0, enabled: true });
    setDrawerOpen(true);
  }

  function openEdit(item: SocialLink): void {
    setEditing(item);
    form.setFieldsValue(item);
    setDrawerOpen(true);
  }

  const columns = [
    {
      title: '图标 URL',
      dataIndex: 'iconKey',
      key: 'iconKey',
      render: (key: string | undefined) =>
        key ? (
          <Avatar src={key} size="small" />
        ) : (
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>默认 / 无</Typography.Text>
        ),
    },
    {
      title: '平台',
      dataIndex: 'platform',
      key: 'platform',
      render: (p: string) => PLATFORM_LABELS[p] ?? p,
    },
    { title: '标题', dataIndex: 'title', key: 'title' },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
      render: (url: string) => (
        <Typography.Link href={url} target="_blank" style={{ fontSize: 12 }}>
          {url}
        </Typography.Link>
      ),
    },
    { title: '排序', dataIndex: 'orderNum', key: 'orderNum', width: 60 },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 60,
      render: (v: boolean) => <Tag color={v ? 'success' : 'default'}>{v ? '是' : '否'}</Tag>,
    },
    {
      title: '操作',
      key: 'actions',
      width: 90,
      render: (_: unknown, r: SocialLink) => (
        <Space size={4}>
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="确认删除？" onConfirm={() => deleteMut.mutate(r.id)}>
            <Button
              type="text"
              size="small"
              icon={<DeleteOutlined />}
              danger
              loading={deleteMut.isPending}
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      {ctx}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 12,
        }}
      >
        <Typography.Text type="secondary">管理首页展示的社交媒体链接</Typography.Text>
        <Button size="small" type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          添加链接
        </Button>
      </div>
      <Table
        rowKey="id"
        size="small"
        dataSource={data ?? []}
        loading={isLoading}
        columns={columns}
        pagination={false}
        scroll={{ x: 520 }}
      />
      <Drawer
        title={editing ? '编辑社交链接' : '添加社交链接'}
        width={400}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button
              type="primary"
              loading={saveMut.isPending}
              onClick={() => form.validateFields().then((v) => saveMut.mutate(v))}
            >
              保存
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="platform"
            label="平台"
            rules={[{ required: true, message: '请选择平台' }]}
          >
            <Select
              options={Object.entries(PLATFORM_LABELS).map(([v, l]) => ({ value: v, label: l }))}
            />
          </Form.Item>
          <Form.Item
            name="title"
            label="显示标题"
            rules={[{ required: true, message: '请填写标题' }]}
          >
            <Input placeholder="GitHub" />
          </Form.Item>
          <Form.Item
            name="url"
            label="URL"
            rules={[{ required: true, message: '请填写 URL' }]}
          >
            <Input placeholder="https://github.com/username" />
          </Form.Item>
          <Form.Item
            name="iconKey"
            label="图标 URL (可选)"
            tooltip="如果是自定义平台或想覆盖默认图标，请填写图片地址"
          >
            <Input placeholder="https://example.com/icon.png" />
          </Form.Item>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="orderNum" label="排序数">
                <InputNumber style={{ width: '100%' }} min={0} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="enabled" label="启用" valuePropName="checked">
                <Switch />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// 3. Footer items tab
// ─────────────────────────────────────────────

function FooterItemsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<FooterItemFormValues>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<FooterItem | null>(null);
  const [linkType, setLinkType] = useState<string>('none');

  const { data, isLoading } = useQuery({
    queryKey: ['footer-items'],
    queryFn: listFooterItems,
  });

  const saveMut = useMutation({
    mutationFn: async (vals: FooterItemFormValues): Promise<void> => {
      if (editing) {
        await updateFooterItem(editing.id, vals);
      } else {
        await createFooterItem(vals);
      }
    },
    onSuccess: () => {
      messageApi.success(editing ? '已更新' : '已添加');
      setDrawerOpen(false);
      queryClient.invalidateQueries({ queryKey: ['footer-items'] });
    },
    onError: () => messageApi.error('保存失败'),
  });

  const deleteMut = useMutation({
    mutationFn: deleteFooterItem,
    onSuccess: () => {
      messageApi.success('已删除');
      queryClient.invalidateQueries({ queryKey: ['footer-items'] });
    },
    onError: () => messageApi.error('删除失败'),
  });

  function openCreate(): void {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({ linkType: 'none', rowNum: 1, orderNum: 0, enabled: true });
    setLinkType('none');
    setDrawerOpen(true);
  }

  function openEdit(item: FooterItem): void {
    setEditing(item);
    form.setFieldsValue(item);
    setLinkType(item.linkType);
    setDrawerOpen(true);
  }

  const LINK_TYPE_LABEL: Record<string, string> = { none: '纯文本', internal: '站内链接', external: '外部链接' };

  const columns = [
    { title: '文字', dataIndex: 'label', key: 'label' },
    {
      title: '链接类型',
      dataIndex: 'linkType',
      key: 'linkType',
      width: 90,
      render: (t: string) => LINK_TYPE_LABEL[t] ?? t,
    },
    { title: '行', dataIndex: 'rowNum', key: 'rowNum', width: 50 },
    { title: '排序', dataIndex: 'orderNum', key: 'orderNum', width: 60 },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 60,
      render: (v: boolean) => <Tag color={v ? 'success' : 'default'}>{v ? '是' : '否'}</Tag>,
    },
    {
      title: '操作',
      key: 'actions',
      width: 90,
      render: (_: unknown, r: FooterItem) => (
        <Space size={4}>
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="确认删除？" onConfirm={() => deleteMut.mutate(r.id)}>
            <Button type="text" size="small" icon={<DeleteOutlined />} danger loading={deleteMut.isPending} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      {ctx}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
        <Typography.Text type="secondary">页脚共 3 行，每行可放置多个链接或文字项</Typography.Text>
        <Button size="small" type="primary" icon={<PlusOutlined />} onClick={openCreate}>添加项目</Button>
      </div>
      <Table
        rowKey="id"
        size="small"
        dataSource={data ?? []}
        loading={isLoading}
        columns={columns}
        pagination={false}
        scroll={{ x: 480 }}
      />
      <Drawer
        title={editing ? '编辑页脚项' : '添加页脚项'}
        width={400}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={saveMut.isPending} onClick={() => form.validateFields().then((v) => saveMut.mutate(v))}>
              保存
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="label" label="显示文字" rules={[{ required: true }]}>
            <Input placeholder="关于我" />
          </Form.Item>
          <Form.Item name="linkType" label="链接类型">
            <Select
              options={[
                { value: 'none', label: '纯文本（无链接）' },
                { value: 'internal', label: '站内页面' },
                { value: 'external', label: '外部链接' },
              ]}
              onChange={(v) => setLinkType(v)}
            />
          </Form.Item>
          {linkType === 'internal' && (
            <Form.Item name="internalArticleSlug" label="页面 Slug" rules={[{ required: true }]}>
              <Input placeholder="about" />
            </Form.Item>
          )}
          {linkType === 'external' && (
            <Form.Item name="externalUrl" label="外部 URL" rules={[{ required: true }]}>
              <Input placeholder="https://example.com" />
            </Form.Item>
          )}
          <Row gutter={12}>
            <Col span={8}>
              <Form.Item name="rowNum" label="行号 (1-3)">
                <Select options={[{ value: 1, label: '第 1 行' }, { value: 2, label: '第 2 行' }, { value: 3, label: '第 3 行' }]} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="orderNum" label="排序">
                <InputNumber style={{ width: '100%' }} min={0} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="enabled" label="启用" valuePropName="checked">
                <Switch />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// 4. Nav items tab
// ─────────────────────────────────────────────

function NavItemsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<NavItemFormValues>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editing, setEditing] = useState<NavItem | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['nav-items'],
    queryFn: listNavItems,
  });

  const saveMut = useMutation({
    mutationFn: async (vals: NavItemFormValues): Promise<void> => {
      if (editing) {
        await updateNavItem(editing.id, vals);
      } else {
        await createNavItem(vals);
      }
    },
    onSuccess: () => {
      messageApi.success(editing ? '已更新' : '已添加');
      setDrawerOpen(false);
      queryClient.invalidateQueries({ queryKey: ['nav-items'] });
    },
    onError: () => messageApi.error('保存失败'),
  });

  const deleteMut = useMutation({
    mutationFn: deleteNavItem,
    onSuccess: () => {
      messageApi.success('已删除');
      queryClient.invalidateQueries({ queryKey: ['nav-items'] });
    },
    onError: () => messageApi.error('删除失败'),
  });

  function openCreate(): void {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({ type: 'link', targetType: 'route', orderNum: 0, enabled: true, parentId: undefined });
    setDrawerOpen(true);
  }

  function openEdit(item: NavItem): void {
    setEditing(item);
    form.setFieldsValue(item);
    setDrawerOpen(true);
  }

  const TYPE_LABEL: Record<string, string> = { menu: '菜单', dropdown: '下拉', link: '链接' };
  const TARGET_LABEL: Record<string, string> = { route: '路由', category: '分类', article: '文章', external: '外链' };

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: 'Key', dataIndex: 'key', key: 'key', render: (v: string) => <Typography.Text code>{v}</Typography.Text> },
    { title: '类型', dataIndex: 'type', key: 'type', width: 70, render: (v: string) => TYPE_LABEL[v] ?? v },
    { title: '目标类型', dataIndex: 'targetType', key: 'targetType', width: 80, render: (v: string) => TARGET_LABEL[v] ?? v },
    { title: '目标值', dataIndex: 'targetValue', key: 'targetValue', render: (v: string) => v || '—' },
    { title: '排序', dataIndex: 'orderNum', key: 'orderNum', width: 60 },
    {
      title: '启用',
      dataIndex: 'enabled',
      key: 'enabled',
      width: 60,
      render: (v: boolean) => <Tag color={v ? 'success' : 'default'}>{v ? '是' : '否'}</Tag>,
    },
    {
      title: '操作',
      key: 'actions',
      width: 90,
      render: (_: unknown, r: NavItem) => (
        <Space size={4}>
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="确认删除？" onConfirm={() => deleteMut.mutate(r.id)}>
            <Button type="text" size="small" icon={<DeleteOutlined />} danger loading={deleteMut.isPending} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      {ctx}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
        <Typography.Text type="secondary">配置博客顶部导航栏的菜单项</Typography.Text>
        <Button size="small" type="primary" icon={<PlusOutlined />} onClick={openCreate}>添加导航项</Button>
      </div>
      <Table
        rowKey="id"
        size="small"
        dataSource={data ?? []}
        loading={isLoading}
        columns={columns}
        pagination={false}
        scroll={{ x: 600 }}
      />
      <Drawer
        title={editing ? '编辑导航项' : '添加导航项'}
        width={420}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={saveMut.isPending} onClick={() => form.validateFields().then((v) => saveMut.mutate(v))}>
              保存
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="显示名称" rules={[{ required: true }]}>
            <Input placeholder="博客" />
          </Form.Item>
          <Form.Item name="key" label="唯一 Key" rules={[{ required: true }]}>
            <Input placeholder="blog" />
          </Form.Item>
          <Form.Item name="parentId" label="父级导航">
            <Select
              allowClear
              placeholder="顶级导航请留空"
              options={(data ?? [])?.map((item) => ({ value: item.id, label: item.name }))}
            />
          </Form.Item>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="type" label="菜单类型">
                <Select options={[{ value: 'link', label: '链接' }, { value: 'menu', label: '菜单' }, { value: 'dropdown', label: '下拉菜单' }]} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="targetType" label="目标类型">
                <Select
                  options={[
                    { value: 'route', label: '路由' },
                    { value: 'category', label: '分类' },
                    { value: 'article', label: '文章' },
                    { value: 'external', label: '外部链接' },
                  ]}
                />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="targetValue" label="目标值">
            <Input placeholder="/blog 或分类 slug 等" />
          </Form.Item>
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="orderNum" label="排序数">
                <InputNumber style={{ width: '100%' }} min={0} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="enabled" label="启用" valuePropName="checked">
                <Switch />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// Page entry
// ─────────────────────────────────────────────

export function SitePage(): ReactElement {
  return (
    <section>
      <Typography.Title level={3} style={{ marginBottom: 16 }}>
        <Space><GlobalOutlined />站点设置</Space>
      </Typography.Title>
      <Tabs
        items={[
          { key: 'settings', label: '基础设置', children: <SettingsTab /> },
          { key: 'comments', label: '评论设置', children: <CommentPolicyTab /> },
          { key: 'social', label: '社交链接', children: <SocialLinksTab /> },
          { key: 'footer', label: '页脚配置', children: <FooterItemsTab /> },
          { key: 'nav', label: '导航菜单', children: <NavItemsTab /> },
        ]}
      />
    </section>
  );
}
