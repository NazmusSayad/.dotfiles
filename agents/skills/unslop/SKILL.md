---
name: unslop
description: Writing guidelines for removing AI patterns and adding human voice. Use when writing, editing, or reviewing substantive content such as summaries and reports, or when the user asks to improve, simplify, or make text more concise.
---

## Audience and purpose

Identify what the audience needs to understand or accomplish, choose a form and structure that supports that goal, and write at their level. Include only the context, terminology, technical depth, and examples they need. Include implementation details only when they help achieve that goal. Ask the user only when the audience or purpose cannot be determined.

Lead with the reader's needed outcome or information, then explain why when useful. Make the path to that outcome clear. Use direct language, clear headings, and the smallest useful accurate examples.

## Process

- Scan for the patterns below.
- Rewrite. Preserve meaning, match intended tone.
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
- **Vague attributions.** "Experts believe", "Industry reports suggest", "Some critics argue". Name the source or delete.
- **Formulaic challenges.** "Despite challenges... continues to thrive." Replace with specific facts.

### Language

- **AI vocabulary.** Additionally, crucial, delve, enduring, enhance, fostering, garner, interplay, intricate, landscape (abstract), pivotal, showcase, tapestry (abstract), testament, underscore, vibrant. Replace with plain words.
- **Fancy ways to say "is".** "serves as", "stands as", "boasts", "features". Just say "is" or "has".
- **"Not just X, but Y."** State the point directly instead.
- **Rule of three.** Forcing ideas into groups of three. Use the natural number.
- **Synonym cycling.** Protagonist, main character, central figure, hero all in one paragraph. Pick one, repeat it.
- **False ranges.** "from X to Y" where X and Y aren't on a meaningful scale. List topics directly.

### Style

- **Em dash overuse.** Avoid em dashes entirely. Use periods or commas only (no parentheses, no en dashes, no hyphen-as-dash substitutes). If a thought needs separation, end the sentence or use a comma.
- **Colon overuse.** Colons are fine before a list or example. Not as mid-sentence connectors. "If you're coming from traditional automation: instead of registering event handlers, you describe conditions" adds nothing with the colon. Rewrite to let the point stand on its own without comparison framing. "Describing when the scheduler should fire works best as plain English." Same meaning, no crutch punctuation.
- **Boldface overuse.** Don't bold every proper noun or acronym.
- **Inline-header lists.** The tell is a bold label and colon that restates the line: "**Performance:** Performance improved...". Convert those to prose. A bold lead-in that ends in a period, names the item, and is followed by genuinely new detail ("**Schema in TypeScript.** Tables live in one file.") is fine, not a tell.
- **Unnecessary numbering.** Do not number headings, sections, or list items by default, whether in documents or replies to the user. Use numbering only when it serves a clear purpose for the reader.
- **Title case headings.** Use sentence case.
- **Decorative emojis.** Remove from headings and bullets.
- **Curly quotes.** Replace with straight quotes.

### Communication artifacts

- **Chatbot phrases.** "I hope this helps!", "Let me know if...", "Of course!", "Certainly!", "Found the smoking gun!" Remove.
- **Cutoff disclaimers.** "While specific details are limited..." Find sources or remove.
- **Sycophantic tone.** "Great question! You're absolutely right!" Respond directly.

### Filler

- **Filler phrases.** "In order to" becomes "To". "Due to the fact that" becomes "Because". "It is important to note that" gets deleted.
- **Excessive hedging.** "could potentially possibly be argued that it might" becomes "may".
- **Generic conclusions.** "The future looks bright." State specific plans or facts.

### Jargon

- **Abstract metaphor nouns.** Substrate, wedge, vector, locus, vantage, nexus, primitive (as noun), harness (as metaphor), surface (as in "API surface"), bedrock, scaffolding (as metaphor), modality, paradigm, gold-plating, ratchet (as metaphor), evacuate (for moving code), endgame, north star, flywheel. These read as technical but usually have a plainer concrete word. "Substrate" becomes "base". "Wedge in" becomes "add". "Vector" becomes "way" or "method". "Gold-plating" becomes "more than the job needs". "Ratchet" becomes the mechanism's real name or "a limit that only tightens". "Evacuate" becomes "move out". "Endgame" becomes "the last phase". Pick the concrete word.

### Plain speech

- **Say what it does, not how it feels.** "the database stays close at hand", "SQL you can read", "types that follow your schema" name a feeling. The fix names the mechanism or a number: "`.toSQL()` returns the exact string sent to the database", "a column rename fails the build". Ask what the sentence tells the reader to do or know, then write that. If you can't restate it as a concrete instruction, fact, or number, cut it. One more check: if the sentence could appear unchanged in another project's docs, it says nothing about this one. Cut it.
- **Shorten or split dense sentences.** If the reader has to backtrack to parse a sentence, break it in two or drop clauses. One idea per sentence.
- **Active voice.** Prefer it. Catch "is/are/was/were + past participle" and name the actor: "queries are validated" becomes "the compiler validates queries", "the file is parsed by the loader" becomes "the loader parses the file". Passive is fine only when the actor is unknown or genuinely doesn't matter.
- **Cut adverbs, or use a stronger verb.** "runs quickly" becomes "is fast" or the number. "significantly improves" becomes the measured delta. An adverb propping up a weak verb means the verb is wrong.
- **Prefer the plain word.** "utilize" becomes "use", "leverage" becomes "use", "facilitate" becomes "help", "numerous" becomes "many", "in the event that" becomes "if". The fancier synonym is rarely clearer.
- **Mannered prose.** Metaphor or flourish where a literal phrase exists: aphorisms ("wire it or delete it"), rhetorical fragments for effect, personified code ("the plan holds it"), figurative verbs ("rides along", "stands on"), stock framing phrases. "A dial worth turning" becomes "a parameter worth varying". Say what you mean. The "Abstract metaphor nouns" rule covers the metaphor nouns.
- **Over-compression.** Dropped articles, verbless fragments, symbol-speak, and abbreviations that make the reader decode instead of read. "Parser rejects bad date → exit 2, no write" becomes "The parser rejects a bad date, exits with code 2, and writes nothing." Write whole sentences with their articles and verbs, and spell out arrows and abbreviations.
