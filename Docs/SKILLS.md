# Quill — Skills Catalogue & Genre Vocabulary

**Companion documents:** [PRD.md](./PRD.md) · [SRS.md](./SRS.md)
**Scope:** the closed genre vocabulary, and the 15 Skills Quill ships.

> **Status: AUTHORED.** All 15 skill bodies and all 20 genre reference files exist in
> `backend/skills/` (35 files). This document is the catalogue and design rationale;
> the files themselves are the source of truth. **Do not re-author them** — implementation
> work consumes them as-is (loading, selection, activation: SRS §7).

---

## 1. What a Skill is (and is not)

A **Skill** is a read-only Markdown file carrying **craft** knowledge. It is curated by us,
activated by the writer at the **universe** level, and **selected by the agent** at review time.

A Skill **is not**:

- **Story memory.** Continuity, contradictions, plot holes, and timeline validation are the job
  of the memory pipeline, which already runs autonomously in the background. There is
  deliberately **no "Continuity Checker" skill** — duplicating that engine would confuse where
  the intelligence actually lives.
- **A prose generator.** No skill rewrites the author's text. Skills produce **margin
  observations**. Quill is the editor, not the author.
- **User-editable.** A writer is a writer, not an engineer.

### 1.1 File format

```yaml
---
name: line-editor
description: >
  Reviews prose at the sentence and paragraph level — rhythm, flow, clarity, word choice,
  and voice. Use this whenever the writer asks about how a passage *reads*, whether prose
  is clunky, whether sentences drag, or whenever a passage feels flat, wordy, or awkward,
  even if they don't name "line editing" explicitly.
genre_tags: []          # empty = applies to every genre
stage: line             # developmental | line | copy | proof | reader | market | craft
---

<instruction body>
```

**The `description` is written for the agent, not for a human.** It is the selection criterion:
the agent reads every active skill's `description` and decides which apply to the passage in
front of it. Descriptions should lean **slightly pushy** — the failure mode in practice is
*under*-triggering, not over-triggering.

### 1.2 Progressive disclosure

1. Only the `description` lines of active skills enter the **selection** prompt.
2. The model picks the applicable skill(s).
3. Only the **selected** skill's full body is loaded into the **review** prompt.

Loading every active skill body would defeat the purpose. This is the token optimisation.

---

## 2. Genre Vocabulary (closed, multi-tag)

Twenty tags. A universe carries **one or more**. They combine freely —
`horror + romance + gothic` is a valid and common universe.

| Tag | Tag | Tag | Tag |
|---|---|---|---|
| `fantasy` | `epic-fantasy` | `urban-fantasy` | `romantasy` |
| `science-fiction` | `space-opera` | `dystopian` | `horror` |
| `gothic` | `paranormal` | `romance` | `mystery` |
| `cozy-mystery` | `thriller` | `crime` | `historical` |
| `literary` | `adventure` | `young-adult` | `coming-of-age` |

**The vocabulary is closed and server-validated.** The LLM may not introduce a tag. If it
could, three books later the system would hold `dark fantasy`, `grimdark`, `fantasía oscura`,
and `low fantasy` as four distinct axes, and preference conditioning (PRD §5.4) would
disintegrate — you cannot filter on an axis that mutates.

---

## 3. The Skills

### 3.1 The editorial team (7)

These mirror the real publishing chain: **developmental → line → copy → proofread**, plus the
readers and the gatekeeper. Each looks at something the others do not.

#### 1. `developmental-editor` — *stage: developmental*

> **description:** Reviews the big picture — story structure, act and beat placement, character
> arcs, stakes, scene purpose, and whether the narrative actually goes anywhere. Use this
> whenever the writer asks about structure, pacing across chapters, whether a scene "earns its
> place", whether an arc lands, or when they say a section feels aimless, saggy, or unearned —
> even if they don't use the word "structure".

**Mandate:** the most comprehensive lens. Does this scene exist for a reason? Does the
protagonist want something, and is something in the way? Is the midpoint doing work? Is this
character changing, or just present?

**Explicitly not:** sentence-level prose. That is the line editor.

---

#### 2. `line-editor` — *stage: line*

> **description:** Reviews prose at the sentence and paragraph level — rhythm, flow, clarity,
> word choice, and voice. Use this whenever the writer asks how a passage *reads*, whether the
> prose is clunky or wordy, whether sentences drag, or whenever a passage feels flat or awkward,
> even if they don't name "line editing" explicitly.

**Mandate:** the music of the sentence. Sentence-length variation. Dead verbs propped up by
adverbs. Clauses in the wrong order. Rhythm that fights the mood of the scene.

**Explicitly not:** grammar correctness. A grammatically perfect sentence can still be ugly.

---

#### 3. `copy-editor` — *stage: copy*

> **description:** Reviews grammar, punctuation, dialogue punctuation, capitalisation,
> hyphenation, number style, and internal consistency of names, spellings, and formatting. Use
> this whenever the writer asks about correctness, house style, or consistency — and whenever a
> passage is being prepared for submission or publication.

**Mandate:** correctness and consistency. Dialogue punctuation is the single most commonly
botched thing in an unedited manuscript and deserves specific attention.

**Explicitly not:** story-level consistency (that is the memory pipeline).

---

#### 4. `proofreader` — *stage: proof*

> **description:** Final-pass surface check for typos, doubled words, missing punctuation, and
> spacing errors. Use this only when the writer says a passage is finished, final, or ready to
> submit — not while they are still drafting.

**Mandate:** the last filter. Cheap, narrow, and it must **not** propose rewrites — by
definition, proofreading happens after the text is settled.

---

#### 5. `beta-reader` — *stage: reader*

> **description:** Reacts as a reader, not an editor — where attention drifted, where confusion
> set in, where belief broke, where a reveal landed or fell flat. Use this whenever the writer
> asks "does this work?", "is this boring?", "will they get it?", or wants a gut reaction rather
> than a technical critique.

**Mandate:** the only skill that reports **experience** rather than **craft**. *"I stopped
believing the character here."* *"I skimmed this paragraph."* Judgements are subjective and
must be stated as such — a beta reader says "I felt", not "this is wrong".

---

#### 6. `sensitivity-reader` — *stage: reader*

> **description:** Reviews representation of cultures, disabilities, gender, race, religion, and
> trauma for harmful stereotype, flat caricature, or unexamined cliché. Use this whenever a
> passage depicts a group or experience outside the writer's own, or handles violence, abuse, or
> identity — and whenever the writer asks whether something "lands badly" or is "handled well".

**Mandate:** flag, explain, and offer the *why* — never moralise, never demand a change. The
writer decides. A note that reads as scolding gets dismissed, and a dismissed note teaches the
system nothing.

---

#### 7. `literary-agent` — *stage: market*

> **description:** Reads as an acquiring agent — does the opening hook, is the premise clear and
> distinctive, does the voice stand out, and is this positioned for a real shelf? Use this
> whenever the writer asks whether the work is sellable, publishable, or ready to query, and
> whenever they are working on an opening chapter, a synopsis, or a pitch.

**Mandate:** the commercial gatekeeper. Does page one earn page two? What does this compete
with, and why would a reader choose this one?

---

### 3.2 Craft (7)

#### 8. `pov-and-tense` — *stage: craft*

> **description:** Enforces point-of-view and tense discipline — first, second, or third person;
> limited or omniscient; head-hopping between viewpoints inside a scene; and slips between past
> and present tense. Use this whenever a passage seems to know things the viewpoint character
> cannot know, whenever the narration drifts between heads, or whenever the writer is
> establishing or changing narrative distance.

**Mandate:** the writer declares a POV and tense per work; this skill defends it. Head-hopping
is the most common unconscious break and the hardest for an author to see in their own prose.

---

#### 9. `show-dont-tell` — *stage: craft*

> **description:** Finds narration that *states* emotion, character, or significance where it
> could be *dramatised* through action, sensory detail, or subtext. Use this whenever prose
> announces what a character feels ("she was furious"), summarises what a scene could enact, or
> whenever a passage feels distant or inert.

**Mandate:** and — importantly — the inverse. **Telling is not always wrong.** Summary is a
legitimate tool for compressing time and skipping the boring parts. A skill that flags every
instance of telling is a nuisance. Flag telling where *dramatising would be better*, and say why.

---

#### 10. `dialogue-and-voice` — *stage: craft*

> **description:** Reviews dialogue for naturalness, subtext, and — crucially — whether distinct
> characters actually *sound* distinct. Use this whenever a passage is dialogue-heavy, whenever
> characters risk sounding like each other or like the author, and whenever the writer asks
> about voice, banter, or how a conversation reads aloud.

**Mandate:** the classic amateur failure is that every character shares one voice — the
author's. This skill has memory of the character's established voice (via the entity graph) and
can flag drift.

---

#### 11. `pacing-and-tension` — *stage: craft*

> **description:** Reviews the rhythm of tension — where a scene stalls, where stakes vanish,
> where a chapter ends on a soft beat instead of a pull, and where relentless action leaves the
> reader no room to breathe. Use this whenever the writer asks if something drags or rushes, and
> whenever a chapter opens or closes.

**Mandate:** including the opposite failure — unbroken tension is exhausting and self-defeating.
Pacing is a wave, not a slope.

---

#### 12. `worldbuilding-and-exposition` — *stage: craft*

> **description:** Reviews how the world is delivered — flagging infodumps, "as you know, Bob"
> dialogue, and lore delivered ahead of the reader's need for it, as well as the reverse: rules
> and stakes the reader was never given. Use this whenever a passage introduces a magic system,
> a technology, a history, or a place, and whenever the writer asks whether they've explained too
> much or too little.

**Mandate:** the world must be **revealed at the speed of need**. This skill reads the universe's
established lore (via the memory subsystem) and can tell the difference between "the reader
doesn't know this yet" and "the reader has been told this three times".

---

#### 13. `prose-economy` — *stage: craft*

> **description:** Hunts filter words ("she saw", "he felt", "it seemed"), adverb propping,
> hedging qualifiers, redundant beats, and cliché. Use this whenever prose feels padded, soft,
> or second-hand, and whenever the writer wants a passage tightened.

**Mandate:** and it **must respect learned preferences**. If the writer has rejected the adverb
note three times, the Writer Memory has promoted that to a deliberate choice — and this skill
**stops raising it**. This skill is the primary source of feedback signals for Writer Memory
(PRD §5), and the primary beneficiary of them.

---

#### 14. `period-register` — *stage: craft*

> **description:** Checks that vocabulary, idiom, technology, social assumptions, and material
> culture belong to the era the work is set in — flagging anachronistic words, objects, and
> attitudes. Use this whenever a work is historical, medieval, or set in any period other than
> the present, and whenever the writer asks whether something "feels modern".

**genre_tags:** `historical`, `epic-fantasy`, `fantasy`, `gothic`

**Mandate:** the subtle failures are not swords and horses — they are *words*. A medieval
character cannot be "stressed", "traumatised", or "focused". This is where the optional
Wikipedia MCP (PRD §4.5) earns its place, and it must degrade gracefully without it.

---

#### 15. `genre-conventions` — *stage: craft*

> **description:** Applies the conventions, reader expectations, and structural promises of the
> work's genre — what this kind of book must deliver, what it must not, and which beats readers
> of it are waiting for. Use this whenever a passage engages a genre expectation: a magic
> system, a murder, a meet-cute, a scare, a reveal, a heist — and whenever the writer asks
> whether something "fits the genre".

**This is the elegant one.** It is **one skill with twenty reference files** — one per genre tag:

```
genre-conventions/
├── SKILL.md                    (the selection logic + how to apply conventions)
└── references/
    ├── fantasy.md
    ├── epic-fantasy.md
    ├── romantasy.md
    ├── cozy-mystery.md
    └── … (20 total, one per tag)
```

The skill loads **only the reference files matching the universe's genre tags**. A
`horror + gothic` universe loads two files, not twenty.

**Why one skill and not twenty:** twenty genre skills would drown the writer in a selection
list, bloat the selection prompt with twenty descriptions, and duplicate the same application
logic twenty times. One skill with progressive disclosure over references gives identical
capability at a fraction of the context — and it is exactly the domain-organisation pattern
that skill authoring recommends for multi-variant skills.

---

## 4. Suggestion on universe creation

When a universe is created with genre tags, Quill **pre-suggests** the skills whose
`genre_tags` intersect them, plus the universally applicable ones. The writer confirms or
adjusts. They should never face a blank list of fifteen and be asked to guess.

---

## 5. Authoring status

| Artefact | Status |
|---|---|
| Genre vocabulary (20 tags) | **Defined** — ready to implement (migration `021`) |
| Skill frontmatter (`name`, `description`, `genre_tags`, `stage`) for all 15 | **Defined** — this document |
| Skill instruction bodies | **To be authored during implementation** |
| `genre-conventions` reference files (20) | **To be authored during implementation** |
