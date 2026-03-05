/**
 * File: embed-node-view.tsx
 * Purpose: Render in-editor preview cards for xPostEmbed and tmdbCardEmbed nodes.
 * Module: frontend-admin-react/components/tiptap-node/embed-node, editor view layer.
 * Related: XPostEmbed, TmdbCardEmbed, embed-node.scss.
 */

import { NodeViewWrapper } from "@tiptap/react"
import type { NodeViewProps } from "@tiptap/react"
import { useEffect, useRef } from "react"

import "@/components/tiptap-node/embed-node/embed-node.scss"

export function EmbedNodeView({ node, selected }: NodeViewProps) {
  const twitterRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (node.type.name === "xPostEmbed" && twitterRef.current) {
      const w = window as any
      if (w.twttr?.widgets) {
        w.twttr.widgets.load(twitterRef.current)
      } else if (!w.twttr) {
        const script = document.createElement('script')
        script.src = 'https://platform.twitter.com/widgets.js'
        script.async = true
        document.body.appendChild(script)
      }
    }
  }, [node])

  if (node.type.name === "xPostEmbed") {
    const { postId, author } = node.attrs as { postId: string; author: string }
    const selectedClass = selected ? " is-selected" : ""
    const url = `https://twitter.com/${author || 'i'}/status/${postId}`

    return (
      <NodeViewWrapper>
        <div
          ref={twitterRef}
          className={`embed-node-card embed-node-x${selectedClass}`}
          data-drag-handle=""
          contentEditable={false}
        >
          <blockquote className="twitter-tweet" data-theme="light">
            <a href={url} target="_blank" rel="noopener noreferrer"></a>
          </blockquote>
        </div>
      </NodeViewWrapper>
    )
  }

  if (node.type.name === "tmdbCardEmbed") {
    const { tmdbId, mediaType, title, overview, posterPath, releaseDate, voteAverage } = node.attrs as {
      tmdbId: string
      mediaType: string
      title: string
      overview: string
      posterPath: string
      releaseDate: string
      voteAverage: number
    }
    const selectedClass = selected ? " is-selected" : ""
    const posterUrl = posterPath ? `https://image.tmdb.org/t/p/w200${posterPath}` : ''

    return (
      <NodeViewWrapper>
        <div
          className={`embed-node-card embed-node-tmdb${selectedClass}`}
          data-drag-handle=""
          contentEditable={false}
        >
          {posterUrl && <img src={posterUrl} alt={title} className="tmdb-poster" />}
          <div className="tmdb-info">
            <div className="tmdb-title">{title || "(no title)"}</div>
            <div className="tmdb-meta">
              {mediaType === 'movie' ? '电影' : '剧集'} · {releaseDate?.slice(0, 4)} · ⭐ {voteAverage?.toFixed(1)}
            </div>
            {overview && <div className="tmdb-overview">{overview}</div>}
          </div>
        </div>
      </NodeViewWrapper>
    )
  }

  return null
}
