/**
 * File: ArticleEditorPage.tsx
 * Purpose: Provide article create and edit form with TipTap rich-text editor.
 * Module: frontend-admin-react/pages/content, presentation layer.
 * Related: articles API module, article types, SimpleEditor, AdminLayout, and ArticlesPage.
 */

import { ArrowLeftOutlined, RobotOutlined } from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Button,
  Card,
  Col,
  Form,
  Input,
  Row,
  Select,
  Space,
  Switch,
  Typography,
  message,
} from 'antd';

import { SimpleEditor } from '../../components/tiptap-templates/simple/simple-editor';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import {
  createArticle,
  generateAiSlug,
  generateAiSummary,
  getArticle,
  updateArticle,
} from '../../api/articles';
import { extractTextFromTiptapJson } from '../../lib/tiptap-utils';
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
  const [slugLoading, setSlugLoading] = useState(false);
  const [summaryLoading, setSummaryLoading] = useState(false);

  // Load existing article data in edit mode
  const { data: article, isLoading: articleLoading } = useQuery({
    queryKey: ['article', id],
    queryFn: () => getArticle(id!),
    enabled: !isNew,
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
      // Fall back to local conversion
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
        onFinish={(values) => saveMut.mutate(values)}
        initialValues={{
          contentKind: 'post',
          status: 'draft',
          visibility: 'public',
          allowComment: true,
          originType: 'original',
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
                style={{ marginBottom: 0 }}
              >
                <Switch />
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
              <Form.Item name="coverImage" label="封面图 URL" style={{ marginBottom: 0 }}>
                <Input placeholder="https://cdn.example.com/cover.jpg" />
              </Form.Item>
            </Card>
          </Col>
        </Row>
      </Form>
    </section>
  );
}
