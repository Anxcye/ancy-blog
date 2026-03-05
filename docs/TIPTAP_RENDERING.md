# TipTap Content Rendering

## Overview

The blog frontend uses `@tiptap/html` to render TipTap JSON content from the backend without requiring the full editor.

## Architecture

### Packages
- `@tiptap/vue-3` - Vue 3 editor integration
- `@tiptap/core` - Base types and utilities
- `@tiptap/starter-kit` - Standard nodes (paragraph, heading, list, blockquote, code, etc.)
- `@tiptap/extension-text-align` - Text alignment (left/center/right/justify)
- `@tiptap/extension-highlight` - Background color highlighting
- `@tiptap/extension-task-list` + `@tiptap/extension-task-item` - Task lists with checkboxes
- `@tiptap/extension-underline` - Underline text
- `@tiptap/extension-subscript` + `@tiptap/extension-superscript` - Sub/superscript
- `@tiptap/extension-typography` - Smart typography replacements

### Custom Extensions
Located in `app/components/tiptap-extensions/`:

1. **x-post-embed.ts** - X (Twitter) post embeds
   - Attributes: `postId`, `author`
   - Renders as blockquote with link to tweet

2. **tmdb-card-embed.ts** - TMDB movie/TV card embeds
   - Attributes: `tmdbId`, `mediaType`, `title`
   - Renders as styled card with link to TMDB

## Usage

```vue
<TiptapRenderer :content="article.content" />
```

The `TiptapRenderer` component:
1. Parses JSON content
2. Calls `generateHTML()` with all extensions
3. Renders HTML with scoped styles

## Adding New Custom Nodes

1. Create extension in `app/components/tiptap-extensions/`
2. Define `name`, `addAttributes()`, `parseHTML()`, `renderHTML()`
3. Import and add to extensions array in `TiptapRenderer.vue`
4. Add CSS styles in `<style scoped>` section

## Sync with Admin

Custom nodes must match the schema in `frontend-admin-react/src/components/tiptap-node/`:
- Same `name` property
- Same `data-embed-type` attribute
- Compatible `parseHTML()` selector
