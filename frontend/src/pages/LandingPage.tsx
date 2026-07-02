import { useRef, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from './LandingPage.module.css'

const features = [
  {
    title: 'Build Worlds',
    text: 'Create rich universes with characters, locations, and lore. Every detail tracked so nothing falls through the cracks.',
  },
  {
    title: 'Write with AI',
    text: 'Qwen-powered intelligence helps you navigate contradictions, fill plot holes, and keep your timeline consistent.',
  },
  {
    title: 'See Connections',
    text: 'An interactive knowledge graph reveals relationships between every entity, event, and chapter in your universe.',
  },
  {
    title: 'Stay Organized',
    text: 'Works, chapters, drafts — all structured like a manuscript. Focus on writing, not on managing files.',
  },
]

export default function LandingPage() {
  const navigate = useNavigate()
  const featureCardsRef = useRef<(HTMLDivElement | null)[]>([])

  // ponytail: IntersectionObserver for scroll-triggered reveals (replaces GSAP ScrollTrigger)
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) {
            entry.target.classList.add(styles.visible)
            observer.unobserve(entry.target)
          }
        }
      },
      { threshold: 0, rootMargin: '0px 0px -80px 0px' }
    )

    // Feature cards
    featureCardsRef.current.forEach((card) => {
      if (card) observer.observe(card)
    })

    // Closing section elements
    document.querySelectorAll('[data-anim="quote"], [data-anim="author"]').forEach((el) => {
      observer.observe(el)
    })

    return () => observer.disconnect()
  }, [])

  const scrollTo = (id: string) => {
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div className={styles.page}>
      {/* Navbar */}
      <nav className={styles.navbar}>
        <span className={styles.logo}>Quill</span>
        <div className={styles.navLinks}>
          <button className={styles.navLink} onClick={() => scrollTo('features')}>
            Features
          </button>
          <button className={styles.navLink} onClick={() => scrollTo('closing')}>
            About
          </button>
          <button className={styles.navCta} onClick={() => navigate('/login')}>
            Start Writing
          </button>
        </div>
      </nav>

      {/* Hero */}
      <section className={styles.hero}>
        <p className={styles.heroTagline} data-anim="tagline">
          For writers who build worlds
        </p>
        <h1 className={styles.heroTitle} data-anim="title">
          Write universes,
          <br />
          not just stories.
        </h1>
        <p className={styles.heroSubtitle} data-anim="subtitle">
          Quill is an AI-powered writing IDE that helps you craft rich, consistent
          fiction — from the first character sketch to the final chapter.
        </p>
        <button
          className={styles.heroCta}
          data-anim="cta"
          onClick={() => navigate('/login')}
        >
          Try the Demo
        </button>
        <div className={styles.heroOrnament} data-anim="ornament" />
      </section>

      {/* Features */}
      <section className={styles.features} id="features">
        <p className={styles.featuresLabel}>What Quill offers</p>
        <h2 className={styles.featuresTitle}>Everything a writer needs</h2>
        <div className={styles.featureGrid}>
          {features.map((f, i) => (
            <div
              key={f.title}
              className={styles.featureCard}
              ref={(el) => { featureCardsRef.current[i] = el }}
            >
              <div className={styles.featureIcon} aria-hidden>
                {['🖋️', '🤖', '🔗', '📜'][i]}
              </div>
              <h3 className={styles.featureCardTitle}>{f.title}</h3>
              <p className={styles.featureCardText}>{f.text}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Closing */}
      <section className={styles.closing} id="closing">
        <p className={styles.closingQuote} data-anim="quote">
          "A writer is a world trapped in a person."
        </p>
        <p className={styles.closingAuthor} data-anim="author">— Victor Hugo</p>
        <button className={styles.closingCta} onClick={() => navigate('/login')}>
          Start Your Universe
        </button>
      </section>

      {/* Footer */}
      <footer className={styles.footer}>
        <p className={styles.footerText}>Quill — Crafted with ink and intelligence</p>
      </footer>
    </div>
  )
}
