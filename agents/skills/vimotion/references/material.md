# Material

Material is everything the video is built from: recordings, screenshots, logos, images, fonts and audio. It is real, and it stays exactly as captured or published.

## Never edit raw material

- Don't trim, crop, scale, speed up, color-correct, blur, join or retouch raw files. Every change happens in Remotion, which reads the raw files and leaves them untouched.
- Recordings are one continuous take per scene, starting a couple of seconds before the first action and ending a couple of seconds after the last. Re-encoding to constant-frame-rate H.264 is the only processing allowed.
- If a take is wrong, record a new one (`-take2`, `-take3`) and keep the old one. Never overwrite or patch a take.

## Real and prepared

- Show the product with real data in a real context, such as its own repository or a realistic real project. No invented sample data.
- Prepare before recording: logos for every place the product shows one, real files, real images. Keep private data off screen at the source: a clean environment, a clean project, test accounts.
- Images, fonts and audio need a licence that allows this use. Record each source and licence in `CREDITS.md`.

## Capture

The source decides the quality: nothing in the edit can fix blurry, choppy or compressed footage.

- Capture at the final quality or higher. Don't assume a tool's recording meets the target; measure it, and switch methods if it falls short.
- Keep the machine quiet while capturing, so nothing competes with the recording.
- When footage looks wrong, find the cause in how it was captured and record again. Don't hide it in the edit.

## Tools

A tool you know can still behave differently in another environment. Don't script from memory or guess: read its help, try each step while watching the result, then script what worked. Judge by the actual result, not by a command or script finishing without errors.
