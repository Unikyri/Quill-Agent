# Quill — Image Asset List

This is the deliverable requested after the frontend restructure: what images to produce, how to name them, and where to save them so they slot into the existing gold/sepia "ancient manuscript" system.

**Status note**: `frontend/public/images/{decorative,entities,landing}/` already exist as empty folders (just `.gitkeep`) from earlier scaffolding — nothing in the code currently references any file inside them yet. Dropping a file in the right place with the right name is **not enough on its own** to make it appear; each image below needs one small `<img src="/images/...">` or CSS `background-image` line added to its page/component. That wiring is NOT done yet — it's a small follow-up (a few lines per page), separate from this list. Say the word and I'll do that pass once real files exist, or now with placeholder-safe fallbacks if you want it wired ahead of time.

## Art direction (applies to every image)

- **Style**: hand-drawn pencil/ink line art — thin, slightly irregular linework, like a woodcut or manuscript marginalia illustration. Not flat vector, not photographic.
- **Palette**: monochrome or duotone, using the ink tone `#262522` as the line/shadow color against transparency — think "line drawing on aged paper," not full-color illustration. This is what lets one image work across both light card backgrounds (`#EBE4D6`/`#F1E5D1`) and darker sidebar backgrounds (`#262522`) if reused.
- **Background**: transparent PNG (no white/paper background baked in) — the app's own paper/card background shows through.
- **Format**: PNG with alpha transparency. SVG is also fine and often better for line art (scales cleanly, smaller file) if your tool can export it — either format works with the same file names below (just change the extension).
- **Edge effect**: the "torn paper" edge mentioned isn't baked into these illustration assets — that's a CSS effect (clip-path or a mask image) applied to containers, not something each individual image needs. Don't spend illustration effort on it.

## Save location convention

Base folder: `frontend/public/images/`. Referenced in code as `/images/<subfolder>/<filename>` (Vite serves `public/` at the root).

| Subfolder | Use |
|---|---|
| `landing/` | Landing page hero illustrations, logo |
| `login/` | Login/Register page background + small logo |
| `dashboard/` | Universe-card watermarks, empty-state decoration |
| `work/` | Per-chapter thumbnail illustrations |
| `entities/` | Character/entity portrait illustrations |
| `decorative/` | Shared/reusable flourishes not tied to one page |

## Asset list

### Landing page → `frontend/public/images/landing/`

| Filename | Subject | Notes |
|---|---|---|
| `logo-color.png` | Quill's logo mark — a feather/quill pen icon, in the gold `#BF9969` accent color | Small, used in navbar; needs a crisp version at ~64×64 and up (SVG ideal) |
| `hero-knight-dragon.png` | Full-bleed hero illustration: a knight facing a dragon, fantasy adventure theme | Wide aspect (matches a full-width hero section), transparent background so the paper texture shows through at the edges |
| `hero-heroic.png` | Full-bleed hero illustration: a general "heroic journey" scene (protagonist silhouette, dramatic pose) | Same wide aspect as above |
| `hero-astronaut.png` | Full-bleed hero illustration: an astronaut/space scene | Represents the sci-fi genre track; same wide aspect |

### Login/Register → `frontend/public/images/login/`

| Filename | Subject | Notes |
|---|---|---|
| `background-elements.png` | Full-page decorative background — subtle scattered manuscript motifs (feathers, ink blots, compass rose) | Very low visual weight — sits behind the login card, must not compete with form legibility. Consider a low-opacity single-tone treatment. |
| `logo-bw.png` | Small black-and-white version of the Quill logo | ~55×55px usage, above the login card |

### Dashboard → `frontend/public/images/dashboard/`

| Filename | Subject | Notes |
|---|---|---|
| `inkwell.png` | An inkwell + quill illustration | Watermark on the "New Universe" empty-state card |
| `castle.png` | A castle illustration (fantasy genre marker) | Watermark on fantasy-genre universe cards |
| `spaceship.png` | A spaceship illustration (sci-fi genre marker) | Watermark on sci-fi-genre universe cards |
| `mountains.png` | A mountain range illustration | Bottom-right decorative flourish over the dashboard's empty background, low opacity (~30%) |

Since universe genre is user-set data (not a fixed enum tied one-to-one to these 3), consider a small set covering the genres your app actually supports, following this same naming pattern (`<genre>.png`) — the 3 above cover fantasy/sci-fi/general; add more (e.g. `mystery.png`, `romance.png`) only if those genres are real options in the universe-creation form.

### Work page → `frontend/public/images/work/`

Per-chapter thumbnail illustrations — small, varied so a chapter list doesn't look repetitive:

| Filename | Subject |
|---|---|
| `chapter-mountains.png` | Mountain landscape thumbnail |
| `chapter-castle-night.png` | Castle at night thumbnail |
| `chapter-ships-harbor.png` | Ships in a harbor thumbnail |
| `chapter-untitled.png` | A generic "?" placeholder mark, for chapters without enough content yet to pick a thumbnail |

These are meant to be assigned somewhat arbitrarily per chapter (there's no backend field driving "which illustration for which chapter" today) — treat them as a rotating decorative set, not content-accurate art.

### Knowledge Graph / Timeline / Entity Card → `frontend/public/images/entities/`

| Filename | Subject | Notes |
|---|---|---|
| `portrait-placeholder.png` | A generic circular character-portrait silhouette | Used as the default node icon in the Knowledge Graph, the Timeline event icon, and the Entity Card header when no specific portrait exists for that character |

If you eventually want per-character portraits instead of one shared placeholder, name them `portrait-<entity-name-slug>.png` (e.g. `portrait-aragorn.png`) — but that's a bigger content-production task, not something to block the current visual pass on.

## What I did NOT include

- Any image for Notes, AI Assistant chat, or the "Orphaned" tab — those features don't exist in the backend and were explicitly excluded from this restructure.
- A torn-paper-edge texture image — that's a CSS effect on containers, not a per-page illustration asset (see Art Direction above).

## Next step, if you want it

Once you have real files at these paths, tell me and I'll wire the actual `<img>`/CSS references into the corresponding page components (Dashboard cards, Login background, Landing hero sections, chapter thumbnails, Knowledge Graph/Timeline/EntityCard portrait). Small change, a few lines per page — I'd run it through the same SDD flow as everything else in this session.
