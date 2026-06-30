import styles from './PageStatus.module.css'

interface PageStatusProps {
  loading?: boolean
  error?: string | null
  onRetry?: () => void
}

export default function PageStatus({ loading, error, onRetry }: PageStatusProps) {
  if (loading) {
    return (
      <div className={styles.statusWrap} data-testid="loading-state">
        <div className={styles.spinner} />
        <p className={styles.loadingText}>Loading…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className={styles.statusWrap} data-testid="error-state">
        <p className={styles.errorText}>{error}</p>
        {onRetry && (
          <button className={styles.retryBtn} onClick={onRetry}>
            Retry
          </button>
        )}
      </div>
    )
  }

  return null
}
