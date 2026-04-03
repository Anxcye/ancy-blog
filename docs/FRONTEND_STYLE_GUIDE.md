# Frontend Style Guide

## Scope
This guide governs both frontend apps:
- `frontend-blog/`: public reading experience.
- `frontend-admin-react/`: content and system management.

## Admin Design Direction
- Tone: clean, calm, efficient, and lightweight.
- Primary accent for admin is fixed fresh teal:
  - `--accent: #2AA889`
  - `--accent-hover: #23967A`
  - `--accent-soft: #E8F7F3`
- Text/background baseline:
  - `--text: #1A1A1A`
  - `--muted: #667085`
  - `--bg: #F7FAF9`
  - `--surface: #FFFFFF`
- Keep visual hierarchy simple; reduce heavy shadows and noisy decorations.

## Public Blog Visual Direction
- Tone: restrained, quiet, editorial, generous whitespace.
- Base background: `#f8f8f8` with subtle grain texture.
- Accent system: one random accent per page load from a soft palette, injected via CSS variable (`--accent-color`).
- Gradient: a very light accent gradient appears near the top-lower area, fading into background.

## Hero & Navigation
- Home top section is a full `100vh` hero.
- Header layout: left small avatar, centered nav, right day/night switch.
- Public blog header should use a layered top-down progressive backdrop blur behind the fixed shell, so the area closest to the viewport top is most blurred and fades into clear content below.
- Hero body split:
  - left: intro text from site settings + social links row;
  - right: profile avatar.
- Add a softly floating down-arrow at hero bottom.
- Navigation wording should be literary and brand-toned, not generic labels.
- Navigation display copy may be more literary than page titles, but should remain understandable on first read.
- Moments page should read as a quiet vertical note stream, not a dashboard timeline or masonry grid.
- Timeline page should behave like grouped editorial chronology, not a generic product timeline:
  - group entries by month or another clear time anchor;
  - mix article and moment rows in one stream;
  - let moment rows reuse the shared moment detail modal instead of inventing a second detail pattern.
  - default to collapsed month groups except for the newest group;
  - expanded entries should stay compact, closer to an archive index than to a card feed.

## Motion & Particles
- Use spring-like easing for interactions and page entry.
- Entry pattern: `slide-up + fade` with `40-60ms` stagger.
- Hover/click: mild lift and press feedback; keep calm and short (`200-400ms`).
- Particle layer: sparse petal-like particles, very low density (about one visible fall event per ~3s), tinted by current accent.
- Comment composer triggers can use stronger spring feedback than article cards, but still need a restrained editorial shell.
- Comment thread lists may use responsive masonry/waterfall stacking when cards have uneven heights.
- Moments detail should open in a route-driven modal instead of expanding inline inside the feed.
- The moments modal should align near the top of the viewport, lock background scroll, and keep wheel focus inside the dialog.
- Moments feed items may use subtle pointer-follow background motion on hover, but should not show explicit "view detail" buttons.
- Public list/info pages such as `articles`, `moments`, `timeline`, and `friends` should share a consistent hero shell:
  - small uppercase eyebrow;
  - large editorial title;
  - restrained explanatory subtitle;
  - optional compact stat pills aligned to the right on wide screens and stacked on mobile.
- Friends page should read as a neighbor directory, not as a form-first utility page:
  - place any editorial intro ahead of the directory list;
  - prefer light roster rows over heavy grid cards;
  - keep the submission form collapsed behind a clear CTA by default.
- Gallery masonry cards should use a subtle image zoom on hover and compact bottom metadata text over a dark gradient, with no bordered or boxed overlay container.
- Gallery masonry stream should use tight inter-photo spacing and square-corner photo tiles so the stream reads more like a continuous photo wall than separate cards.
- Gallery photo viewer details should use typography, spacing, and font-weight hierarchy for section separation, not border lines or card-like background blocks in the right panel.
- Gallery photo viewer images should stay shadowless and square-cornered to preserve the original photo silhouette.
- Gallery photo viewer should keep image edge spacing tight so the photo occupies more of the viewport.
- Gallery photo viewer stage should fill the available main area, and the image itself should use `object-fit: contain` inside that full-size stage to avoid accidental shrink-wrapping.
- Gallery photo viewer should use only the full-screen BlurHash wash behind the photo, not a second BlurHash layer inside the main image stage.
- Gallery photo viewer should render the current photo's BlurHash as a full-screen blurred backdrop, and place metadata on a floating translucent frosted-glass panel with generous outer margins, rounded corners, and a desktop collapse toggle so the image can reclaim more viewport space.
- Gallery routes (`/gallery` and locale-prefixed variants) should force dark mode while active, then restore the user's previous theme preference after navigating away, so the gallery consistently reads on a black cinematic background without permanently changing site theme.
- Desktop gallery viewer actions should be visually grouped inside the floating metadata panel header when the panel is expanded; external chrome should be kept minimal and only appear as a restore control after collapse.
- In the gallery viewer, keep the global close action as a floating top-left window control outside the metadata panel; the desktop metadata panel header should only carry the collapse action on the right.
- On mobile gallery viewer, avoid putting action buttons inside the bottom-sheet metadata panel; use external window-level controls and gesture/tap collapse instead.

## Responsive & SEO Constraints
- Mobile-first layout and touch targets >= `40px`.
- Admin and blog must both provide strong mobile adaptation.
- i18n URL strategy: Chinese as default root (`/`), English under `/en/*`.
- Ensure canonical, hreflang, meta, and structured SEO output for both locales.

## Admin Mobile Full-Feature Rule
- Mobile admin is full-featured, not a reduced companion mode.
- Required mobile shell pattern:
  - Top fixed app bar for title/context actions.
  - Bottom fixed tab bar with five primary entries: `Workbench`, `Content`, `Site`, `Interaction`, `System`.
- Editing pages on mobile must support:
  - Content editing.
  - Metadata editing.
  - Preview switching.
  - Draft/publish/schedule actions.
- Do not hide core admin capabilities behind desktop-only interactions.

## Implementation Rules
- Keep style tokens in CSS variables; avoid hardcoded accent colors.
- Respect `prefers-reduced-motion` for particle and spring effects.
- Any style-direction change must update this file and be logged in `docs/SESSION_LOG.md`.

- Homepage hero social links should use circular icon buttons rather than plain glyph links.
- Homepage hero secondary line should stay on one line on wider screens by scaling typography down with available width; mobile may wrap naturally.
- Article detail metadata should read as one compact editorial line:
  - publish time first;
  - update time only when it meaningfully differs from publish time;
  - category inline with the same shell;
  - use icon-led inline items instead of capsule chips;
  - AI disclosure shown as an inline item with hover explanation.
