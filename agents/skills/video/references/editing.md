# Editing

## Footage in

- Takes arrive raw and unclipped: one continuous recording per scene, constant 60 fps H.264. No trimming, speed changes or joins before the edit.
- Keep failed takes and add `-take2`, `-take3`. Don't overwrite.
- Before using a take, check it with `ffprobe` (1920x1080, 60/1) and look at a contact sheet of frames across the whole take. A take can finish "successfully" and still be wrong, for example when the input had no focus and nothing was typed.
- If the footage shows something the user's real setup doesn't (gaps, padding, cut-off panels), re-record with the setup fixed. Don't crop or scale it away.

## Remotion

Do all editing in Remotion. Don't build a custom renderer.

- Keep the whole edit as data in one file, such as `src/edit.ts`. Each shot has its cue, an enter transition, captions timed relative to the shot, and clips as `{ file, in, out }` in raw-file seconds.
- Symlink the footage folder into `public/`, and generate an index of every file's size, fps and duration with `ffprobe`. Warn when a file isn't 60 fps.
- Play clips with `<OffthreadVideo trimBefore={Math.round(in * fps)} muted />` so each source frame maps to one output frame.
- Add a check that throws when a shot's cut misses its beat cue, naming the shot to lengthen. When you swap clips, keep the shot's total length unless you also move the cue.

## Cutting

- Find cut points by looking at frames around each moment, not from logs or guesses. Apps react later than the keystroke.
- Cut into a moment just before it happens, and out just after the result is readable.
- Several short clips from one take beat one long clip: show the action, then jump to the result.
- Check edits with stills of the changed frames (`npx remotion still <comp> out.png --frame=N --scale=0.5`) before rendering the whole thing.
