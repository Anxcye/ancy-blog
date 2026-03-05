/**
 * File: embed-dropdown-menu.tsx
 * Purpose: Toolbar button that opens an Ant Design modal to insert embed blocks.
 * Module: frontend-admin-react/components/tiptap-ui/embed-dropdown-menu, editor UI layer.
 * Related: XPostEmbed, TmdbCardEmbed, simple-editor toolbar.
 */

import { useState } from "react"
import { Form, Input, Modal, Segmented, Select, message } from "antd"
import { useCurrentEditor } from "@tiptap/react"

import { Button } from "@/components/tiptap-ui-primitive/button"
import type { TmdbMediaType } from "@/components/tiptap-node/embed-node/tmdb-card-node-extension"
import { httpClient } from "@/lib/http"

type EmbedType = "x_post" | "tmdb_card"

interface EmbedFormValues {
  // x_post fields
  postId: string
  author?: string
  // tmdb fields
  tmdbId: string
  mediaType: TmdbMediaType
  title?: string
}

export function EmbedDropdownMenu() {
  const { editor } = useCurrentEditor()
  const [open, setOpen] = useState(false)
  const [embedType, setEmbedType] = useState<EmbedType>("x_post")
  const [form] = Form.useForm<EmbedFormValues>()
  const [loading, setLoading] = useState(false)

  const handleOpen = () => {
    setOpen(true)
  }

  const handleTypeChange = (val: EmbedType) => {
    setEmbedType(val)
    form.resetFields()
  }

  const handleInsert = async () => {
    try {
      const values = await form.validateFields()
      setLoading(true)

      if (embedType === "x_post") {
        editor?.commands.insertContent({
          type: "xPostEmbed",
          attrs: {
            postId: values.postId.trim(),
            author: (values.author ?? "").trim(),
          },
        })
      } else {
        // Fetch TMDB metadata from backend
        const { data } = await httpClient.get(`/admin/integrations/tmdb/${values.mediaType}/${values.tmdbId.trim()}`)
        const metadata = data.data

        editor?.commands.insertContent({
          type: "tmdbCardEmbed",
          attrs: {
            tmdbId: metadata.id.toString(),
            mediaType: values.mediaType ?? "movie",
            title: metadata.title,
            overview: metadata.overview,
            posterPath: metadata.posterPath,
            releaseDate: metadata.releaseDate,
            voteAverage: metadata.voteAverage,
          },
        })
      }
      setOpen(false)
      form.resetFields()
      editor?.commands.focus()
    } catch (error: any) {
      if (error?.response?.data?.message) {
        message.error(error.response.data.message)
      }
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = () => {
    setOpen(false)
    form.resetFields()
  }

  return (
    <>
      <Button
        variant="ghost"
        tooltip="Insert embed block"
        onClick={handleOpen}
        disabled={!editor?.isEditable}
        aria-label="Insert embed"
      >
        <span style={{ fontSize: 12, fontWeight: 600, letterSpacing: "0.02em" }}>
          Embed
        </span>
      </Button>

      <Modal
        title="Insert embed block"
        open={open}
        onOk={handleInsert}
        onCancel={handleCancel}
        okText="Insert"
        cancelText="Cancel"
        confirmLoading={loading}
        width={460}
        destroyOnClose
        // Prevent editor from losing focus handling
        getContainer={false}
      >
        <div style={{ marginBottom: 16 }}>
          <Segmented<EmbedType>
            value={embedType}
            onChange={handleTypeChange}
            options={[
              { label: "𝕏  X Post", value: "x_post" },
              { label: "🎬  TMDB Card", value: "tmdb_card" },
            ]}
            block
          />
        </div>

        <Form form={form} layout="vertical" requiredMark={false}>
          {embedType === "x_post" ? (
            <>
              <Form.Item
                name="postId"
                label="Post ID"
                rules={[{ required: true, message: "请填写 Post ID" }]}
              >
                <Input placeholder="e.g. 1860000000000000000" autoFocus />
              </Form.Item>
              <Form.Item name="author" label="Author（选填）">
                <Input placeholder="@username" />
              </Form.Item>
            </>
          ) : (
            <>
              <Form.Item
                name="tmdbId"
                label="TMDB ID"
                rules={[{ required: true, message: "请填写 TMDB ID" }]}
              >
                <Input placeholder="e.g. 550" autoFocus />
              </Form.Item>
              <Form.Item name="mediaType" label="类型" initialValue="movie">
                <Select
                  options={[
                    { value: "movie", label: "Movie 电影" },
                    { value: "tv", label: "TV Show 电视剧" },
                  ]}
                />
              </Form.Item>
              <Form.Item name="title" label="Title（选填）">
                <Input placeholder="e.g. Fight Club" />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </>
  )
}
