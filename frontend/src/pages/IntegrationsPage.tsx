import { useState } from 'react'
import styles from './IntegrationsPage.module.css'

// MCP tool names/descriptions copied verbatim from backend/internal/mcp/server.go
// so this screen never drifts from the real tool schema.
const TOOLS = [
  { name: 'search_memory', description: 'Search semantically similar manuscript passages in a universe.' },
  { name: 'query_entities', description: 'Find an entity and its graph neighbours in a universe.' },
  { name: 'recall', description: "Run Quill's hybrid memory recall for a universe." },
] as const

// Derived client-side, not fetched — matches the real POST /api/v1/mcp route
// registered in backend/cmd/server/main.go. No backend addition needed.
const MCP_ENDPOINT = `${window.location.origin}/api/v1/mcp`

export default function IntegrationsPage() {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(MCP_ENDPOINT)
      setCopied(true)
      return
    }

    // navigator.clipboard is undefined outside secure contexts (plain HTTP
    // on a public IP, e.g. this app's hackathon deployment) — fall back to
    // the legacy execCommand copy path.
    const textarea = document.createElement('textarea')
    textarea.value = MCP_ENDPOINT
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.focus()
    textarea.select()
    const ok = document.execCommand('copy')
    document.body.removeChild(textarea)
    if (ok) setCopied(true)
  }

  return (
    <main className={styles.wrap}>
      <div className={styles.heading}>
        <p className={styles.eyebrow}>Account</p>
        <h1>Connect Quill&apos;s memory to other AI tools</h1>
      </div>

      <div className={styles.endpointRow}>
        <span className={styles.endpointLabel}>MCP endpoint</span>
        <code className={styles.endpointValue}>{MCP_ENDPOINT}</code>
        <button type="button" className={styles.copyButton} onClick={handleCopy}>
          {copied ? 'Copied' : 'Copy'}
        </button>
      </div>

      <div>
        <h2 className={styles.toolsTitle}>Available tools</h2>
        <ul className={styles.toolList}>
          {TOOLS.map((tool) => (
            <li key={tool.name} className={styles.toolRow}>
              <code className={styles.toolName}>{tool.name}</code>
              <span className={styles.toolDescription}>{tool.description}</span>
            </li>
          ))}
        </ul>
      </div>

      <p className={styles.intro}>
        Paste this endpoint into Claude Desktop / Cursor&apos;s MCP settings to let another AI assistant query this
        universe&apos;s memory.
      </p>
    </main>
  )
}
