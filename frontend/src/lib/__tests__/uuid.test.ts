import { describe, it, expect, vi, afterEach } from 'vitest'
import { randomUUIDCompat } from '../uuid'

describe('randomUUIDCompat', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('uses crypto.randomUUID when available (secure context)', () => {
    expect(randomUUIDCompat()).toMatch(/^[0-9a-f-]{36}$/)
  })

  it('falls back to crypto.getRandomValues when randomUUID is undefined (insecure context)', () => {
    const original = globalThis.crypto
    vi.stubGlobal('crypto', { getRandomValues: original.getRandomValues.bind(original) })

    expect(randomUUIDCompat()).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/)
  })

  it('returns null when neither API is available', () => {
    vi.stubGlobal('crypto', {})
    expect(randomUUIDCompat()).toBeNull()
  })
})
