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

## Recording setup

- Record at the final resolution and frame rate or higher. On a 2x display, a 960x540 window captures as a sharp 1920x1080.
- Use the user's real product settings, and act at a human pace with pauses long enough to read each result.

## Catches

The source decides the quality: nothing in the edit fixes blurry, choppy or compressed footage.

- Built-in browser recorders didn't deliver a true 60 fps. Capture the real, visible window natively instead.
- Heavy work during a capture (renders, encodes) drops frames.
- Measure every take, and look at it.
- If the footage doesn't match the user's real setup, the setup is wrong. Fix it and re-record; don't crop.

## Tools

A tool you know can still behave differently in another environment. Don't script from memory: read its help, try each step while watching the result, then script what worked. Judge by what happened on screen, not by the script finishing. With Playwright, for example, a click doesn't always focus the input, so a take can "succeed" with nothing typed.
