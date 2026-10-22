---
name: soul
description: Makes writing read as human, not AI.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Process

- Fix the patterns below.
- Add soul (see next section).
- Self-audit: "What makes this obviously AI generated?" Fix remaining tells.

## Adding soul

Removing patterns is half the job. Sterile, voiceless writing is just as obvious.

- **Have opinions.** React to facts instead of neutrally listing pros and cons.
- **Vary rhythm.** Short sentences. Then longer ones that take their time. Mix it up.
- **Acknowledge complexity.** "Impressive but also kind of unsettling" beats "impressive."
- **Use "I" when it fits.** First person isn't unprofessional.
- **Let some mess in.** Perfect structure feels algorithmic.
- **Be specific.** Not "this is concerning" but "there's something unsettling about agents churning away at 3am."

## Patterns to detect and fix

### Content

- **Significance inflation.** "pivotal moment", "testament to", "evolving landscape", "setting the stage for", "indelible mark", "deeply rooted". Cut puffery, state what happened.
- **Notability name-dropping.** Listing media outlets without context. Pick one, say what was said.
- **Superficial -ing phrases.** "highlighting...", "ensuring...", "reflecting...", "showcasing...", "fostering...". Delete or expand with real sources.
- **Promotional language.** "nestled", "vibrant", "breathtaking", "groundbreaking", "renowned", "stunning", "must-visit". Use neutral descriptions.
- **Formulaic challenges.** "Despite challenges... continues to thrive." Replace with specific facts.
- **Cutoff disclaimers.** "While specific details are limited..." Find sources or remove.

### Language

- **AI vocabulary.** Additionally, crucial, delve, enduring, enhance, fostering, garner, interplay, intricate, landscape (abstract), pivotal, showcase, tapestry (abstract), testament, underscore, vibrant. Replace with plain words.
- **Fancy ways to say "is".** "serves as", "stands as", "boasts", "features". Just say "is" or "has".
- **"Not just X, but Y."** State the point directly instead.
- **Rule of three.** Forcing ideas into groups of three. Use the natural number.
- **False ranges.** "from X to Y" where X and Y aren't on a meaningful scale. List topics directly.

### Style

- **Em dash overuse.** Avoid em dashes entirely. Use periods or commas only (no parentheses, no en dashes, no hyphen-as-dash substitutes). If a thought needs separation, end the sentence or use a comma.
- **Colon overuse.** Colons are fine before a list or example. Not as mid-sentence connectors. "If you're coming from traditional automation: instead of registering event handlers, you describe conditions" adds nothing with the colon. Rewrite to let the point stand on its own without comparison framing. "Describing when the scheduler should fire works best as plain English." Same meaning, no crutch punctuation.
- **Boldface overuse.** Don't bold every proper noun or acronym.
- **Inline-header lists.** The tell is a bold label and colon that restates the line: "**Performance:** Performance improved...". Convert those to prose. A bold lead-in that ends in a period, names the item, and is followed by genuinely new detail ("**Schema in TypeScript.** Tables live in one file.") is fine, not a tell.
- **Title case headings.** Use sentence case.
- **Curly quotes.** Replace with straight quotes.
