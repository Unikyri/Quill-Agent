import { expect, test, type Page } from '../../frontend/node_modules/@playwright/test'

const runID = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
const email = `e2e-${runID}@quill.local`
const password = 'E2E-passphrase-123'

async function register(page: Page) {
  await page.goto('/login')
  await page.getByRole('button', { name: 'Register' }).click()
  await page.getByPlaceholder('Author Name').fill('E2E Author')
  await page.getByPlaceholder('Email').fill(email)
  await page.getByPlaceholder('Password').fill(password)
  await page.getByRole('button', { name: 'Create My Account' }).click()
  await expect(page.getByRole('button', { name: 'Create Universe' })).toBeVisible()
}

async function createWritingWorkspace(page: Page) {
  await page.getByPlaceholder(/Universe Name/).fill(`E2E Universe ${runID}`)
  await page.locator('select[multiple]').selectOption(['fantasy', 'mystery'])
  await page.getByRole('button', { name: 'Create Universe' }).click()

  await page.getByRole('link', { name: /Works & Chapters/ }).click()
  await page.getByRole('button', { name: '+ New Work' }).click()
  await page.getByPlaceholder('Work title').fill('E2E Novel')
  await page.locator('select').last().selectOption('novel')
  await page.getByRole('button', { name: 'Create' }).click()
  await page.getByRole('button', { name: '+ New Chapter' }).click()
  await page.getByPlaceholder('Chapter title').fill('Chapter One')
  await page.getByRole('button', { name: 'Create' }).click()
  await page.getByText('Chapter One', { exact: true }).first().click()
  await expect(page.locator('.ProseMirror')).toBeVisible()
}

test.describe.configure({ mode: 'serial' })

test.beforeEach(async ({ page }) => {
  test.skip(!process.env.QWEN_API_KEY, 'QWEN_API_KEY is required: this suite verifies the real model-backed loop.')
  await register(page)
})

test('EE-1: editor submission produces an entity-bearing terminal result and entity browser entry', async ({ page }) => {
  const analysisResults: Array<{ entities?: unknown[] }> = []
  const analysisFailures: Array<{ reason?: string; cause?: string }> = []
  let resolveTerminal!: () => void
  let rejectTerminal!: (error: Error) => void
  const terminalAnalysis = new Promise<void>((resolve, reject) => {
    resolveTerminal = resolve
    rejectTerminal = reject
  })
  page.on('websocket', (socket) => {
    socket.on('framereceived', (frame) => {
      if (typeof frame.payload !== 'string') return
      try {
        const message = JSON.parse(frame.payload)
        if (message.type === 'analysis_result') {
          analysisResults.push(message.payload)
          resolveTerminal()
        }
        if (message.type === 'analysis_failed') {
          analysisFailures.push(message.payload)
          rejectTerminal(new Error(`analysis_failed: ${message.payload.reason ?? 'unknown reason'}${message.payload.cause ? ` (${message.payload.cause})` : ''}`))
        }
      } catch { /* non-JSON heartbeat */ }
    })
  })

  await createWritingWorkspace(page)
  await page.locator('.ProseMirror').fill('Captain Aster carries the Starfall Compass through the city of Elaris.')
  await expect(page.getByTestId('analysis-submission-status')).toHaveAttribute('data-phase', /submitted|analyzing/)
  await Promise.race([
    terminalAnalysis,
    page.waitForTimeout(120_000).then(() => { throw new Error('analysis_result did not arrive within 120 seconds') }),
  ])
  expect(analysisFailures).toHaveLength(0)
  expect(analysisResults.some((result) => Array.isArray(result.entities) && result.entities.length > 0)).toBeTruthy()
  await expect(page.getByTestId('analysis-submission-status')).toHaveAttribute('data-phase', 'done')

  await page.getByRole('link', { name: 'Entities' }).click()
  await expect(page.getByText(/Aster|Starfall Compass|Elaris/).first()).toBeVisible()
})
