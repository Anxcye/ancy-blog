/**
 * File: image-url-popover.tsx
 * Purpose: Insert or update editor images from a remote image URL.
 * Module: frontend-admin-react/components/tiptap-ui/image-url-popover, presentation layer.
 * Related: simple-editor.tsx, @tiptap/extension-image, and image-upload-button.
 */

import { forwardRef, useCallback, useEffect, useState } from "react"
import type { Editor } from "@tiptap/react"
import "@tiptap/extension-image"

import { useTiptapEditor } from "@/hooks/use-tiptap-editor"
import { ImagePlusIcon } from "@/components/tiptap-icons/image-plus-icon"
import { CornerDownLeftIcon } from "@/components/tiptap-icons/corner-down-left-icon"
import { Button, ButtonGroup } from "@/components/tiptap-ui-primitive/button"
import type { ButtonProps } from "@/components/tiptap-ui-primitive/button"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/tiptap-ui-primitive/popover"
import {
  Card,
  CardBody,
  CardItemGroup,
} from "@/components/tiptap-ui-primitive/card"
import { Input, InputGroup } from "@/components/tiptap-ui-primitive/input"

export interface ImageUrlPopoverProps extends Omit<ButtonProps, "type"> {
  editor?: Editor | null
}

function isValidImageUrl(value: string): boolean {
  if (!value.trim()) return false
  try {
    const parsed = new URL(value.trim())
    return parsed.protocol === "http:" || parsed.protocol === "https:"
  } catch {
    return false
  }
}

export const ImageUrlButton = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <Button
        type="button"
        className={className}
        variant="ghost"
        role="button"
        tabIndex={-1}
        aria-label="Insert image by URL"
        tooltip="Insert image by URL"
        ref={ref}
        {...props}
      >
        {children || (
          <>
            <ImagePlusIcon className="tiptap-button-icon" />
            <span className="tiptap-button-text">URL</span>
          </>
        )}
      </Button>
    )
  },
)

ImageUrlButton.displayName = "ImageUrlButton"

export const ImageUrlPopover = forwardRef<HTMLButtonElement, ImageUrlPopoverProps>(
  ({ editor: providedEditor, onClick, children, ...buttonProps }, ref) => {
    const { editor } = useTiptapEditor(providedEditor)
    const [isOpen, setIsOpen] = useState(false)
    const [src, setSrc] = useState("")
    const [alt, setAlt] = useState("")

    const imageAttrs = (() => {
      if (!editor?.isActive("image")) {
        return { src: "", alt: "" }
      }
      const attrs = editor.getAttributes("image") as {
        src?: string
        alt?: string
      }
      return {
        src: attrs.src ?? "",
        alt: attrs.alt ?? "",
      }
    })()

    useEffect(() => {
      if (!isOpen) return
      setSrc(imageAttrs.src)
      setAlt(imageAttrs.alt)
    }, [imageAttrs.alt, imageAttrs.src, isOpen])

    const canApply = !!editor && editor.isEditable && isValidImageUrl(src)

    const applyImage = useCallback(() => {
      if (!editor || !canApply) return
      const normalizedSrc = src.trim()
      const normalizedAlt = alt.trim()
      const chain = editor.chain().focus()

      if (editor.isActive("image")) {
        chain
          .updateAttributes("image", {
            src: normalizedSrc,
            alt: normalizedAlt || null,
          })
          .run()
      } else {
        chain
          .setImage({ src: normalizedSrc, alt: normalizedAlt || undefined })
          .run()
      }

      setIsOpen(false)
    }, [alt, canApply, editor, src])

    const handleClick = useCallback(
      (event: React.MouseEvent<HTMLButtonElement>) => {
        onClick?.(event)
        if (event.defaultPrevented) return
        setIsOpen((current) => !current)
      },
      [onClick]
    )

    const handleKeyDown = useCallback(
      (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key !== "Enter") return
        event.preventDefault()
        applyImage()
      },
      [applyImage]
    )

    return (
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <ImageUrlButton
            ref={ref}
            onClick={handleClick}
            aria-pressed={isOpen}
            disabled={!editor?.isEditable}
            {...buttonProps}
          >
            {children}
          </ImageUrlButton>
        </PopoverTrigger>
        <PopoverContent align="start">
          <Card>
            <CardBody>
              <CardItemGroup orientation="vertical">
                <InputGroup>
                  <Input
                    type="url"
                    placeholder="https://example.com/image.jpg"
                    value={src}
                    onChange={(event) => setSrc(event.target.value)}
                    onKeyDown={handleKeyDown}
                    autoFocus
                    autoComplete="off"
                    autoCorrect="off"
                    autoCapitalize="off"
                  />
                </InputGroup>
                <InputGroup>
                  <Input
                    type="text"
                    placeholder="Alt text (optional)"
                    value={alt}
                    onChange={(event) => setAlt(event.target.value)}
                    onKeyDown={handleKeyDown}
                  />
                </InputGroup>
                <ButtonGroup orientation="horizontal">
                  <Button
                    type="button"
                    variant="primary"
                    onClick={applyImage}
                    disabled={!canApply}
                  >
                    <CornerDownLeftIcon className="tiptap-button-icon" />
                    <span className="tiptap-button-text">Apply</span>
                  </Button>
                </ButtonGroup>
              </CardItemGroup>
            </CardBody>
          </Card>
        </PopoverContent>
      </Popover>
    )
  }
)

ImageUrlPopover.displayName = "ImageUrlPopover"
