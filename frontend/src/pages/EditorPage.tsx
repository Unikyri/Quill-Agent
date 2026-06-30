import { useEffect, useRef, useCallback, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useEditorStore } from '../stores/editorStore'
import { useWSStore } from '../stores/wsStore'
import { useWS } from '../hooks/useWS'
import { api } from '../lib/api'
import TipTapEditor from '../components/editor/TipTapEditor'
import ContextPanel from '../components/context-panel/ContextPanel'

export default function EditorPage() {
  const { chapterId } = useParams<{ chapterId: string }>()
  const { content, wordCount, isSaving, lastSavedAt, setContent, saveContent } = useEditorStore()
  const wsStatus = useWSStore((s) => s.status)
  const saveTimerRef = useRef<ReturnType<typeof setTimeout>>()
  const [workId, setWorkId] = useState<string>('')
  const [universeId, setUniverseId] = useState<string>('')

  // Open WS connection
  useWS()

  // Load chapter data on mount
  useEffect(() => {
    if (chapterId) {
      api.getChapter(chapterId).then(({ chapter }) => {
        setContent(chapter.content || '', chapter.raw_text || '')
        // Derive workId and universeId from chapter if available
        if (chapter.work_id) setWorkId(chapter.work_id)
        if (chapter.universe_id) setUniverseId(chapter.universe_id)
      })
    }
  }, [chapterId])

  // Auto-save after 5 seconds of inactivity via editorStore
  const handleContentChange = useCallback((_html: string, text: string) => {
    setContent(_html, text)

    if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
    saveTimerRef.current = setTimeout(() => {
      if (chapterId) saveContent(chapterId)
    }, 5000)
  }, [chapterId, setContent, saveContent])

  // Traffic light indicator
  const statusIndicator = wsStatus === 'open' ? '🟢' : wsStatus === 'reconnecting' ? '🟡' : '🔴'

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      {/* Editor */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <div style={{
          padding: '12px 24px',
          borderBottom: '1px solid #333',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}>
          <span style={{ color: '#888', display: 'flex', alignItems: 'center', gap: 8 }}>
            Chapter Editor
            <span title={`WS: ${wsStatus}`} style={{ fontSize: 12 }}>{statusIndicator}</span>
          </span>
          <div style={{ display: 'flex', gap: 16, color: '#888', fontSize: 12 }}>
            <span>{wordCount} words</span>
            <span>{isSaving ? 'Saving...' : lastSavedAt ? `Saved ${lastSavedAt.toLocaleTimeString()}` : ''}</span>
          </div>
        </div>

        {chapterId && workId && universeId ? (
          <TipTapEditor
            chapterId={chapterId}
            workId={workId}
            universeId={universeId}
            initialContent={content}
            onContentChange={handleContentChange}
          />
        ) : (
          <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#555' }}>
            Loading editor…
          </div>
        )}
      </div>

      <ContextPanel status={wsStatus} />
    </div>
  )
}
