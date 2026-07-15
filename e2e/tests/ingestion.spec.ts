import { expect, test, type Page } from '../../frontend/node_modules/@playwright/test'

const runID = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
const email = `ingestion-${runID}@quill.local`
const password = 'E2E-passphrase-123'

async function registerAndCreateUniverse(page: Page) {
	const universeName = `Ingestion Universe ${runID}`
  await page.goto('/login')
  await page.getByRole('button', { name: 'Register' }).click()
  await page.getByPlaceholder('Author Name').fill('E2E Ingestion Author')
  await page.getByPlaceholder('Email').fill(email)
  await page.getByPlaceholder('Password').fill(password)
  await page.getByRole('button', { name: 'Create My Account' }).click()
  await page.getByPlaceholder(/Universe Name/).fill(universeName)
  await page.locator('select[multiple]').selectOption(['science-fiction', 'fantasy'])
  await Promise.all([
    page.waitForURL(/\/universe\/[0-9a-f-]{36}(?:\/|$)/i),
    page.getByRole('button', { name: 'Create Universe' }).click(),
  ])
  await expect(page.getByRole('button', { name: new RegExp(universeName) })).toBeVisible()
}

test('EE-2: manuscript ingestion emits terminal progress and persists canonical entities, chapters, and graph edges', async ({ page }) => {
  test.skip(!process.env.QWEN_API_KEY, 'QWEN_API_KEY is required: this suite verifies the real model-backed loop.')
  const terminalProgress: Array<{ status?: string; error?: string; error_message?: string }> = []
  let resolveTerminal!: () => void
  let rejectTerminal!: (error: Error) => void
  const terminalIngestion = new Promise<void>((resolve, reject) => {
    resolveTerminal = resolve
    rejectTerminal = reject
  })
  page.on('websocket', (socket) => {
    socket.on('framereceived', (frame) => {
      if (typeof frame.payload !== 'string') return
      try {
        const message = JSON.parse(frame.payload)
        if (message.type === 'ingestion_progress') {
          terminalProgress.push(message.payload)
          if (message.payload.status === 'completed') resolveTerminal()
          if (message.payload.status === 'failed') {
            rejectTerminal(new Error(`ingestion failed: ${message.payload.error_message ?? message.payload.error ?? 'unknown reason'}`))
          }
        }
      } catch { /* non-JSON heartbeat */ }
    })
  })
  await registerAndCreateUniverse(page)
  const universePath = new URL(page.url()).pathname
  const universeMatch = universePath.match(/^\/universe\/([0-9a-f-]{36})(?:\/|$)/i)
  expect(universeMatch, `expected a universe UUID in route ${universePath}`).not.toBeNull()
  const universeID = universeMatch?.[1]
  if (!universeID) throw new Error(`universe UUID missing from route ${universePath}`)
  await page.getByRole('link', { name: 'Ingestion' }).click()
  await page.getByTestId('ingest-file-input').setInputFiles('../e2e/fixtures/world.md')
  await Promise.race([
    terminalIngestion,
    page.waitForTimeout(120_000).then(() => { throw new Error('ingestion did not reach a terminal status within 120 seconds') }),
  ])
  expect(terminalProgress.at(-1)).toMatchObject({ status: 'completed' })
  await expect(page.getByText('Completed')).toBeVisible()

  const persisted = await page.evaluate(async (id) => {
    const token = localStorage.getItem('token')
    const headers = token ? { Authorization: `Bearer ${token}` } : {}
    const request = async (path: string) => {
      const response = await fetch(`/api/v1${path}`, { headers })
      if (!response.ok) throw new Error(`${path}: ${response.status}`)
      return response.json()
    }
    const works = await request(`/universes/${id}/works`)
    const chapterLists = await Promise.all(works.works.map((work: { id: string }) => request(`/works/${work.id}/chapters`)))
    const entities = await request(`/universes/${id}/entities?limit=100`)
    const graph = await request(`/universes/${id}/graph`)
    return {
      chapterCount: chapterLists.flatMap((entry: { chapters: unknown[] }) => entry.chapters).length,
      entityTypes: entities.entities.map((entity: { type: string }) => entity.type),
      edgeCount: graph.edges.length,
    }
  }, universeID)
  expect(persisted.chapterCount).toBeGreaterThanOrEqual(10)
  expect(persisted.entityTypes.length).toBeGreaterThan(0)
  const canonicalTypes = new Set(['character', 'place', 'object', 'faction', 'event', 'world_rule', 'plot_arc'])
  expect(persisted.entityTypes.every((type: string) => canonicalTypes.has(type))).toBeTruthy()
  expect(persisted.edgeCount).toBeGreaterThan(0)

  await page.getByRole('link', { name: 'Entities' }).click()
  await expect(page.getByText(/Mira Voss|Aurelia Station|Sun Key/).first()).toBeVisible()
})
