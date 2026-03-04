/**
 * File: SystemPage.tsx
 * Purpose: Provide integration provider management, translation job queue, and content override.
 * Module: frontend-admin-react/pages/system, presentation layer.
 * Related: system API module, system types, and AdminLayout.
 */

import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  CloudOutlined,
  PlusOutlined,
  ReloadOutlined,
  RobotOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Card,
  Col,
  Drawer,
  Form,
  Input,
  InputNumber,
  Row,
  Select,
  Space,
  Switch,
  Table,
  Tabs,
  Tag,
  Tooltip,
  Typography,
  message,
} from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';

import {
  createTranslationJob,
  listProviders,
  listTranslationContents,
  listTranslationJobs,
  retryTranslationJob,
  testProvider,
  updateProvider,
  updateTranslationContent,
} from '../../api/system';
import type {
  CreateTranslationJobPayload,
  IntegrationProvider,
  TranslationContent,
  TranslationContentListParams,
  TranslationContentStatus,
  TranslationJob,
  TranslationJobListParams,
  TranslationJobStatus,
  TranslationSourceType,
  UpdateTranslationContentPayload,
} from '../../types/system';

function fmtDate(iso?: string): string {
  if (!iso) return '—';
  const d = new Date(iso);
  const p = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`;
}

// Sentinel value the backend uses for masked secrets
const MASKED = '***';

// ─────────────────────────────────────────────
// Integration tab
// ─────────────────────────────────────────────

const PROVIDER_META: Record<string, { icon: ReactElement; label: string }> = {
  cloudflare_r2: { icon: <CloudOutlined />, label: 'Cloudflare R2' },
  openai_compatible: { icon: <RobotOutlined />, label: 'OpenAI Compatible' },
};

// R2-specific config form fields
function R2ConfigFields(): ReactElement {
  return (
    <>
      <Form.Item name="endpoint" label="Endpoint URL" rules={[{ required: true }]}>
        <Input placeholder="https://xxxx.r2.cloudflarestorage.com" />
      </Form.Item>
      <Form.Item name="bucket" label="Bucket 名称" rules={[{ required: true }]}>
        <Input placeholder="my-bucket" />
      </Form.Item>
      <Form.Item name="accessKeyId" label="Access Key ID" rules={[{ required: true }]}>
        <Input placeholder="your-access-key-id" />
      </Form.Item>
      <Form.Item name="secretAccessKey" label="Secret Access Key">
        <Input.Password placeholder="已配置，输入新值以更新" />
      </Form.Item>
      <Form.Item name="publicUrlBase" label="Public URL Base">
        <Input placeholder="https://cdn.example.com" />
      </Form.Item>
    </>
  );
}

// OpenAI-compatible config form fields
function OpenAIConfigFields(): ReactElement {
  return (
    <>
      <Form.Item name="apiBase" label="API Base URL" rules={[{ required: true }]}>
        <Input placeholder="https://api.openai.com/v1" />
      </Form.Item>
      <Form.Item name="apiKey" label="API Key">
        <Input.Password placeholder="已配置，输入新值以更新" />
      </Form.Item>
      <Form.Item name="defaultModel" label="默认模型">
        <Input placeholder="gpt-4.1-mini" />
      </Form.Item>
    </>
  );
}

function IntegrationsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingProvider, setEditingProvider] = useState<IntegrationProvider | null>(null);
  const [saving, setSaving] = useState(false);
  const [testingKey, setTestingKey] = useState<string | null>(null);

  const { data: providers = [], isLoading } = useQuery({
    queryKey: ['providers'],
    queryFn: () => listProviders(),
  });

  function openConfig(provider: IntegrationProvider): void {
    setEditingProvider(provider);
    // Pre-fill form: keep masked values as-is so backend can treat "***" as "no change"
    form.setFieldsValue({ ...provider.configJson, enabled: provider.enabled });
    setDrawerOpen(true);
  }

  async function handleSaveAndTest(): Promise<void> {
    try {
      const vals = await form.validateFields();
      const { enabled, ...configFields } = vals;
      setSaving(true);
      await updateProvider(editingProvider!.providerKey, {
        enabled: enabled as boolean,
        configJson: configFields as Record<string, unknown>,
      });
      messageApi.success('配置已保存，正在测试连接…');
      setSaving(false);
      setDrawerOpen(false);
      queryClient.invalidateQueries({ queryKey: ['providers'] });

      // Auto-run connection test after save
      setTestingKey(editingProvider!.providerKey);
      try {
        const result = await testProvider(editingProvider!.providerKey);
        if (result.ok) {
          messageApi.success(`连接正常 (${result.latencyMs}ms)`);
        } else {
          messageApi.warning(`连接测试失败: ${result.message}`);
        }
      } finally {
        setTestingKey(null);
        queryClient.invalidateQueries({ queryKey: ['providers'] });
      }
    } catch {
      setSaving(false);
      messageApi.error('保存失败，请检查表单');
    }
  }

  async function handleTest(providerKey: string): Promise<void> {
    setTestingKey(providerKey);
    try {
      const result = await testProvider(providerKey);
      if (result.ok) {
        messageApi.success(`连接正常 (${result.latencyMs}ms)`);
      } else {
        messageApi.error(`连接失败: ${result.message}`);
      }
      queryClient.invalidateQueries({ queryKey: ['providers'] });
    } catch {
      messageApi.error('测试请求失败');
    } finally {
      setTestingKey(null);
    }
  }

  return (
    <>
      {ctx}
      {isLoading ? (
        <Row gutter={[16, 16]}>
          {[1, 2].map((i) => (
            <Col xs={24} md={12} key={i}>
              <Card loading style={{ borderRadius: 14 }} />
            </Col>
          ))}
        </Row>
      ) : (
        <Row gutter={[16, 16]}>
          {providers.map((provider) => {
            const meta = PROVIDER_META[provider.providerKey];
            const lastTest = provider.metaJson;
            const isTesting = testingKey === provider.providerKey;

            return (
              <Col xs={24} md={12} key={provider.providerKey}>
                <Card
                  style={{ borderRadius: 14 }}
                  styles={{ body: { padding: '20px 24px' } }}
                >
                  {/* Card header */}
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: 14 }}>
                    <Space size={10}>
                      <div
                        style={{
                          width: 40,
                          height: 40,
                          borderRadius: 10,
                          background: provider.providerType === 'llm'
                            ? 'rgba(124,58,237,0.10)'
                            : 'rgba(31,143,138,0.10)',
                          color: provider.providerType === 'llm' ? '#7c3aed' : '#1f8f8a',
                          display: 'grid',
                          placeItems: 'center',
                          fontSize: 18,
                          flexShrink: 0,
                        }}
                      >
                        {meta?.icon}
                      </div>
                      <div>
                        <Typography.Text strong style={{ fontSize: 15 }}>
                          {meta?.label ?? provider.name}
                        </Typography.Text>
                        <br />
                        <Tag
                          color={provider.enabled ? 'success' : 'default'}
                          style={{ marginTop: 2 }}
                        >
                          {provider.enabled ? '已启用' : '未启用'}
                        </Tag>
                      </div>
                    </Space>
                    <Space size={6}>
                      <Button
                        size="small"
                        icon={<ReloadOutlined />}
                        loading={isTesting}
                        onClick={() => handleTest(provider.providerKey)}
                      >
                        测试
                      </Button>
                      <Button size="small" type="primary" onClick={() => openConfig(provider)}>
                        配置
                      </Button>
                    </Space>
                  </div>

                  {/* Last test result */}
                  {lastTest?.lastTestAt ? (
                    <Space size={8}>
                      {lastTest.lastTestOk ? (
                        <CheckCircleOutlined style={{ color: '#1f8f8a' }} />
                      ) : (
                        <CloseCircleOutlined style={{ color: '#ef4444' }} />
                      )}
                      <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                        {lastTest.lastTestOk
                          ? `连接正常 · ${lastTest.latencyMs ?? '—'}ms`
                          : `连接失败 · ${lastTest.lastTestMsg ?? ''}`}
                      </Typography.Text>
                      <Typography.Text type="secondary" style={{ fontSize: 11 }}>
                        {fmtDate(lastTest.lastTestAt)}
                      </Typography.Text>
                    </Space>
                  ) : (
                    <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                      尚未测试连接
                    </Typography.Text>
                  )}
                </Card>
              </Col>
            );
          })}
        </Row>
      )}

      {/* Provider config drawer */}
      <Drawer
        title={`配置 ${PROVIDER_META[editingProvider?.providerKey ?? '']?.label ?? ''}`}
        width={440}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button type="primary" loading={saving} onClick={handleSaveAndTest}>
              保存并测试
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          {editingProvider?.providerKey === 'cloudflare_r2' && <R2ConfigFields />}
          {editingProvider?.providerKey === 'openai_compatible' && <OpenAIConfigFields />}
          <Form.Item
            name="enabled"
            label="启用此集成"
            valuePropName="checked"
            style={{ marginTop: 8 }}
          >
            <Switch />
          </Form.Item>
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            密钥字段留空表示不修改现有配置。
          </Typography.Text>
        </Form>
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// Translation jobs tab
// ─────────────────────────────────────────────

const JOB_STATUS_COLOR: Record<string, string> = {
  queued: 'default',
  running: 'processing',
  succeeded: 'success',
  failed: 'error',
};
const JOB_STATUS_LABEL: Record<string, string> = {
  queued: '排队中',
  running: '运行中',
  succeeded: '已完成',
  failed: '失败',
};

function TranslationJobsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<CreateTranslationJobPayload>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [params, setParams] = useState<TranslationJobListParams>({ page: 1, pageSize: 20 });

  const { data, isLoading } = useQuery({
    queryKey: ['translation-jobs', params],
    queryFn: () => listTranslationJobs(params),
    refetchInterval: 8000, // auto-refresh for running jobs
  });

  const createMut = useMutation({
    mutationFn: createTranslationJob,
    onSuccess: () => {
      messageApi.success('翻译任务已创建');
      setDrawerOpen(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['translation-jobs'] });
    },
    onError: () => messageApi.error('创建失败'),
  });

  const retryMut = useMutation({
    mutationFn: retryTranslationJob,
    onSuccess: () => {
      messageApi.success('任务已重新加入队列');
      queryClient.invalidateQueries({ queryKey: ['translation-jobs'] });
    },
    onError: () => messageApi.error('重试失败'),
  });

  const columns = [
    {
      title: '类型',
      dataIndex: 'sourceType',
      key: 'sourceType',
      width: 70,
      render: (t: TranslationSourceType) => (
        <Tag color={t === 'article' ? 'blue' : 'purple'}>{t === 'article' ? '文章' : '瞬间'}</Tag>
      ),
    },
    {
      title: 'Source ID',
      dataIndex: 'sourceId',
      key: 'sourceId',
      render: (id: string) => (
        <Typography.Text code style={{ fontSize: 11 }}>
          {id.slice(0, 8)}…
        </Typography.Text>
      ),
    },
    {
      title: '语言',
      key: 'locales',
      width: 110,
      render: (_: unknown, r: TranslationJob) => (
        <Typography.Text style={{ fontSize: 12 }}>
          {r.sourceLocale} → {r.targetLocale}
        </Typography.Text>
      ),
    },
    {
      title: '模型',
      dataIndex: 'modelName',
      key: 'modelName',
      width: 130,
      render: (m: string) => (
        <Typography.Text style={{ fontSize: 12 }}>{m}</Typography.Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (s: TranslationJobStatus) => (
        <Tag
          color={JOB_STATUS_COLOR[s]}
          icon={s === 'running' ? <SyncOutlined spin /> : undefined}
        >
          {JOB_STATUS_LABEL[s] ?? s}
        </Tag>
      ),
    },
    {
      title: '重试',
      key: 'retry',
      width: 80,
      render: (_: unknown, r: TranslationJob) => (
        <Typography.Text type="secondary" style={{ fontSize: 12 }}>
          {r.retryCount}/{r.maxRetries}
        </Typography.Text>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 140,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      render: (_: unknown, r: TranslationJob) =>
        r.status === 'failed' ? (
          <Tooltip title="重新加入队列">
            <Button
              type="text"
              size="small"
              icon={<ReloadOutlined />}
              loading={retryMut.isPending}
              onClick={() => retryMut.mutate(r.id)}
            />
          </Tooltip>
        ) : null,
    },
  ];

  return (
    <>
      {ctx}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 12 }}>
        <Space>
          <Select
            style={{ width: 120 }}
            placeholder="状态"
            allowClear
            options={[
              { value: 'queued', label: '排队中' },
              { value: 'running', label: '运行中' },
              { value: 'succeeded', label: '已完成' },
              { value: 'failed', label: '失败' },
            ]}
            onChange={(v) =>
              setParams((p) => ({ ...p, status: (v as TranslationJobStatus) || '', page: 1 }))
            }
          />
          <Select
            style={{ width: 100 }}
            placeholder="类型"
            allowClear
            options={[
              { value: 'article', label: '文章' },
              { value: 'moment', label: '瞬间' },
            ]}
            onChange={(v) =>
              setParams((p) => ({ ...p, sourceType: (v as TranslationSourceType) || '', page: 1 }))
            }
          />
        </Space>
        <Button
          type="primary"
          size="small"
          icon={<PlusOutlined />}
          onClick={() => { form.resetFields(); setDrawerOpen(true); }}
        >
          创建任务
        </Button>
      </div>

      <Table
        rowKey="id"
        size="small"
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
        scroll={{ x: 720 }}
      />

      {/* Create job drawer */}
      <Drawer
        title="创建翻译任务"
        width={440}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button
              type="primary"
              loading={createMut.isPending}
              onClick={() => form.validateFields().then((v) => createMut.mutate(v))}
            >
              提交任务
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Row gutter={12}>
            <Col span={12}>
              <Form.Item name="sourceType" label="内容类型" rules={[{ required: true }]}>
                <Select options={[{ value: 'article', label: '文章' }, { value: 'moment', label: '瞬间' }]} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="sourceLocale" label="源语言" rules={[{ required: true }]}>
                <Select options={[{ value: 'zh-CN', label: '中文' }, { value: 'en-US', label: 'English' }]} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="sourceId" label="Source ID (UUID)" rules={[{ required: true }]}>
            <Input placeholder="文章或瞬间的 UUID" />
          </Form.Item>
          <Form.Item name="targetLocale" label="目标语言" rules={[{ required: true }]}>
            <Select options={[{ value: 'en-US', label: 'English (en-US)' }, { value: 'zh-CN', label: '中文 (zh-CN)' }, { value: 'ja-JP', label: '日本語 (ja-JP)' }]} />
          </Form.Item>
          <Row gutter={12}>
            <Col span={14}>
              <Form.Item name="providerKey" label="提供商" rules={[{ required: true }]}>
                <Select options={[{ value: 'openai_compatible', label: 'OpenAI Compatible' }]} />
              </Form.Item>
            </Col>
            <Col span={10}>
              <Form.Item name="maxRetries" label="最大重试">
                <InputNumber style={{ width: '100%' }} min={0} max={5} defaultValue={3} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="modelName" label="模型名称" rules={[{ required: true }]}>
            <Input placeholder="gpt-4.1-mini" />
          </Form.Item>
          <Form.Item name="autoPublish" label="完成后自动发布" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// Translation contents tab
// ─────────────────────────────────────────────

const CONTENT_STATUS_COLOR: Record<string, string> = {
  draft: 'default',
  published: 'success',
  archived: 'warning',
};
const CONTENT_STATUS_LABEL: Record<string, string> = {
  draft: '草稿',
  published: '已发布',
  archived: '已归档',
};

function TranslationContentsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [form] = Form.useForm<UpdateTranslationContentPayload>();
  const [params, setParams] = useState<TranslationContentListParams>({ sourceType: 'article', page: 1, pageSize: 20 });
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingContent, setEditingContent] = useState<TranslationContent | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['translation-contents', params],
    queryFn: () => listTranslationContents(params),
  });

  const saveMut = useMutation({
    mutationFn: (vals: UpdateTranslationContentPayload) => updateTranslationContent(vals),
    onSuccess: () => {
      messageApi.success('译文已保存');
      setDrawerOpen(false);
      queryClient.invalidateQueries({ queryKey: ['translation-contents'] });
    },
    onError: () => messageApi.error('保存失败'),
  });

  function openEdit(content: TranslationContent): void {
    setEditingContent(content);
    form.setFieldsValue(content);
    setDrawerOpen(true);
  }

  const columns = [
    {
      title: '类型',
      dataIndex: 'sourceType',
      key: 'sourceType',
      width: 70,
      render: (t: TranslationSourceType) => (
        <Tag color={t === 'article' ? 'blue' : 'purple'}>{t === 'article' ? '文章' : '瞬间'}</Tag>
      ),
    },
    {
      title: '语言',
      dataIndex: 'locale',
      key: 'locale',
      width: 90,
    },
    {
      title: '标题 / 内容预览',
      key: 'preview',
      render: (_: unknown, r: TranslationContent) => (
        <Typography.Paragraph ellipsis={{ rows: 1 }} style={{ margin: 0, fontSize: 13, maxWidth: 300 }}>
          {r.title ?? r.content}
        </Typography.Paragraph>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (s: TranslationContentStatus) => (
        <Tag color={CONTENT_STATUS_COLOR[s]}>{CONTENT_STATUS_LABEL[s] ?? s}</Tag>
      ),
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      key: 'updatedAt',
      width: 140,
      render: (v: string) => fmtDate(v),
    },
    {
      title: '操作',
      key: 'actions',
      width: 80,
      render: (_: unknown, r: TranslationContent) => (
        <Button type="link" size="small" style={{ padding: 0 }} onClick={() => openEdit(r)}>
          编辑
        </Button>
      ),
    },
  ];

  return (
    <>
      {ctx}
      <Row gutter={[12, 12]} style={{ marginBottom: 12 }}>
        <Col xs={12} sm={6} md={4}>
          <Select
            style={{ width: '100%' }}
            value={params.sourceType}
            options={[{ value: 'article', label: '文章' }, { value: 'moment', label: '瞬间' }]}
            onChange={(v) => setParams((p) => ({ ...p, sourceType: v as TranslationSourceType, page: 1 }))}
          />
        </Col>
        <Col xs={12} sm={6} md={4}>
          <Select
            style={{ width: '100%' }}
            placeholder="语言"
            allowClear
            options={[
              { value: 'en-US', label: 'en-US' },
              { value: 'zh-CN', label: 'zh-CN' },
              { value: 'ja-JP', label: 'ja-JP' },
            ]}
            onChange={(v) => setParams((p) => ({ ...p, locale: (v as string) || undefined, page: 1 }))}
          />
        </Col>
      </Row>

      <Table
        rowKey="id"
        size="small"
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
        scroll={{ x: 600 }}
      />

      {/* Content edit drawer */}
      <Drawer
        title="编辑译文"
        width={520}
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        extra={
          <Space>
            <Button onClick={() => setDrawerOpen(false)}>取消</Button>
            <Button
              type="primary"
              loading={saveMut.isPending}
              onClick={() =>
                form.validateFields().then((vals) =>
                  saveMut.mutate({
                    ...vals,
                    sourceType: editingContent!.sourceType,
                    sourceId: editingContent!.sourceId,
                    locale: editingContent!.locale,
                  }),
                )
              }
            >
              保存
            </Button>
          </Space>
        }
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          {editingContent?.sourceType === 'article' && (
            <>
              <Form.Item name="title" label="标题译文">
                <Input placeholder="Translated title" />
              </Form.Item>
              <Form.Item name="summary" label="摘要译文">
                <Input.TextArea rows={3} placeholder="Translated summary" />
              </Form.Item>
            </>
          )}
          <Form.Item name="content" label="正文译文" rules={[{ required: true }]}>
            <Input.TextArea
              rows={16}
              placeholder="Translated content…"
              style={{ fontFamily: 'monospace', fontSize: 13 }}
            />
          </Form.Item>
          <Form.Item name="status" label="发布状态">
            <Select
              options={[
                { value: 'draft', label: '草稿' },
                { value: 'published', label: '立即发布' },
                { value: 'archived', label: '归档' },
              ]}
            />
          </Form.Item>
        </Form>
      </Drawer>
    </>
  );
}

// ─────────────────────────────────────────────
// Page entry
// ─────────────────────────────────────────────

export function SystemPage(): ReactElement {
  return (
    <section>
      <Typography.Title level={3} style={{ marginBottom: 16 }}>
        系统设置
      </Typography.Title>
      <Tabs
        items={[
          {
            key: 'integrations',
            label: <Space size={6}><CloudOutlined />集成中心</Space>,
            children: <IntegrationsTab />,
          },
          {
            key: 'jobs',
            label: <Space size={6}><SyncOutlined />翻译任务</Space>,
            children: <TranslationJobsTab />,
          },
          {
            key: 'contents',
            label: <Space size={6}><RobotOutlined />译文管理</Space>,
            children: <TranslationContentsTab />,
          },
        ]}
      />
    </section>
  );
}

// Suppress unused import warning for masked sentinel
void MASKED;
