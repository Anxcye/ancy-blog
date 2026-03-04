/**
 * File: ArticleEditorPage.tsx
 * Purpose: Provide article create and edit form with TipTap rich-text editor.
 * Module: frontend-admin-react/pages/content, presentation layer.
 * Related: articles API module, article types, SimpleEditor, AdminLayout, and ArticlesPage.
 */

import { ArrowLeftOutlined, PlusOutlined, RobotOutlined, UploadOutlined } from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Card,
  Col,
  DatePicker,
  Divider,
  Form,
  Input,
  Row,
  Select,
  Space,
  Switch,
  Typography,
  Upload,
  message,
} from 'antd';
import type { UploadChangeParam } from 'antd/es/upload';

import { SimpleEditor } from '../../components/tiptap-templates/simple/simple-editor';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import dayjs from 'dayjs';

import {
  createArticle,
  createCategory,
  createTag,
  generateAiSlug,
  generateAiSummary,
  getArticle,
  listCategories,
  listTags,
  updateArticle,
} from '../../api/articles';
import { extractTextFromTiptapJson } from '../../lib/tiptap-utils';
import { useAuthStore } from '../../store/auth';
import type { ArticleFormValues } from '../../types/article';

// Convert a Chinese/English title to a basic URL slug (local fallback)
function titleToSlug(title: string): string {
  return title
    .toLowerCase()
    .replace(/[\s_]+/g, '-')
    .replace(/[^\w-]+/g, '')
    .replace(/--+/g, '-')
    .replace(/^-+|-+$/, '');
}

export function ArticleEditorPage(): ReactElement {
  const { id } = useParams<{ id: string }>();
  const isNew = !id;
  const [form] = Form.useForm<ArticleFormValues>();
  const [messageApi, ctx] = message.useMessage();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const accessToken = useAuthStore((s) => s.accessToken);
  const [slugLoading, setSlugLoading] = useState(false);
  const [summaryLoading, setSummaryLoading] = useState(false);
  const [coverUploading, setCoverUploading] = useState(false);
  const [newCatName, setNewCatName] = useState('');
  const [newCatSlug, setNewCatSlug] = useState('');
  const [newTagName, setNewTagName] = useState('');
  const [newTagSlug, setNewTagSlug] = useState('');
  const watchStatus = Form.useWatch('status', form);

  // Load existing article data in edit mode
  const { data: article, isLoading: articleLoading } = useQuery({
    queryKey: ['article', id],
    queryFn: () => getArticle(id!),
    enabled: !isNew,
  });

  // Load categories and tags for selectors
  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: listCategories,
  });
  const { data: tags = [] } = useQuery({
    queryKey: ['tags'],
    queryFn: listTags,
  });

  useEffect(() => {
    if (article) {
      form.setFieldsValue({
        title: article.title,
        slug: article.slug,
        summary: article.summary,
        content: article.content,
        contentKind: article.contentKind,
        status: article.status,
        visibility: article.visibility,
        allowComment: article.allowComment,
        coverImage: article.coverImage,
        originType: article.originType,
        sourceUrl: article.sourceUrl,
        categorySlug: article.categorySlug,
        tagSlugs: article.tagSlugs ?? [],
        isPinned: article.isPinned ?? false,
        isFeatured: article.isFeatured ?? false,
        publishedAt: article.publishedAt,
      });
    }
  }, [article, form]);

  const saveMut = useMutation({
    mutationFn: (values: ArticleFormValues) =>
      isNew ? createArticle(values) : updateArticle(id!, values),
    onSuccess: () => {
      messageApi.success(isNew ? '文章已创建' : '文章已保存');
      queryClient.invalidateQueries({ queryKey: ['articles'] });
      if (isNew) {
        navigate('/content/articles');
      }
    },
    onError: () => messageApi.error('保存失败，请重试'),
  });

  // Inline category create from editor
  const createCatMut = useMutation({
    mutationFn: ({ name, slug }: { name: string; slug: string }) => createCategory({ name, slug }),
    onSuccess: (cat) => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      form.setFieldValue('categorySlug', cat.slug);
      setNewCatName('');
      setNewCatSlug('');
      messageApi.success('分类已创建');
    },
    onError: () => messageApi.error('分类创建失败'),
  });

  // Inline tag create from editor
  const createTagMut = useMutation({
    mutationFn: ({ name, slug }: { name: string; slug: string }) => createTag({ name, slug }),
    onSuccess: (tag) => {
      queryClient.invalidateQueries({ queryKey: ['tags'] });
      const current: string[] = form.getFieldValue('tagSlugs') ?? [];
      form.setFieldValue('tagSlugs', [...current, tag.slug]);
      setNewTagName('');
      setNewTagSlug('');
      messageApi.success('标签已创建');
    },
    onError: () => messageApi.error('标签创建失败'),
  });

  const handleGenSlug = useCallback(async () => {
    const title: string = form.getFieldValue('title') ?? '';
    if (!title.trim()) {
      messageApi.warning('请先填写标题');
      return;
    }
    setSlugLoading(true);
    try {
      const result = await generateAiSlug(title);
      form.setFieldValue('slug', result.slug);
      if (result.fallbackUsed) {
        messageApi.info('AI 不可用，已使用本地转换');
      }
    } catch {
      form.setFieldValue('slug', titleToSlug(title));
      messageApi.info('AI 不可用，已使用本地转换');
    } finally {
      setSlugLoading(false);
    }
  }, [form, messageApi]);

  const handleGenSummary = useCallback(async () => {
    const title: string = form.getFieldValue('title') ?? '';
    const content: string = extractTextFromTiptapJson(form.getFieldValue('content') ?? '');
    if (!title.trim() && !content.trim()) {
      messageApi.warning('请先填写标题或正文');
      return;
    }
    setSummaryLoading(true);
    try {
      const result = await generateAiSummary(title, content);
      form.setFieldValue('summary', result.summary);
      if (result.fallbackUsed) {
        messageApi.info('AI 不可用，已使用截断摘要');
      }
    } catch {
      messageApi.error('摘要生成失败');
    } finally {
      setSummaryLoading(false);
    }
  }, [form, messageApi]);

  // Handle cover image upload via the image upload API
  const handleCoverUploadChange = useCallback(
    (info: UploadChangeParam) => {
      if (info.file.status === 'uploading') {
        setCoverUploading(true);
        return;
      }
      setCoverUploading(false);
      if (info.file.status === 'done') {
        const url: string = info.file.response?.data?.url ?? '';
        if (url) {
          form.setFieldValue('coverImage', url);
          messageApi.success('封面上传成功');
        } else {
          messageApi.error('上传成功但未获取到 URL');
        }
      } else if (info.file.status === 'error') {
        messageApi.error('封面上传失败');
      }
    },
    [form, messageApi],
  );

  // Serialize form values before submit: convert dayjs publishedAt to ISO string
  const handleFinish = useCallback(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (values: any) => {
      const payload = { ...values } as ArticleFormValues;
      if (values.publishedAt && typeof values.publishedAt.toISOString === 'function') {
        payload.publishedAt = (values.publishedAt as dayjs.Dayjs).toISOString();
      }
      saveMut.mutate(payload);
    },
    [saveMut],
  );

  return (
    <section>
      {ctx}

      {/* Header */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: 20,
        }}
      >
        <Space>
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/content/articles')}
          />
          <Typography.Title level={3} style={{ margin: 0 }}>
            {isNew ? '写文章' : '编辑文章'}
          </Typography.Title>
        </Space>
        <Space>
          <Button onClick={() => navigate('/content/articles')}>取消</Button>
          <Button
            type="primary"
            loading={saveMut.isPending}
            onClick={() => form.submit()}
          >
            {isNew ? '发布' : '保存修改'}
          </Button>
        </Space>
      </div>

      <Form
        form={form}
        layout="vertical"
        onFinish={handleFinish}
        initialValues={{
          contentKind: 'post',
          status: 'draft',
          visibility: 'public',
          allowComment: true,
          isPinned: false,
          isFeatured: false,
          originType: 'original',
          tagSlugs: [],
        }}
        disabled={articleLoading}
      >
        <Row gutter={[20, 0]}>
          {/* Left column — main content */}
          <Col xs={24} lg={16}>
            <Form.Item
              name="title"
              label="标题"
              rules={[{ required: true, message: '请填写标题' }]}
            >
              <Input placeholder="文章标题" size="large" />
            </Form.Item>

            <Row gutter={12}>
              <Col flex={1}>
                <Form.Item
                  name="slug"
                  label="Slug"
                  rules={[
                    { required: true, message: '请填写 Slug' },
                    {
                      pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/,
                      message: '只允许小写字母、数字和连字符',
                    },
                  ]}
                >
                  <Input placeholder="url-friendly-slug" />
                </Form.Item>
              </Col>
              <Col style={{ paddingTop: 30 }}>
                <Button
                  icon={<RobotOutlined />}
                  loading={slugLoading}
                  onClick={handleGenSlug}
                  size="small"
                >
                  AI 生成
                </Button>
              </Col>
            </Row>

            <Form.Item
              name="summary"
              label={
                <Space>
                  摘要
                  <Button
                    size="small"
                    type="link"
                    icon={<RobotOutlined />}
                    loading={summaryLoading}
                    style={{ padding: 0, height: 'auto' }}
                    onClick={handleGenSummary}
                  >
                    AI 生成
                  </Button>
                </Space>
              }
            >
              <Input.TextArea rows={3} placeholder="文章摘要（可选）" />
            </Form.Item>

            <Form.Item name="content" label="正文">
              <SimpleEditor />
            </Form.Item>
          </Col>

          {/* Right column — metadata & publish settings */}
          <Col xs={24} lg={8}>
            <Card size="small" title="发布设置" style={{ marginBottom: 16 }}>
              <Form.Item name="status" label="状态" style={{ marginBottom: 12 }}>
                <Select
                  options={[
                    { value: 'draft', label: '草稿' },
                    { value: 'published', label: '立即发布' },
                    { value: 'scheduled', label: '定时发布' },
                    { value: 'archived', label: '归档' },
                  ]}
                />
              </Form.Item>

              {/* Scheduled publish time — only shown when status = scheduled */}
              {watchStatus === 'scheduled' && (
                <Form.Item
                  name="publishedAt"
                  label="定时发布时间"
                  rules={[{ required: true, message: '请选择发布时间' }]}
                  style={{ marginBottom: 12 }}
                >
                  <DatePicker
                    showTime
                    format="YYYY-MM-DD HH:mm"
                    style={{ width: '100%' }}
                    disabledDate={(d) => d.isBefore(dayjs(), 'day')}
                  />
                </Form.Item>
              )}

              <Form.Item name="contentKind" label="内容类型" style={{ marginBottom: 12 }}>
                <Select
                  options={[
                    { value: 'post', label: '博客文章 (post)' },
                    { value: 'page', label: '独立页面 (page)' },
                  ]}
                />
              </Form.Item>

              <Form.Item name="visibility" label="可见性" style={{ marginBottom: 12 }}>
                <Select
                  options={[
                    { value: 'public', label: '公开' },
                    { value: 'unlisted', label: '不列出' },
                    { value: 'private', label: '仅自己可见' },
                  ]}
                />
              </Form.Item>

              <Form.Item
                name="allowComment"
                label="允许评论"
                valuePropName="checked"
                style={{ marginBottom: 8 }}
              >
                <Switch />
              </Form.Item>
              <Form.Item
                name="isPinned"
                label="置顶"
                valuePropName="checked"
                style={{ marginBottom: 8 }}
              >
                <Switch />
              </Form.Item>
              <Form.Item
                name="isFeatured"
                label="精选"
                valuePropName="checked"
                style={{ marginBottom: 0 }}
              >
                <Switch />
              </Form.Item>
            </Card>

            <Card size="small" title="分类与标签" style={{ marginBottom: 16 }}>
              <Form.Item name="categorySlug" label="分类" style={{ marginBottom: 12 }}>
                <Select
                  allowClear
                  placeholder="选择分类（可选）"
                  options={categories.map((c) => ({ value: c.slug, label: c.name }))}
                  dropdownRender={(menu) => (
                    <>
                      {menu}
                      <Divider style={{ margin: '4px 0' }} />
                      <div style={{ padding: '4px 8px', display: 'flex', flexDirection: 'column', gap: 4 }}>
                        <Input
                          size="small"
                          placeholder="分类名称"
                          value={newCatName}
                          onChange={(e) => setNewCatName(e.target.value)}
                          onKeyDown={(e) => e.stopPropagation()}
                        />
                        <Input
                          size="small"
                          placeholder="slug（英文）"
                          value={newCatSlug}
                          onChange={(e) => setNewCatSlug(e.target.value)}
                          onKeyDown={(e) => e.stopPropagation()}
                        />
                        <Button
                          type="primary"
                          size="small"
                          icon={<PlusOutlined />}
                          loading={createCatMut.isPending}
                          disabled={!newCatName.trim() || !newCatSlug.trim()}
                          onClick={() => createCatMut.mutate({ name: newCatName.trim(), slug: newCatSlug.trim() })}
                        >
                          添加分类
                        </Button>
                      </div>
                    </>
                  )}
                />
              </Form.Item>

              <Form.Item name="tagSlugs" label="标签" style={{ marginBottom: 0 }}>
                <Select
                  mode="multiple"
                  allowClear
                  placeholder="选择标签（可多选）"
                  options={tags.map((t) => ({ value: t.slug, label: t.name }))}
                  dropdownRender={(menu) => (
                    <>
                      {menu}
                      <Divider style={{ margin: '4px 0' }} />
                      <div style={{ padding: '4px 8px', display: 'flex', flexDirection: 'column', gap: 4 }}>
                        <Input
                          size="small"
                          placeholder="标签名称"
                          value={newTagName}
                          onChange={(e) => setNewTagName(e.target.value)}
                          onKeyDown={(e) => e.stopPropagation()}
                        />
                        <Input
                          size="small"
                          placeholder="slug（英文）"
                          value={newTagSlug}
                          onChange={(e) => setNewTagSlug(e.target.value)}
                          onKeyDown={(e) => e.stopPropagation()}
                        />
                        <Button
                          type="primary"
                          size="small"
                          icon={<PlusOutlined />}
                          loading={createTagMut.isPending}
                          disabled={!newTagName.trim() || !newTagSlug.trim()}
                          onClick={() => createTagMut.mutate({ name: newTagName.trim(), slug: newTagSlug.trim() })}
                        >
                          添加标签
                        </Button>
                      </div>
                    </>
                  )}
                />
              </Form.Item>
            </Card>

            <Card size="small" title="内容来源" style={{ marginBottom: 16 }}>
              <Form.Item name="originType" label="来源类型" style={{ marginBottom: 12 }}>
                <Select
                  options={[
                    { value: 'original', label: '原创' },
                    { value: 'repost', label: '转载' },
                    { value: 'translation', label: '翻译' },
                  ]}
                />
              </Form.Item>

              <Form.Item name="sourceUrl" label="原文链接" style={{ marginBottom: 0 }}>
                <Input placeholder="https://..." />
              </Form.Item>
            </Card>

            <Card size="small" title="封面与元数据">
              <Form.Item name="coverImage" label="封面图" style={{ marginBottom: 8 }}>
                <Input placeholder="https://cdn.example.com/cover.jpg" />
              </Form.Item>
              <Upload
                name="file"
                accept="image/*"
                showUploadList={false}
                action="/api/v1/admin/upload/image"
                headers={{ Authorization: `Bearer ${accessToken}` }}
                onChange={handleCoverUploadChange}
              >
                <Button
                  icon={<UploadOutlined />}
                  loading={coverUploading}
                  size="small"
                  style={{ marginBottom: 0 }}
                >
                  上传封面图
                </Button>
              </Upload>
            </Card>
          </Col>
        </Row>
      </Form>
    </section>
  );
}
