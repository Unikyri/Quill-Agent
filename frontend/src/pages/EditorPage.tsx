import { useEffect, useRef, useCallback, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useEditorStore } from '../stores/editorStore'
import { useWSStore } from '../stores/wsStore'
import { useWS } from '../hooks/useWS'
import { api } from '../lib/api'
import TipTapEditor from '../components/editor/TipTapEditor'
import ContextPanel from '../components/context-panel/ContextPanel'
import styles from './EditorPage.module.css'

export default function EditorPage() {
  const { chapterId } = useParams<{ chapterId: string }>()
  const { content, wordCount, isSaving, lastSavedAt, setContent, saveContent } = useEditorStore()
  const wsStatus = useWSStore((s) => s.status)
  const saveTimerRef = useRef<ReturnType<typeof setTimeout>>()
  const [workId, setWorkId] = useState<string>('')
  const [universeId, setUniverseId] = useState<string>('')

  useWS()

  useEffect(() => {
    if (chapterId) {
      api.getChapter(chapterId).then(({ chapter }) => {
        setContent(chapter.content || '', chapter.raw_text || '')
        if (chapter.work_id) setWorkId(chapter.work_id)
        if (chapter.universe_id) setUniverseId(chapter.universe_id)
      })
    }
  }, [chapterId])

  const handleContentChange = useCallback((_html: string, text: string) => {
    setContent(_html, text)

    if (saveTimerRef.current) clearTimeout(saveTimerRef.current)
    saveTimerRef.current = setTimeout(() => {
      if (chapterId) saveContent(chapterId)
    }, 5000)
  }, [chapterId, setContent, saveContent])

  const statusIndicator = wsStatus === 'open' ? '\u{1F7E2}' : wsStatus === 'reconnecting' ? '\u{1F7E1}' : '\u{1F534}'

  return (
    <div className={styles.wrap}>
      <div className={styles.editorPanel}>
        <div className={styles.headerBar}>
          <span className={styles.headerLeft}>
            Chapter Editor
            <span className={styles.wsIndicator} title={`WS: ${wsStatus}`}>{statusIndicator}</span>
          </span>
          <div className={styles.headerRight}>
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
          <div className={styles.loading}>Loading editor…</div>
        )}
      </div>

      <ContextPanel status={wsStatus} />
    </div>
  )
}
