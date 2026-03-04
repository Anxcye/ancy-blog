# Frontend Style Guide

## Scope
This guide governs both frontend apps:
- `frontend-blog/`: public reading experience.
- `frontend-admin/`: content and system management.

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
- Hero body split:
  - left: intro text from site settings + social links row;
  - right: profile avatar.
- Add a softly floating down-arrow at hero bottom.
- Navigation wording should be literary and brand-toned, not generic labels.

## Motion & Particles
- Use spring-like easing for interactions and page entry.
- Entry pattern: `slide-up + fade` with `40-60ms` stagger.
- Hover/click: mild lift and press feedback; keep calm and short (`200-400ms`).
- Particle layer: sparse petal-like particles, very low density (about one visible fall event per ~3s), tinted by current accent.

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
