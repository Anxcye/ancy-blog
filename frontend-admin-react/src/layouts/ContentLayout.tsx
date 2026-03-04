/**
 * File: ContentLayout.tsx
 * Purpose: Wrap content module routes with a mobile-only segment switcher for articles/moments.
 * Module: frontend-admin-react/layouts/content, presentation frame layer.
 * Related: ArticlesPage, MomentsPage, ArticleEditorPage, and AdminLayout.
 */

import { FileTextOutlined, ThunderboltOutlined } from '@ant-design/icons';
import type { ReactElement } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';

// Routes on which the section tab bar should appear
const LIST_PATHS = ['/content/articles', '/content/moments'];

export function ContentLayout(): ReactElement {
  const location = useLocation();
  const navigate = useNavigate();

  // Only show the tab bar on the two list pages, not inside the editor
  const showTabs = LIST_PATHS.includes(location.pathname);
  const isArticles = !location.pathname.startsWith('/content/moments');

  return (
    <>
      {showTabs && (
        <div className="content-mobile-tabs" role="tablist" aria-label="内容分区">
          <button
            role="tab"
            aria-selected={isArticles}
            className={`cmt-tab${isArticles ? ' is-active' : ''}`}
            onClick={() => navigate('/content/articles')}
          >
            <FileTextOutlined />
            <span>文章</span>
          </button>
          <button
            role="tab"
            aria-selected={!isArticles}
            className={`cmt-tab${!isArticles ? ' is-active' : ''}`}
            onClick={() => navigate('/content/moments')}
          >
            <ThunderboltOutlined />
            <span>瞬间</span>
          </button>
        </div>
      )}
      <Outlet />
    </>
  );
}
