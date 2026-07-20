// crypto.randomUUID() is gated to secure contexts (HTTPS or localhost).
// This app is routinely deployed over plain HTTP on a public IP for judging,
// where that API is undefined — crypto.getRandomValues() has no such
// restriction, so it's the fallback. Returns null only if neither exists.
export function randomUUIDCompat(): string | null {
  if (typeof globalThis.crypto?.randomUUID === 'function') {
    return globalThis.crypto.randomUUID()
  }

  if (typeof globalThis.crypto?.getRandomValues === 'function') {
    const bytes = globalThis.crypto.getRandomValues(new Uint8Array(16))
    bytes[6] = (bytes[6] & 0x0f) | 0x40
    bytes[8] = (bytes[8] & 0x3f) | 0x80
    const hex = Array.from(bytes, (b) => b.toString(16).padStart(2, '0')).join('')
    return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`
  }

  return null
}
