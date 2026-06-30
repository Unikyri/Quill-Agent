import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import ContextCard from '../ContextCard'

beforeEach(() => {
  vi.useFakeTimers()
})

afterEach(() => {
  vi.useRealTimers()
})

describe('ContextCard', () => {
  it('renders recall card with correct content', () => {
    render(
      <ContextCard
        id="1"
        type="recall"
        title="Memory"
        detail="Something recalled"
        onDismiss={vi.fn()}
      />
    )
    expect(screen.getByText(/Memory/)).toBeInTheDocument()
    expect(screen.getByText('Something recalled')).toBeInTheDocument()
  })

  it('renders contradiction with severity badge', () => {
    render(
      <ContextCard
        id="2"
        type="contradiction"
        title="Conflict"
        detail="Two facts contradict"
        severity="high"
        onDismiss={vi.fn()}
      />
    )
    expect(screen.getByText('HIGH')).toBeInTheDocument()
  })

  it('renders entity card with NEW badge', () => {
    render(
      <ContextCard
        id="3"
        type="entity"
        title="Alice"
        detail="Type: character"
        isNew
        onDismiss={vi.fn()}
      />
    )
    expect(screen.getByText('NEW')).toBeInTheDocument()
  })

  it('calls onDismiss when X button clicked', () => {
    const onDismiss = vi.fn()
    render(
      <ContextCard
        id="4"
        type="recall"
        title="Test"
        detail="Test detail"
        onDismiss={onDismiss}
      />
    )
    fireEvent.click(screen.getByTitle('Dismiss'))
    // The dismiss triggers a 400ms fade animation before calling onDismiss
    vi.advanceTimersByTime(500)
    expect(onDismiss).toHaveBeenCalledWith('4')
  })

  it('auto-fades after 10 seconds', () => {
    const onDismiss = vi.fn()
    render(
      <ContextCard
        id="5"
        type="recall"
        title="Auto"
        detail="Auto fade"
        onDismiss={onDismiss}
      />
    )

    // Advance 10s to trigger fade start
    vi.advanceTimersByTime(10000)
    // After 10s, the fade timer fires, setting fading=true
    // Then 400ms more to finish the fade animation
    vi.advanceTimersByTime(400)

    expect(onDismiss).toHaveBeenCalledWith('5')
  })

  it('clears auto-fade timer on manual dismiss', () => {
    const onDismiss = vi.fn()
    render(
      <ContextCard
        id="6"
        type="recall"
        title="Manual"
        detail="Manual dismiss"
        onDismiss={onDismiss}
      />
    )

    // Advance only 5s (not enough for auto-fade)
    vi.advanceTimersByTime(5000)
    fireEvent.click(screen.getByTitle('Dismiss'))

    // Wait for the fade-out animation
    vi.advanceTimersByTime(500)

    expect(onDismiss).toHaveBeenCalledWith('6')
  })
})
