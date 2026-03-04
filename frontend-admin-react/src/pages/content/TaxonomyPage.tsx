/**
 * File: TaxonomyPage.tsx
 * Purpose: Manage categories and tags with create/delete and inline slug generation.
 * Module: frontend-admin-react/pages/content, presentation layer.
 * Related: articles API module (category/tag CRUD), ContentLayout.
 */

import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Button, Card, Col, Input, Popconfirm, Row, Space, Table, Tabs, Typography, message } from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';

import {
  createCategory,
  createTag,
  deleteCategory,
  deleteTag,
  listCategories,
  listTags,
} from '../../api/articles';
import type { Category, Tag } from '../../api/articles';

function toSlug(name: string): string {
  return name
    .toLowerCase()
    .trim()
    .replace(/[\s_]+/g, '-')
    .replace(/[^\w-]/g, '')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');
}

function CategoriesTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const [slugEdited, setSlugEdited] = useState(false);

  const { data: categories = [], isLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: listCategories,
  });

  const addMut = useMutation({
    mutationFn: createCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      setName('');
      setSlug('');
      setSlugEdited(false);
      messageApi.success('分类已创建');
    },
    onError: () => messageApi.error('创建失败'),
  });

  const delMut = useMutation({
    mutationFn: deleteCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      messageApi.success('分类已删除');
    },
    onError: () => messageApi.error('删除失败'),
  });

  function handleNameChange(val: string) {
    setName(val);
    if (!slugEdited) setSlug(toSlug(val));
  }

  function handleAdd() {
    const s = slug || toSlug(name);
    if (!name.trim() || !s) return;
    addMut.mutate({ name: name.trim(), slug: s });
  }

  return (
    <>
      {ctx}
      <Space.Compact style={{ marginBottom: 16, width: '100%' }}>
        <Input
          placeholder="分类名称"
          value={name}
          onChange={(e) => handleNameChange(e.target.value)}
          onPressEnter={handleAdd}
          style={{ maxWidth: 200 }}
        />
        <Input
          placeholder="slug"
          value={slug}
          onChange={(e) => { setSlug(e.target.value); setSlugEdited(true); }}
          onPressEnter={handleAdd}
          style={{ maxWidth: 200 }}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd} loading={addMut.isPending}>
          添加
        </Button>
      </Space.Compact>
      <Table<Category>
        rowKey="id"
        dataSource={categories}
        loading={isLoading}
        pagination={false}
        size="small"
        columns={[
          { title: '名称', dataIndex: 'name', key: 'name' },
          { title: 'Slug', dataIndex: 'slug', key: 'slug' },
          {
            title: '操作',
            key: 'action',
            width: 80,
            render: (_, record) => (
              <Popconfirm title="确认删除？" onConfirm={() => delMut.mutate(record.id)}>
                <Button type="text" danger icon={<DeleteOutlined />} size="small" />
              </Popconfirm>
            ),
          },
        ]}
      />
    </>
  );
}

function TagsTab(): ReactElement {
  const [messageApi, ctx] = message.useMessage();
  const queryClient = useQueryClient();
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const [slugEdited, setSlugEdited] = useState(false);

  const { data: tags = [], isLoading } = useQuery({
    queryKey: ['tags'],
    queryFn: listTags,
  });

  const addMut = useMutation({
    mutationFn: createTag,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] });
      setName('');
      setSlug('');
      setSlugEdited(false);
      messageApi.success('标签已创建');
    },
    onError: () => messageApi.error('创建失败'),
  });

  const delMut = useMutation({
    mutationFn: deleteTag,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] });
      messageApi.success('标签已删除');
    },
    onError: () => messageApi.error('删除失败'),
  });

  function handleNameChange(val: string) {
    setName(val);
    if (!slugEdited) setSlug(toSlug(val));
  }

  function handleAdd() {
    const s = slug || toSlug(name);
    if (!name.trim() || !s) return;
    addMut.mutate({ name: name.trim(), slug: s });
  }

  return (
    <>
      {ctx}
      <Space.Compact style={{ marginBottom: 16, width: '100%' }}>
        <Input
          placeholder="标签名称"
          value={name}
          onChange={(e) => handleNameChange(e.target.value)}
          onPressEnter={handleAdd}
          style={{ maxWidth: 200 }}
        />
        <Input
          placeholder="slug"
          value={slug}
          onChange={(e) => { setSlug(e.target.value); setSlugEdited(true); }}
          onPressEnter={handleAdd}
          style={{ maxWidth: 200 }}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd} loading={addMut.isPending}>
          添加
        </Button>
      </Space.Compact>
      <Table<Tag>
        rowKey="id"
        dataSource={tags}
        loading={isLoading}
        pagination={false}
        size="small"
        columns={[
          { title: '名称', dataIndex: 'name', key: 'name' },
          { title: 'Slug', dataIndex: 'slug', key: 'slug' },
          {
            title: '操作',
            key: 'action',
            width: 80,
            render: (_, record) => (
              <Popconfirm title="确认删除？" onConfirm={() => delMut.mutate(record.id)}>
                <Button type="text" danger icon={<DeleteOutlined />} size="small" />
              </Popconfirm>
            ),
          },
        ]}
      />
    </>
  );
}

export function TaxonomyPage(): ReactElement {
  return (
    <Row justify="center">
      <Col xs={24} md={20} lg={16}>
        <Typography.Title level={4} style={{ marginBottom: 16 }}>分类与标签</Typography.Title>
        <Card>
          <Tabs
            defaultActiveKey="categories"
            items={[
              { key: 'categories', label: '分类', children: <CategoriesTab /> },
              { key: 'tags', label: '标签', children: <TagsTab /> },
            ]}
          />
        </Card>
      </Col>
    </Row>
  );
}
