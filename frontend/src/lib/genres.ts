export const GENRE_OPTIONS = [
  { value: 'fantasy', label: 'Fantasy' },
  { value: 'epic-fantasy', label: 'Epic Fantasy' },
  { value: 'urban-fantasy', label: 'Urban Fantasy' },
  { value: 'romantasy', label: 'Romantasy' },
  { value: 'science-fiction', label: 'Science Fiction' },
  { value: 'space-opera', label: 'Space Opera' },
  { value: 'dystopian', label: 'Dystopian' },
  { value: 'horror', label: 'Horror' },
  { value: 'gothic', label: 'Gothic' },
  { value: 'paranormal', label: 'Paranormal' },
  { value: 'romance', label: 'Romance' },
  { value: 'mystery', label: 'Mystery' },
  { value: 'cozy-mystery', label: 'Cozy Mystery' },
  { value: 'thriller', label: 'Thriller' },
  { value: 'crime', label: 'Crime' },
  { value: 'historical', label: 'Historical' },
  { value: 'literary', label: 'Literary' },
  { value: 'adventure', label: 'Adventure' },
  { value: 'young-adult', label: 'Young Adult' },
  { value: 'coming-of-age', label: 'Coming of Age' },
] as const

export const WORK_FORMAT_OPTIONS = [
  { value: 'novel', label: 'Novel' },
  { value: 'novella', label: 'Novella' },
  { value: 'short-story', label: 'Short Story' },
] as const

export function selectedValues(select: HTMLSelectElement): string[] {
  return Array.from(select.selectedOptions, (option) => option.value)
}
