import { useEffect, useRef, useState } from 'react'
import { useEditor, EditorContent, type Editor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import Highlight from '@tiptap/extension-highlight'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import { useWSStore } from '../../stores/wsStore'
import styles from './TipTapEditor.module.css'

interface TipTapEditorProps {
  chapterId: string
  workId: string
  universeId: string
  initialContent?: string
  onContentChange?: (html: string, text: string) => void
}

function ToolbarButton({
  active, title, children, onClick,
}: {
  active?: boolean; title: string; children: React.ReactNode; onClick: () => void
}) {
  return (
    <button
      type="button"
      className={`${styles.toolbarBtn} ${active ? styles.toolbarBtnActive : ''}`}
      title={title}
      onClick={onClick}
      tabIndex={-1}
    >
      {children}
    </button>
  )
}

function Toolbar({ editor, fontSize, setFontSize }: { editor: Editor | null, fontSize: number, setFontSize: (s: number) => void }) {
  if (!editor) return null
  return (
    <div className={styles.toolbar}>
      <ToolbarButton title="Decrease font size" onClick={() => setFontSize(Math.max(12, fontSize - 1))}>
        <span style={{ fontSize: 13 }}>A-</span>
      </ToolbarButton>
      <ToolbarButton title="Increase font size" onClick={() => setFontSize(Math.min(32, fontSize + 1))}>
        <span style={{ fontSize: 15 }}>A+</span>
      </ToolbarButton>
      <div className={styles.toolbarDivider} />
      <ToolbarButton
        active={editor.isActive('bold')}
        title="Bold (⌘B)"
        onClick={() => editor.chain().focus().toggleBold().run()}
      >
        <b>B</b>
      </ToolbarButton>
      <ToolbarButton
        active={editor.isActive('italic')}
        title="Italic (⌘I)"
        onClick={() => editor.chain().focus().toggleItalic().run()}
      >
        <i>I</i>
      </ToolbarButton>
      <ToolbarButton
        active={editor.isActive('underline')}
        title="Underline (⌘U)"
        onClick={() => editor.chain().focus().toggleUnderline().run()}
      >
        <u>U</u>
      </ToolbarButton>
      <div className={styles.toolbarDivider} />
      <ToolbarButton
        active={editor.isActive('heading', { level: 1 })}
        title="Heading 1"
        onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
      >
        H1
      </ToolbarButton>
      <ToolbarButton
        active={editor.isActive('heading', { level: 2 })}
        title="Heading 2"
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
      >
        H2
      </ToolbarButton>
      <div className={styles.toolbarDivider} />
      <ToolbarButton
        active={editor.isActive('bulletList')}
        title="Bullet list"
        onClick={() => editor.chain().focus().toggleBulletList().run()}
      >
        •
      </ToolbarButton>
      <ToolbarButton
        active={editor.isActive('orderedList')}
        title="Numbered list"
        onClick={() => editor.chain().focus().toggleOrderedList().run()}
      >
        1.
      </ToolbarButton>
      <div className={styles.toolbarDivider} />
      <ToolbarButton
        active={editor.isActive('highlight')}
        title="Highlight"
        onClick={() => editor.chain().focus().toggleHighlight().run()}
      >
        <span style={{ fontSize: 11 }}>▐</span>
      </ToolbarButton>
      <ToolbarButton
        active={editor.isActive('blockquote')}
        title="Quote block"
        onClick={() => editor.chain().focus().toggleBlockquote().run()}
      >
        "
      </ToolbarButton>
    </div>
  )
}

export default function TipTapEditor({
  chapterId,
  workId,
  universeId,
  initialContent,
  onContentChange,
}: TipTapEditorProps) {
  const send = useWSStore((s) => s.send)
  const submitTimerRef = useRef<ReturnType<typeof setTimeout>>()
  const lastParagraphTextRef = useRef<string>('')
  const [fontSize, setFontSize] = useState(17)

  const editor = useEditor({
    extensions: [
      StarterKit,
      Placeholder.configure({ placeholder: 'Start writing…' }),
      Highlight,
      Underline,
      Link.configure({ openOnClick: false }),
    ],
    content: initialContent || '',
    onUpdate: ({ editor }) => {
      const html = editor.getHTML()
      const text = editor.getText()
      onContentChange?.(html, text)

      // Debounced paragraph submit for live AI analysis
      if (submitTimerRef.current) clearTimeout(submitTimerRef.current)
      submitTimerRef.current = setTimeout(() => {
        const paragraph = getParagraphAtCursor(editor)
        if (paragraph && paragraph.text.trim() && paragraph.text !== lastParagraphTextRef.current) {
          lastParagraphTextRef.current = paragraph.text
          send({
            type: 'paragraph_submit',
            payload: {
              work_id: workId,
              chapter_id: chapterId,
              universe_id: universeId,
              text: paragraph.text,
            },
          })
        }
      }, 5000)
    },
  })

  // Sync initial content when chapter changes
  useEffect(() => {
    if (editor && initialContent !== undefined && editor.getHTML() !== initialContent) {
      if (!editor.getText() || initialContent !== editor.getHTML()) {
        editor.commands.setContent(initialContent || '')
      }
    }
  }, [chapterId]) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    return () => {
      if (submitTimerRef.current) clearTimeout(submitTimerRef.current)
    }
  }, [])

  return (
    <div className={styles.wrapper}>
      <Toolbar editor={editor} fontSize={fontSize} setFontSize={setFontSize} />
      <div className={`${styles.editorContent} q-scroll`} style={{ fontSize: `${fontSize}px` }}>
        <EditorContent editor={editor} />
      </div>
    </div>
  )
}

function getParagraphAtCursor(editor: Editor): { text: string } | null {
  const { from } = editor.state.selection
  const doc = editor.state.doc
  const resolved = doc.resolve(from)
  let node = resolved.parent
  while (node && !node.isBlock && resolved.depth > 0) {
    const parentResolved = doc.resolve(resolved.before(resolved.depth))
    node = parentResolved.parent
  }
  if (!node || !node.isBlock) return null
  return { text: node.textContent }
}
