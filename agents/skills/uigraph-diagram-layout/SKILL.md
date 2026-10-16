---
name: uigraph-diagram-layout
description: Hand-layout and styling of diagrams in the UIGraph viewer (`npx uigraph`, React Flow editor at http://127.0.0.1:7357) driven through playwright-cli. Use when generated UIGraph diagrams look mangled, overlapping, unbalanced or harshly colored and need to be repositioned, resized, recolored or re-routed so they look good, for example before recording them for a video.
---

## The standard

The goal is a diagram a person finds clean at a glance: clear reading direction, aligned boxes, straight or simply-jogged lines, readable labels, and calm color. Layout is judged by eye, not computed. Users reject scripted auto-layout: plan positions yourself, move nodes in the live browser, and judge every result from a screenshot taken after a reload. Scripts exist only to make precise mouse input possible; they never decide where something goes.

Diagrams live in `<project>/.uigraph/diagrams/<slug>/` as `<slug>.mmd` (source) plus `context.json` (positions, sizes, styles, edge handles). The viewer autosaves every edit to both files.

## Workflow

1. Open the diagram (`/diagram/<project>--<slug>`) in a headless 1600x1000 session and take a `-before.png` screenshot. Read the `.mmd` and run `scripts/layout.sh` to get node ids, sizes and edge ids.
2. Decide the layout on paper before touching anything: which tiers exist (callers → core → stores), what sits in the center, and which side of each box every edge should use. Write target flow coordinates for every node and check them against the viewer rules below.
3. Style nodes, then move them with `scripts/flowto.sh`, then reconnect edges whose saved handles no longer face their target (`scripts/reconnect.sh`), then style edges.
4. Reload the page, re-run `layout.sh`, and screenshot. Positions that changed after reload mean a viewer rule pushed them; respace and repeat. The diagram is done when it survives a reload unchanged and the screenshot looks clean.
5. Save `-after.png` next to the before shot.

## Layout judgment

- Pick one reading direction. Flow and job chains read top to bottom with triggers in a row at the top. Hub diagrams put the hub in the center, callers on the left, the hub's dependencies above and below, and external systems in a column on the right.
- Align on centers, not edges, so straight edges really are straight: a node feeding straight down shares the target's center x; side-by-side neighbors share center y. Give siblings the same width and spacing, and keep mirrored branches symmetric around the parent's center.
- Each edge should leave and enter on the facing sides (source-bottom → target-top for something below). Design so each node has at most one edge per side when that edge should be straight; several edges fanning out of one side are fine when they diverge anyway.
- Leave room for labels. A smoothstep label sits at the middle of the path, so keep the gap between tiers large enough (≈150px+) that labels clear both boxes and any boundary line. Two labels on the same horizontal run must not overlap.
- Nothing may cross a box. Reroute or reorder before accepting a line through a node.
- Size boxes to their text with even padding; widen a box rather than let a label wrap awkwardly. Widths that differ by a few pixels look like mistakes. Use equal widths or clearly different ones.

## Viewer rules that move things behind your back

These come from the viewer bundle (`~/.npm/_npx/*/node_modules/uigraph/dist/viewer/assets/index-*.js`). Recheck the constants there if behavior differs.

- **Row push on load.** Nodes are sorted by y. A node whose y is within 80px of a row's first node joins that row, and each node in a row must start at least `previous.x + previous.width + 140`, or it gets pushed right on reload. Either keep rows ≥ 81px apart in y or ≥ 140px apart in x, with a safety margin of 10px or more, because drags land ±1px.
- **Hub offset.** Edges on the same *geometric* side of a node (decided by comparing center distances, `|dx| > |dy|` means left/right, regardless of the saved handle) are spread 24px apart (±12 for two). An intended straight line then jogs by 12px. Compensate by shifting the other node 12px, or avoid sharing that side.
- **Handles are saved per edge and are not reassigned when nodes move.** After moving, edges keep exiting old sides and loop around. Reconnect them.
- **C4 boundaries.**
  - Children are clamped inside the boundary while dragging, and the boundary does not grow on its own. Resize it first with `scripts/resize-boundary.sh`, then move children, then shrink it to fit (about 40px padding, more at the bottom for the boundary label).
  - Dropping a node with its center inside a boundary adopts it, and the autosave rewrites the `.mmd` to put it inside the boundary block. Once adopted, a node can't be dragged back out. Keep outside nodes clear of the boundary while moving. If one gets adopted, navigate away, move its line back out in the `.mmd`, convert its `___position` to absolute (child position + boundary position), and tell the user you touched the source.
- **Sequence diagrams are laid out automatically from the mermaid source.** Participants and messages can't be positioned. Message width is about `max(140, min(240, 24 + 7 * max(ceil(len/2), longestWord)))`, so a single long token (e.g. `StorageTemplateMigrationSingle`) breaks mid-word. Fixing them means editing the `.mmd` wording, which changes the source; tell the user.
- Saving drops `name`/`description` keys from `context.json`.
- The viewer opens zoomed in tightly, with top nodes hidden under the toolbar. Zoom out a step for screenshots and recordings.

## Color

Users found saturated fills (`#1e3a8a` with `#60a5fa` borders and so on) "brutal". Aim for subtle: nearly-background fills with a faint hue, muted mid-tone borders, off-white text, and 1.5px lines. Color encodes role, and an edge takes the color of the node it points to.

| Role | Fill | Border | Text |
|---|---|---|---|
| Core service | `#1a2742` | `#5a7fb8` | `#e2e8f0` |
| Clients / triggers / people | `#1d2330` | `#6b7689` | `#e2e8f0` |
| ML / async workers | `#241f3a` | `#8577b8` | `#ddd8ee` |
| Data stores / sync | `#152b2c` | `#4b9189` | `#d3e6e3` |
| External systems / notifications (dashed) | `#211e1a` | `#a88848` | `#e6dcc2` |
| Delete / destructive | `#2f1d24` | `#a8636f` | `#ecd9dc` |

C4 nodes only expose Color (border tint) and Text Color (`scripts/style-c4.sh`). Their default dark-blue text is hard to read on the dark canvas, so set the text color. When the user says to only reorganize a diagram, leave its colors as they are.

Pure color-value changes across finished diagrams can be done by remapping hex values in `context.json`, but only while no browser tab has that diagram open, because an open tab autosaves its in-memory state over the file. Back the files up first, then reload and screenshot to confirm. Never edit positions that way.

## Driving the viewer reliably

All scripts use the session `${UIG_SESSION:-uig}` and take flow coordinates.

- **Style panel inputs** commit only on click, Meta+A, type, Tab, with ~350ms pauses. `fill()` silently saves wrong values.
  - Node fields: `getByText(label).locator('xpath=../..')`, then the visible input, skipping range sliders and hidden inputs.
  - Edge fields: `locator('label').getByText(label).locator('xpath=..')`; Stroke Width needs `../..`.
  - A loose selector can write a color into the Edge Label field; check labels afterwards.
- **Select an edge** with `dispatchEvent('click')` on its `.react-flow__edge-interaction` path. Edge ids look like `edge-<src>-<dst>-<n>` or, for C4, `c4-rel-<src>-<dst>-<n>`.
- **Reconnect edges** by dragging `.react-flow__edgeupdater-source|target` onto `.react-flow__handle[data-handleid=...]`. Dragging the edge itself moves the node instead.
- **Moving nodes.** Grab a node near its left edge (`x + 8`), not the center label, which enters text edit. Drags slip ~4 screen px along the initial move; `flowto.sh` compensates and iterates. Residuals under ~1px can't be corrected at low zoom; zoom in for final alignment of straight edges.
- **Zoom and visibility.** Ctrl+wheel zooms and a plain wheel pans; arrow keys don't nudge nodes. Off-screen nodes can't be clicked, so zoom out (`scripts/zoom.sh 0.5`) before batch operations.
- **Batch long operations** in groups of a few nodes; one playwright `run-code` call times out at 120s.
