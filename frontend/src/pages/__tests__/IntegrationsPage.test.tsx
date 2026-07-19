import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import IntegrationsPage from '../IntegrationsPage'

vi.mock('../IntegrationsPage.module.css', () => ({ default: new Proxy({}, { get: (_, key) => key }) }))

const writeText = vi.fn().mockResolvedValue(undefined)

beforeEach(() => {
  writeText.mockClear()
})

describe('IntegrationsPage', () => {
  it('renders the MCP endpoint derived from the current origin', () => {
    render(<IntegrationsPage />)
    expect(screen.getByText('http://localhost:3000/api/v1/mcp')).toBeInTheDocument()
  })

  it('lists the three MCP tools with their real names and descriptions', () => {
    render(<IntegrationsPage />)
    expect(screen.getByText('search_memory')).toBeInTheDocument()
    expect(screen.getByText('Search semantically similar manuscript passages in a universe.')).toBeInTheDocument()
    expect(screen.getByText('query_entities')).toBeInTheDocument()
    expect(screen.getByText('Find an entity and its graph neighbours in a universe.')).toBeInTheDocument()
    expect(screen.getByText('recall')).toBeInTheDocument()
    expect(screen.getByText("Run Quill's hybrid memory recall for a universe.")).toBeInTheDocument()
  })

  it('copies the endpoint to the clipboard and shows a confirmation', async () => {
    const user = userEvent.setup()
    // userEvent.setup() installs its own clipboard stub on navigator.clipboard;
    // override it after setup so our spy wins.
    Object.defineProperty(navigator, 'clipboard', { value: { writeText }, configurable: true })
    render(<IntegrationsPage />)

    await user.click(screen.getByRole('button', { name: /copy/i }))

    expect(writeText).toHaveBeenCalledWith('http://localhost:3000/api/v1/mcp')
    expect(await screen.findByText(/copied/i)).toBeInTheDocument()
  })
})
