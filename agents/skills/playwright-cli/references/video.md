## Scripted demo videos

For a polished recording, first walk the flow with `playwright-cli` to collect locators, then write one script and run it with `playwright-cli run-code --filename=demo.js`. A script controls pacing and annotations better than individual commands.

Playwright's `page.screencast` API:

| Method                                                      | Use                                                                                                                                                       |
| ----------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `start({ path, size, fps })` / `stop()`                     | Record to WebM                                                                                                                                            |
| `showActions({ cursor, duration, position, style })`        | Animated cursor travelling to each action, plus optional click point, target highlight, and action title. Actions are paced by `duration` (default 500ms) |
| `hideActions()`                                             | Stop annotating actions and hide the cursor                                                                                                               |
| `showChapter(title, { description, duration, styleSheet })` | Full-screen card over a blurred page; blocks until `duration` ends                                                                                        |
| `showOverlay(html, { duration })`                           | Custom HTML overlay; without `duration` it stays until `.dispose()`                                                                                       |
| `hideOverlays()` / `showOverlays()`                         | Temporarily hide or show all overlays                                                                                                                     |

- `showActions` `cursor` is `'pointer'` (default) or `'none'`. `position` is `top-left|top|top-right|bottom-left|bottom|bottom-right`.
- `style.point`, `style.highlight` and `style.title` are CSS declaration strings. `point` and `highlight` are hidden unless set. `point` is zero-sized and centred on the click point, so give it a size. Prefer `outline` over `border` for `highlight`. Use `title: 'display: none'` to keep the cursor without the callout.
- Overlays are `pointer-events: none`, so sticky overlays can stay visible while the script clicks and types.
- Position an overlay around an element with its `locator.boundingBox()` and absolutely positioned HTML.
- Use `pressSequentially(text, { delay: 60 })` for natural typing and short `waitForTimeout` pauses between steps.

```js
;async (page) => {
  await page.screencast.start({
    path: "demo.webm",
    size: { width: 1280, height: 800 },
    fps: 60
  })
  await page.screencast.showActions({
    duration: 800,
    style: { title: "display: none" }
  })
  await page.goto("https://demo.playwright.dev/todomvc")
  await page.screencast.showChapter("Adding items", {
    description: "Add a todo item.",
    duration: 2000
  })
  await page
    .getByRole("textbox", { name: "What needs to be done?" })
    .pressSequentially("Walk the dog", { delay: 60 })
  await page
    .getByRole("textbox", { name: "What needs to be done?" })
    .press("Enter")
  const note = await page.screencast.showOverlay(
    '<div style="position:absolute;top:8px;right:8px;padding:6px 12px;background:rgba(0,0,0,.7);color:white;border-radius:8px">Item added</div>'
  )
  await page.waitForTimeout(1500)
  await note.dispose()
  await page.screencast.stop()
}
```
