import '@testing-library/jest-dom/vitest'

// ponyail: ReactFlow v11 needs ResizeObserver in jsdom; mock it globally
globalThis.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
} as unknown as typeof ResizeObserver
