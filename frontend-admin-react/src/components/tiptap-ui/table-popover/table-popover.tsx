/**
 * File: table-popover.tsx
 * Purpose: Provide table insertion and basic row or column management controls.
 * Module: frontend-admin-react/components/tiptap-ui/table-popover, presentation layer.
 * Related: simple-editor.tsx and TipTap table extensions.
 */

import { forwardRef, useCallback, useMemo, useState } from "react"
import type { Editor } from "@tiptap/react"
import "@tiptap/extension-table"

import { useTiptapEditor } from "@/hooks/use-tiptap-editor"
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
  CardGroupLabel,
  CardItemGroup,
} from "@/components/tiptap-ui-primitive/card"

export interface TablePopoverProps extends Omit<ButtonProps, "type"> {
  editor?: Editor | null
}

export const TableButton = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <Button
        type="button"
        className={className}
        variant="ghost"
        role="button"
        tabIndex={-1}
        aria-label="Table"
        tooltip="Table"
        ref={ref}
        {...props}
      >
        {children || <span className="tiptap-button-text">Table</span>}
      </Button>
    )
  },
)

TableButton.displayName = "TableButton"

export const TablePopover = forwardRef<HTMLButtonElement, TablePopoverProps>(
  ({ editor: providedEditor, onClick, children, ...buttonProps }, ref) => {
    const { editor } = useTiptapEditor(providedEditor)
    const [isOpen, setIsOpen] = useState(false)
    const isInTable = !!editor?.isActive("table")

    const actions = useMemo(
      () => [
        {
          label: "Insert 3x3",
          enabled: !!editor?.isEditable,
          run: () =>
            editor
              ?.chain()
              .focus()
              .insertTable({ rows: 3, cols: 3, withHeaderRow: true })
              .run(),
        },
        {
          label: "Toggle header",
          enabled: !!editor?.isEditable && isInTable,
          run: () => editor?.chain().focus().toggleHeaderRow().run(),
        },
        {
          label: "Add row",
          enabled: !!editor?.isEditable && isInTable,
          run: () => editor?.chain().focus().addRowAfter().run(),
        },
        {
          label: "Add column",
          enabled: !!editor?.isEditable && isInTable,
          run: () => editor?.chain().focus().addColumnAfter().run(),
        },
        {
          label: "Delete row",
          enabled: !!editor?.isEditable && isInTable,
          run: () => editor?.chain().focus().deleteRow().run(),
        },
        {
          label: "Delete column",
          enabled: !!editor?.isEditable && isInTable,
          run: () => editor?.chain().focus().deleteColumn().run(),
        },
        {
          label: "Delete table",
          enabled: !!editor?.isEditable && isInTable,
          run: () => editor?.chain().focus().deleteTable().run(),
        },
      ],
      [editor, isInTable]
    )

    const handleClick = useCallback(
      (event: React.MouseEvent<HTMLButtonElement>) => {
        onClick?.(event)
        if (event.defaultPrevented) return
        setIsOpen((current) => !current)
      },
      [onClick]
    )

    const handleAction = useCallback((run: () => boolean | undefined) => {
      run()
      setIsOpen(false)
    }, [])

    return (
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <TableButton
            ref={ref}
            onClick={handleClick}
            aria-pressed={isOpen}
            disabled={!editor?.isEditable}
            {...buttonProps}
          >
            {children}
          </TableButton>
        </PopoverTrigger>
        <PopoverContent align="start">
          <Card>
            <CardBody>
              <CardItemGroup orientation="vertical">
                <CardGroupLabel>Table</CardGroupLabel>
                <ButtonGroup orientation="vertical">
                  {actions.map((action) => (
                    <Button
                      key={action.label}
                      type="button"
                      variant="ghost"
                      disabled={!action.enabled}
                      onClick={() => handleAction(action.run)}
                    >
                      <span className="tiptap-button-text">{action.label}</span>
                    </Button>
                  ))}
                </ButtonGroup>
              </CardItemGroup>
            </CardBody>
          </Card>
        </PopoverContent>
      </Popover>
    )
  }
)

TablePopover.displayName = "TablePopover"
