# Quill Design System

## Core Concept
"Writing a universe as if it were in an ancient manuscript, powered by modern intelligence."

This UI is not a typical application. It is a narrative space.

---

## Color System

Base:
- Ink: #262522   Paper: #FBF6F1   Card: #EBE4D6   Workspace: #F1E5D1

Accents:
- Gold: #BF9969   Sand: #D9BD9C   Stone: #D9D0C7

Semantic:
- Danger #9E3B2E   Warning #B37E2A   Success/Progress #5B7A4A

Knowledge-graph node types:
- Character #7A5C86   Place #6E8B4E   Event #C07B3A   Faction #A8894E   WorldRule #4E8073   PlotArc #9C4A38

Rules:
- All colors live as `:root` tokens in `index.css`; components consume `var(--token)` only — never hardcode hex.
- Opacity variants (`--ink-80`..`-10`, `--paper-60` etc.) for tints instead of new greys.

---

## Typography

- Titles: EB Garamond (`--serif`)
- Body: Lora (`--sans`)
- Marginal/microcopy: Caveat (`--micro`)

Cursor:
- Use a feather/quill style blinking cursor instead of a line

---

## Surfaces

- Paper textures with subtle grain
- Cards feel like manuscript pages
- Sidebar darker, like ink or leather

---

## Layout

- Editorial spacing
- Generous margins
- Not rigid like SaaS dashboards

---

## Components

Cards:
- Look like paper fragments
- Include illustrations and metadata

Buttons:
- Outline or ink-filled
- Hover feels like “inking”

Inputs:
- Underlined style
- Writing animation on focus

Progress:
- Ink strokes instead of modern bars

---

## Animations

- Slow and intentional
- Organic motion
- Fade and subtle movement

---

## Illustrations

- Black and white
- Engraving/sketch style
- No modern or colorful graphics

---

## Emotional Tone

- Calm
- Creative control
- Intimate
- Thoughtful

---

## Consistency Rules

Ask:
1. Does it feel written or generated?
2. Does it feel physical?
3. Does it have space?
4. Could it exist in an ancient book?

If not, redesign.

---
## Animations System (GSAP + ScrollTrigger)

### Technology Stack

The application uses:

- GSAP (GreenSock Animation Platform)  
- ScrollTrigger plugin  

These tools enable precise, performant, and scroll-driven animations aligned with the narrative nature of the product.

---

### Core Principle

Animations are not decorative — they are part of the storytelling experience.

The interface should feel like it is being written, revealed, and explored as the user scrolls.

---

### Animation Approach

All animations are:

- Scroll-driven (linked to user movement)  
- Smooth and slow-paced  
- Subtle and intentional  
- Non-distracting  

Avoid:

- Instant or aggressive transitions  
- Looping animations without user interaction  
- GIFs or video-based animations  

---

### Supported Animation Types

#### Fade + Vertical Reveal  
Elements appear with:
- Opacity: `0 → 1`  
- Transform: `translateY(20px → 0)`  

Used for:
- Sections  
- Cards  
- Text blocks  

---

#### Writing Effect (Text)  
Text appears progressively, simulating writing.

Used for:
- Headlines  
- Descriptions  
- Narrative content  

---

#### SVG Drawing (Ink Effect)  
SVG paths animate using:

- `stroke-dasharray`  
- `stroke-dashoffset`  

Creates a hand-drawn ink effect.

Used for:
- Dividers  
- Icons  
- Illustrations  

---

#### Layered Reveal  
Elements appear in sequence:

1. Title  
2. Divider (draws itself)  
3. Text (writes in)  
4. Illustration (fades in)  

---

#### Subtle Parallax  
Background elements move slower than foreground content.

Used for:
- Decorative elements (mountains, maps, etc.)  
- Background depth  

---

### Interaction Rule

Every animation must respond to scroll.

If an animation does not react to user input, it should not exist.

---

### Performance Guidelines

- Use lightweight assets (SVG, optimized PNG)  
- Avoid large video or GIF files  
- Prefer `opacity` and `transform` for animations  
- Limit simultaneous animations  

---

### Implementation Note

ScrollTrigger controls:

- When animations start  
- Their duration relative to scroll  
- Synchronization between elements  

---

### Summary

Animations should feel like:

- Ink appearing on paper  
- Pages being revealed  
- A story unfolding  

The user is not navigating a UI —  
they are writing and discovering a universe.

## Summary

An editorial, manuscript-inspired interface where modern UI disappears behind a timeless writing experience.
