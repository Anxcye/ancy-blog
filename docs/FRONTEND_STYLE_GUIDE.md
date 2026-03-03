# Frontend Style Guide (Blog)

## Purpose
Define a consistent visual and motion system for the public blog frontend (`frontend-blog/`), aligned with "静谧美学" and modern minimal Japanese-inspired design.

## Visual Direction
- Tone: restrained, quiet, editorial, generous whitespace.
- Base background: `#f8f8f8`.
- Surface feeling: "paper floating above page" with soft elevation and thin borders.
- Typography: elegant serif + neutral sans pairing, moderate line-height, no heavy decorative effects.

## Color System
- Use CSS variables as the single source of truth.
- One random theme accent is generated at app bootstrap and injected to root variables.
- Required variables:
  - `--bg`, `--surface`, `--text`, `--muted`
  - `--theme-h`, `--theme-s`, `--theme-l`
  - `--theme`, `--theme-soft`, `--theme-ring`
- Keep saturation controlled; avoid high-contrast neon palettes.

## Background & Atmosphere
- Add subtle grain/noise texture on base background.
- Add floating particle layer (petal-like points/shapes).
- Particle tint and glow must derive from `--theme`.
- Particle density must be lower on mobile to preserve readability and FPS.

## Motion Rules
- All interactive motion uses spring-like easing with slight overshoot.
- Recommended easing tokens:
  - `--ease-spring: cubic-bezier(0.22, 1.25, 0.32, 1)`
  - `--ease-soft: cubic-bezier(0.2, 0.8, 0.2, 1)`
- Entry animation standard:
  - `slide-up-spring` + opacity fade.
  - Stagger children by `40-80ms`.
- Keep animations short and calm; avoid aggressive parallax.

## Layout Rules
- Desktop: large margins, clear reading column, deliberate negative space.
- Mobile: reduce decorative density, keep tap targets >= `40px`, maintain vertical rhythm.
- Components should feel lightweight and breathable, not crowded.

## Accessibility & Performance
- Respect `prefers-reduced-motion` and disable particle/overshoot when enabled.
- Ensure text contrast meets WCAG AA.
- Particle animation should run on transform/opacity and be frame-budget friendly.

## Implementation Note
- This guide is mandatory for all new `frontend-blog` pages/components.
- Any style direction change must update this file and be logged in `docs/SESSION_LOG.md`.
