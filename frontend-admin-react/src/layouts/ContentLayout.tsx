/**
 * File: ContentLayout.tsx
 * Purpose: Wrap content module routes with a mobile-only segment switcher for articles/moments.
 * Module: frontend-admin-react/layouts/content, presentation frame layer.
 * Related: ArticlesPage, MomentsPage, ArticleEditorPage, and AdminLayout.
 */

import { FileTextOutlined, TagsOutlined, ThunderboltOutlined } from '@ant-design/icons';
import type { ReactElement } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

// Routes on which the section tab bar should appear
const LIST_PATHS = ['/content/articles', '/content/moments', '/content/taxonomy'];

export function ContentLayout(): ReactElement {
  const location = useLocation();
  const navigate = useNavigate();

  // Only show the tab bar on list pages, not inside the editor
  const showTabs = LIST_PATHS.includes(location.pathname);
  const activeTab = location.pathname.startsWith('/content/moments')
    ? 'moments'
    : location.pathname.startsWith('/content/taxonomy')
      ? 'taxonomy'
      : 'articles';

  return (
    <>
      {showTabs && (
        <div className="content-mobile-tabs" role="tablist" aria-label="内容分区">
          <button
            role="tab"
            aria-selected={activeTab === 'articles'}
            className={`cmt-tab${activeTab === 'articles' ? ' is-active' : ''}`}
            onClick={() => navigate('/content/articles')}
          >
            <FileTextOutlined />
            <span>文章</span>
          </button>
          <button
            role="tab"
            aria-selected={activeTab === 'moments'}
            className={`cmt-tab${activeTab === 'moments' ? ' is-active' : ''}`}
            onClick={() => navigate('/content/moments')}
          >
            <ThunderboltOutlined />
            <span>瞬间</span>
          </button>
          <button
            role="tab"
            aria-selected={activeTab === 'taxonomy'}
            className={`cmt-tab${activeTab === 'taxonomy' ? ' is-active' : ''}`}
            onClick={() => navigate('/content/taxonomy')}
          >
            <TagsOutlined />
            <span>分类/标签</span>
          </button>
        </div>
      )}
      <Outlet />
    </>
  );
}
