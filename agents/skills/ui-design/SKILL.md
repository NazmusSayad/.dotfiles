---
name: ui-design
description: Designs, styles, redesigns, and reviews web interfaces. Use for product UI, websites, components, layout, visual hierarchy, typography, color, responsive behavior, interaction, motion, and visual polish.
---

## Understand the work

Before editing, inspect the existing components, tokens, assets, and a representative screen. Decide whether the task is new work, an extension, a redesign, or a review.

Identify the audience, primary task, content, usage frequency, and constraints. Ask a question only when the answer would materially change the result.

Preserve established behavior, content, and visual language unless the user asked to replace them. A local addition to an existing interface is not an excuse to invent a new design system.

Treat product UI and marketing pages differently. Product UI prioritizes clarity, state, and repeated use. Marketing pages may use stronger composition and expression, but the offer and primary action must remain obvious.

## Choose a direction

Derive the visual direction from the subject, audience, and available material. For new work, choose one clear direction and one memorable idea. Keep everything else disciplined. Do not present a menu of aesthetics unless the user asked to explore alternatives.

Use real product content and assets when available. Do not invent claims, metrics, testimonials, customers, product screens, or brand marks. Use official logos and assets rather than recreating them.

## Build the interface

- Establish the reading order and hierarchy before decorating.
- Let content and interaction determine the layout. Do not force every section into cards, equal columns, or a centered composition.
- Use a deliberate type scale, comfortable line lengths, and typefaces appropriate to the subject. Typography should carry hierarchy without decorative labels above every heading.
- Give color clear roles. Maintain readable contrast and avoid gradients, glows, or glass effects unless the direction calls for them.
- Use spacing to show relationships: tight within groups, generous between groups, and more space above a heading than below it.
- Use familiar icons when visual recognition is faster than text or when they improve scanning, orientation, or feedback. Use one coherent icon family. Do not decorate every label or action with an icon, and do not use an ambiguous symbol where clear text is better. Give icon-only controls an accessible name.
- Reuse existing components and tokens. Do not create a parallel visual system for one feature.
- Include hover, focus, active, disabled, loading, empty, and error states when relevant.
- Make every layout work at narrow and wide widths. Resolve overflow, wrapping, touch targets, and content extremes explicitly.
- Prefer semantic elements and visible keyboard focus. Accessibility and responsiveness are part of the design, not a later pass.

## Use motion deliberately

Before adding motion, name its purpose: feedback, state change, spatial continuity, explanation, or preventing a jarring change. If it has no purpose, do not animate it. Frequent and keyboard-driven actions should be instant or nearly instant.

Use the simplest suitable tool. Prefer CSS transitions for small state changes and add a motion library only for gestures, springs, layout transitions, or coordinated exits. Animate transform and opacity when possible, name transition properties explicitly, and keep routine UI motion under 300ms. Entrances and exits usually use ease-out; movement within the screen usually uses ease-in-out.

Ship reduced-motion behavior with the animation. Gate hover effects to devices that support hover.

## Avoid AI design slop

Do not place pills, badges, icons, or eyebrow labels above headings or titles.

Do not repeat the same or similar design patterns over and over. Do not add low-quality elements just to add variety.

Do not reuse the same landing-page composition across pages, especially repeated centered sections or one-sided layouts that leave the other side empty. Do not create variety by right-aligning English text.

Do not wrap every section or content item in a card, or nest cards inside cards.

Do not leave side-by-side columns or cards with visibly unbalanced heights or large empty areas beneath one side.

Do not use large explanatory text blocks when a small diagram would communicate the same information more clearly.

Do not number sections, headings, cards, or content items by default or use numbers as decoration. Use numbering only when it serves a clear purpose for the user.

Do not fill sections with repeated "Learn more" buttons or buttons that merely jump to another section.

Do not fabricate product interfaces or substitute simulated screens for the actual product UI.

Do not create, draw, or hand-code brand logos from scratch instead of using the real logo from an official source or a library.

Avoid generic gradient decoration and identical reveal effects on every section. Any pattern listed here is acceptable when the brief or content genuinely requires it, except fabricating product interfaces or brand assets.
