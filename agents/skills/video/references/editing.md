# Editing

## Footage in

Before using a take, measure it and look at frames across the whole take. A take can run without errors and still show the wrong thing.

## Remotion

Do all editing in Remotion. Don't build a custom renderer.

- Keep the whole edit as data in one file: each shot has its beat cue, a transition, captions timed relative to the shot, and clips as `{ file, in, out }` in raw-file seconds.
- Play clips with `<OffthreadVideo trimBefore={Math.round(in * fps)} muted />` so each source frame maps to one output frame.
- Add a check that fails when a shot's cut misses its beat cue. When you swap clips, keep the shot's length unless you also move the cue.

## Cutting

- Find cut points by looking at frames around each moment. Apps react later than the input that triggers them.
- Cut in just before a moment and out just after its result is readable.
- Several short clips from one take beat one long clip: show the action, then jump to the result.
