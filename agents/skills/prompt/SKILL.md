---
name: prompt
description: Creates and tunes prompts that express the user's intended task without adding the model's own goals, requirements, or process.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Audience

The audience for a prompt is an AI agent, not the user who requested it. An agent is nondeterministic: the prompt guides how it interprets the task and makes decisions.

## Preserve the intended task

A prompt should make the user's task clearer, not replace it with a different task that seems better to the writer. Preserve the intended outcome, scope, facts, actors, qualifiers, uncertainty, negation, requested action, output, tone, and force where they matter. A cleaner sentence is worse if it changes who does what, turns a statement into a task, answers a question that should be preserved, or softens the meaning.

Treat examples, references, and surrounding conversation as evidence of intent, not as requirements to copy wholesale. Do not invent capabilities, source material, audiences, constraints, success criteria, or work the user did not request. A prompt for discussion, evaluation, or explanation must not silently become a request to implement.

## Specify what matters

Pin down choices that determine whether the result is correct. Leave the rest to the model's judgment. Include context when it prevents a plausible misunderstanding, and add a constraint only when it comes from the user's request or the real environment. Do not define every case or solve the task inside the prompt.

State the desired behavior directly. Name a failure mode when it is realistic and important, not to build an exhaustive list of everything the model must avoid. Do not explain familiar concepts or add personas, reasoning rituals, retries, self-review loops, tool instructions, or rigid schemas unless the task genuinely needs them.

## Let structure follow the job

There is no universal prompt template. Do not force every prompt into role, context, task, constraints, and output sections. Keep a simple request short. Use headings, lists, examples, or a required output shape only when they make distinct instructions easier to follow or the result easier to use.

Examples should clarify a real boundary rather than decorate the prompt. Prefer actual representative inputs and outputs when they exist. An invented example can introduce behavior the user never wanted.

## Tune from evidence

When a prompt already exists, preserve what works and change the part responsible for the observed failure. Tune it instead of rewriting it around a preferred prompting style. Make the smallest change that fixes the problem without weakening useful behavior elsewhere.

Test each prompt against its own purpose using real or representative cases. Do not copy a fix into other prompts merely because they look related. One failure does not justify a broad rule unless the evidence shows the same underlying problem.

Every instruction in the finished prompt should trace back to the user's goal, the actual context, or an observed failure. Remove anything added only because it sounds like good prompt engineering.
