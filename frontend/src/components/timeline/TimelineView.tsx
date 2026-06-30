import styles from './TimelineView.module.css'

export interface TimelineEvent {
  id: string
  label: string
  timestamp: string
  description: string
  chapter?: string
}

interface TimelineViewProps {
  events: TimelineEvent[]
}

export default function TimelineView({ events }: TimelineViewProps) {
  return (
    <div className={styles.timeline}>
      {events.map((event) => (
        <div key={event.id} className={styles.eventCard}>
          <div className={styles.eventDot} />
          {event.chapter && <p className={styles.eventChapter}>{event.chapter}</p>}
          <h4 className={styles.eventLabel}>{event.label}</h4>
          <p className={styles.eventDesc}>{event.description}</p>
        </div>
      ))}
    </div>
  )
}
