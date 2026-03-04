/**
 * File: x-post-node-extension.ts
 * Purpose: Define TipTap block node extension for X (Twitter) post embeds.
 * Module: frontend-admin-react/components/tiptap-node/embed-node, editor schema layer.
 * Related: TmdbCardEmbed, EmbedNodeView, EmbedDropdownMenu, simple-editor.
 */

import { Node, mergeAttributes } from "@tiptap/core"
import { ReactNodeViewRenderer } from "@tiptap/react"

import { EmbedNodeView } from "./embed-node-view"

export const XPostEmbed = Node.create({
  name: "xPostEmbed",
  group: "block",
  atom: true,
  selectable: true,
  draggable: true,

  addAttributes() {
    return {
      postId: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-post-id") ?? "",
        renderHTML: (attrs) => ({ "data-post-id": attrs.postId }),
      },
      author: {
        default: "",
        parseHTML: (el) => el.getAttribute("data-author") ?? "",
        renderHTML: (attrs) => ({ "data-author": attrs.author }),
      },
    }
  },

  parseHTML() {
    return [{ tag: 'div[data-embed-type="x_post"]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return [
      "div",
      mergeAttributes(HTMLAttributes, { "data-embed-type": "x_post" }),
    ]
  },

  addNodeView() {
    return ReactNodeViewRenderer(EmbedNodeView)
  },
})
