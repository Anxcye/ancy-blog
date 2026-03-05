/**
 * File: tmdb-card-node-extension.ts
 * Purpose: Define TipTap block node extension for TMDB movie/TV card embeds.
 * Module: frontend-admin-react/components/tiptap-node/embed-node, editor schema layer.
 * Related: XPostEmbed, EmbedNodeView, EmbedDropdownMenu, simple-editor.
 */

import { Node, mergeAttributes } from "@tiptap/core"
import { ReactNodeViewRenderer } from "@tiptap/react"

import { EmbedNodeView } from "./embed-node-view"

export type TmdbMediaType = "movie" | "tv"

export const TmdbCardEmbed = Node.create({
  name: "tmdbCardEmbed",
  group: "block",
  atom: true,
  selectable: true,
  draggable: true,

  addAttributes() {
    return {
      tmdbId: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-tmdb-id") ?? "",
        renderHTML: (attrs) => ({ "data-tmdb-id": attrs.tmdbId }),
      },
      mediaType: {
        default: "movie" as TmdbMediaType,
        parseHTML: (el) => (el.getAttribute("data-media-type") ?? "movie") as TmdbMediaType,
        renderHTML: (attrs) => ({ "data-media-type": attrs.mediaType }),
      },
      title: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-title") ?? "",
        renderHTML: (attrs) => ({ "data-title": attrs.title }),
      },
      overview: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-overview") ?? "",
        renderHTML: (attrs) => ({ "data-overview": attrs.overview }),
      },
      posterPath: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-poster-path") ?? "",
        renderHTML: (attrs) => ({ "data-poster-path": attrs.posterPath }),
      },
      releaseDate: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-release-date") ?? "",
        renderHTML: (attrs) => ({ "data-release-date": attrs.releaseDate }),
      },
      voteAverage: {
        default: 0,
        parseHTML: (el) => parseFloat(el.getAttribute("data-vote-average") ?? "0"),
        renderHTML: (attrs) => ({ "data-vote-average": attrs.voteAverage }),
      },
    }
  },

  parseHTML() {
    return [{ tag: 'div[data-embed-type="tmdb_card"]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return [
      "div",
      mergeAttributes(HTMLAttributes, { "data-embed-type": "tmdb_card" }),
    ]
  },

  addNodeView() {
    return ReactNodeViewRenderer(EmbedNodeView)
  },
})
