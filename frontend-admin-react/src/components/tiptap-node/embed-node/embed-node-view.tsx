/**
 * File: embed-node-view.tsx
 * Purpose: Render in-editor preview cards for xPostEmbed and tmdbCardEmbed nodes.
 * Module: frontend-admin-react/components/tiptap-node/embed-node, editor view layer.
 * Related: XPostEmbed, TmdbCardEmbed, embed-node.scss.
 */

import { NodeViewWrapper } from "@tiptap/react"
import type { NodeViewProps } from "@tiptap/react"

import "@/components/tiptap-node/embed-node/embed-node.scss"

export function EmbedNodeView({ node, selected }: NodeViewProps) {
  if (node.type.name === "xPostEmbed") {
    const { postId, author } = node.attrs as { postId: string; author: string }
    const selectedClass = selected ? " is-selected" : ""

    return (
      <NodeViewWrapper>
        <div
          className={`embed-node-card embed-node-x${selectedClass}`}
          data-drag-handle=""
          contentEditable={false}
        >
          <span className="embed-node-icon" aria-hidden="true">𝕏</span>
          <div className="embed-node-body">
            <span className="embed-node-label">X Post</span>
            <span className="embed-node-primary">{postId || "(no Post ID)"}</span>
            {author && (
              <span className="embed-node-secondary">
                @{author.replace(/^@/, "")}
              </span>
            )}
          </div>
        </div>
      </NodeViewWrapper>
    )
  }

  if (node.type.name === "tmdbCardEmbed") {
    const { tmdbId, mediaType, title } = node.attrs as {
      tmdbId: string
      mediaType: string
      title: string
    }
    const selectedClass = selected ? " is-selected" : ""
    const typeLabel = (mediaType || "movie").toUpperCase()

    return (
      <NodeViewWrapper>
        <div
          className={`embed-node-card embed-node-tmdb${selectedClass}`}
          data-drag-handle=""
          contentEditable={false}
        >
          <span className="embed-node-icon" aria-hidden="true">🎬</span>
          <div className="embed-node-body">
            <span className="embed-node-label">TMDB · {typeLabel}</span>
            <span className="embed-node-primary">{title || "(no title)"}</span>
            {tmdbId && (
              <span className="embed-node-secondary">ID: {tmdbId}</span>
            )}
          </div>
        </div>
      </NodeViewWrapper>
    )
  }

  return null
}
