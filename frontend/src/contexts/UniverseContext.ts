import { createContext } from 'react'

export interface UniverseContextValue {
  universe: { id: string; name: string; genre_tags?: string[]; genre?: string; description?: string } | null
  works: { id: string; title: string; type: string; order_index: number }[]
  refetchWorks: () => Promise<void>
}

export const UniverseContext = createContext<UniverseContextValue>({
  universe: null,
  works: [],
  refetchWorks: async () => {},
})
