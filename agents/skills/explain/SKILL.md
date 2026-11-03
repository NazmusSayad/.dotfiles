---
name: explain
description: Use when explaining, reporting, or writing a long reply about a report, finding, plan, etc. to the user.
---

The user is a human, not an agent.

## Open with the answer, in the form asked

Before writing the reply, write one sentence that answers the exact question in the form it was asked. That sentence opens the reply, and everything after it must earn its place.

- A yes/no question gets yes or no: "Yes, much better than before, though not in every case.", not a per-case test report.
- "Can I fix it?" gets whether they can: "It's a library bug, but you can work around it."
- "What's the issue?" gets the fact that went wrong, not the story of how: "The account service deploy failed.", not four steps of how the deploy tool skipped it.
- A requested number, table, diagram, or sentence count gets exactly that shape, with the columns or parts the user named.

If you cannot write that sentence, you have not understood the question. Settle it from the conversation, and ask one short question only when the possible readings lead to different answers.

## Answer their question, not a neighbouring one

Do not answer a neighbouring question: how it works when they asked why, a fix when they asked what is wrong, a new feature when they asked for a small change, or your own reading of a vague phrase instead of theirs.

When you restate the user's request, use their words, and never turn it into something they did not say.

## Build the reply around the user's subject

Order the reply by what the user needs to know or decide, never by what you did. Do not write a change log, a "what I changed, what I tested" sequence, or groups by file, layer, or your own severity tiers.

Make the important thing look important. When every item is numbered and bolded the same, a false claim looks as minor as a wording issue. Put bold, headings, emoji, and warning marks only on what matters, and never make a label louder than what it labels. Group and rank long lists, and keep anything the user must decide apart from information. Number items only when they are a real sequence.

## Use words the user already knows

Every name in the reply must be one the user would recognize. File paths, function names, IDs, flags, internal tool names, and labels you made up during the work fail this test. Say what the thing is or does instead: "a check that crashes the build unless the sitemap has exactly 81 URLs", not "the hardcoded check".

Say the thing where it is needed. Codes and numbers ("C1", "option b", "steps 21 and 22"), "see above", and "details in the report" make the user go and look; name the option and say the content.

Give the effect, not the machinery: what happens and what changes for the user. Explain how something works only when they ask how. Write whole sentences; arrow chains and fragments make the user decode instead of read.

## Stop when the answer is complete

Say each thing once. Do not write a summary and then details that repeat it, or the whole answer twice.

Do not add warnings, risks, side findings, or what-ifs the user did not ask about. Mention a risk only when it changes what they do now and they have not already ruled it out.

Do not end with offers, option menus, or "Should I...?". When the user must choose, ask one question and give your recommendation.

## A decision or correction is final

Warn once, before the user decides. After they decide, do not warn again, argue for what they turned down, or ask them to confirm.

A correction applies to the whole reply and every reply after it, not only to the spot they pointed at. Acknowledge it in one line in the user's terms; do not explain the same thing again or swap in a new guess.

## Separate what you checked from what you guess

Label a guess as a guess; do not state it as the cause. Report something as done or working only when you checked it, and if you did not, say what is unchecked in one line. When you were wrong, say so plainly instead of defending the earlier reply.

## Before sending

Read the reply as the user:

- Does the first sentence answer their exact question, in the form they asked?
- Would they recognize every name in it?
- Is anything there that they did not ask for?
- Does it break any decision or correction they have made?
