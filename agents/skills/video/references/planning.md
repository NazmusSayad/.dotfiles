# Planning

## Audience

Decide who is watching before anything else. Every later choice (story, features, claims, pace, captions, audio) serves them.

- Name the audience concretely: for example, developers who live in the terminal and run coding agents, not "users".
- Write down what they already know, what they care about, and what problem they feel today.
- Decide what they should think or do after watching (try `npx sshtty`, star the repo), and build toward it.
- Decide where they will watch (README, social feed, YouTube, a landing page). That sets the length and whether it must work without sound.
- Speak their language: show the tools and workflows they use, and use their words in captions. Skip anything they don't care about, even if it's impressive.
- If the audience isn't clear from the request or the product, ask the user. It's the one question worth asking up front.

## Story

- Lead with what the product is in one line, then show why it matters, then how to get it.
- Pick a handful of features that look good on screen. A feature that does something visual (an image appearing, a shell surviving a restart) beats one that only prints text.
- Find the one bold, true claim and give it an attention hit: a statement, then a short beat that underlines it ("Better images than Windows Terminal." → "Yes, on Windows.").
- Let length follow content. A rough target (say 60 s) is fine, but don't squeeze a scene to hit it.

## Claims

- Check every claim against the code or a real run before it goes in the storyboard.
- Footage from one OS can stand in for another only if the product behaves the same there. Never present it as a different product.

## Scenes

- One idea per scene.
- Keep scenes simple. For an AI agent: open it, give one short prompt, cut after the first few messages. Don't wait for the full answer.
- Plan what real content each scene needs before recording: logos, files, images, a clean project, a matching config.

## Storyboard

Write `STORYBOARD.md` in a promo folder outside the product repo:

- the audience, what they should do after watching, and where they'll watch
- the one-liner and the rules
- a scene table: rough time, scene, what is shown, the real source, and the raw file name
- where raw footage and per-take notes go

When several agents work in parallel, this file is their shared brief.

## Takes

For each scene, plan one raw take: what happens, in what order, and what the screen must look like before it starts. Each take should note the raw-file time of each real moment (when something appears, not when the key was pressed) so the edit can cut on them.
