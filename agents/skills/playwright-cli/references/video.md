## Demo videos

For a quick recording of steps you run as commands, use `video-start`, `video-chapter`, `video-show-actions` and `video-stop`. For a polished demo with overlays, timed pauses, and one continuous take, walk the flow with `playwright-cli` to collect locators, then run the whole recording as one `run-code` function. It uses Playwright's `page.screencast` API:

| Method                                                      | Use                                                                                                                                                       |
| ----------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `start({ path, size, fps })` / `stop()`                     | Record to WebM; convert to `.mp4` afterwards                                                                                                              |
| `showActions({ cursor, duration, position, style })`        | Animated cursor travelling to each action, plus optional click point, target highlight, and action title. Actions are paced by `duration` (default 500ms) |
| `hideActions()`                                             | Stop annotating actions and hide the cursor                                                                                                               |
| `showChapter(title, { description, duration, styleSheet })` | Full-screen card over a blurred page; blocks until `duration` ends                                                                                        |
| `showOverlay(html, { duration })`                           | Custom HTML overlay; without `duration` it stays until `.dispose()`                                                                                       |
| `hideOverlays()` / `showOverlays()`                         | Temporarily hide or show all overlays                                                                                                                     |

- `showActions` `cursor` is `'pointer'` (default) or `'none'`. `position` is `top-left|top|top-right|bottom-left|bottom|bottom-right`.
- `style.point`, `style.highlight` and `style.title` are CSS declaration strings. `point` and `highlight` are hidden unless set. `point` is zero-sized and centred on the click point, so give it a size. Prefer `outline` over `border` for `highlight`. Use `title: 'display: none'` to keep the cursor without the callout.
- Overlays are `pointer-events: none`, so sticky overlays can stay visible while the code clicks and types.
- Use `pressSequentially(text, { delay: 60 })` for natural typing and `waitForTimeout(500)` to `waitForTimeout(1500)` between steps.
- Give `start` an absolute `.webm` path under `.playwright/<group>/`, then convert it to `.mp4` with `ffmpeg` and delete the `.webm`.

```bash
playwright-cli run-code "$(cat <<'EOF'
async (page) => {
  await page.screencast.start({
    path: "/path/to/project/.playwright/todo/add-item.webm",
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
  const input = page.getByRole("textbox", { name: "What needs to be done?" })
  await input.pressSequentially("Walk the dog", { delay: 60 })
  await input.press("Enter")
  const note = await page.screencast.showOverlay(
    '<div style="position:absolute;top:8px;right:8px;padding:6px 12px;background:rgba(0,0,0,.7);color:white;border-radius:8px">Item added</div>'
  )
  await page.waitForTimeout(1500)
  await note.dispose()
  await page.screencast.stop()
}
EOF
)"
```
