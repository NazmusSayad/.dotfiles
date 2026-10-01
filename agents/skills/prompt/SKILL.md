---
name: prompt
description: Creates and revises agent-facing prompts, including skills, system instructions, agent files, task prompts, and reusable rules.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Audience

The audience for a prompt is an AI agent, not the user who requested it. An agent is nondeterministic: a prompt shapes how it interprets context and makes decisions rather than mechanically determining one result.

A task prompt, skill, system instruction, agent file, or reusable rule is an agent-facing prompt. They share that purpose but govern different situations and timescales.

## Define the document's job

Identify the behavior the document should influence, when it applies, and where it belongs. A task prompt directs one task. A skill teaches subject-specific judgment across a class of tasks. Agent and system instructions establish broader behavior within their scope.

Put an instruction at the narrowest scope where it remains true. Do not repeat guidance already supplied by a broader source. For reusable prompts, names, descriptions, paths, and activation metadata are part of the design when they determine whether the instructions are discovered or loaded.

A skill should teach its subject, not explain its invocation or read like documentation for a human newcomer. Generic agent workflow does not belong unless it is part of the skill's actual subject.

## Preserve the intended behavior

Clarify and organize the user's intent without replacing it with the prompt writer's preferences. Preserve the outcome, scope, facts, actors, qualifiers, uncertainty, negation, tone, force, and whether the user wants discussion, evaluation, creation, or execution. Semantic fidelity matters more than literal wording.

Separate outcomes, constraints, and methods. Do not prescribe a method merely because it sounds prudent or resembles prompt-engineering practice. Include it when the method itself matters, the environment requires it, or evidence shows that it prevents a relevant failure.

Treat references, examples, and surrounding conversation as evidence of intent, not requirements to copy wholesale. Distinguish instructions from quoted material, facts, examples, and other content when the agent could plausibly confuse them.

## Guide judgment

State the decisions that matter and leave harmless choices to the agent. Broad guidance should provide useful criteria, distinctions, priorities, and boundaries without trying to define every future case. If the prompt solves every decision in advance, the agent has no judgment left to apply.

Use force deliberately. Reserve absolute language for hard requirements and express defaults or preferences as such. When important concerns can conflict, state which one takes priority instead of relying on repetition or emphasis.

Every instruction needs a basis in the intended outcome, actual environment, domain, or evidence from use. A brief reason can help the agent apply a broad rule to situations the writer cannot predict. Do not add explanations for concepts the agent already understands unless the explanation changes how it should decide.

Prefer a positive selection principle over a catalog of prohibitions. Name a failure mode when it is plausible, consequential, and not already prevented by clearer guidance. Add personas, procedures, tool rules, schemas, examples, or review steps only when they materially improve the intended behavior.

## Let form follow purpose

There is no universal prompt template. Use sections, lists, examples, schemas, or required output shapes when they help the agent distinguish instructions or produce a usable result. Do not force every document into role, context, task, constraints, and output sections.

Examples should clarify a real boundary, not decorate the prompt or define the entire valid range by accident. Make clear what an example demonstrates when its incidental details could be mistaken for rules.

Keep the document concise by removing onboarding, repetition, generic rituals, and details that do not affect behavior. Do not remove subject-specific distinctions merely to make it shorter.

## Create, revise, and evaluate

Creating a prompt requires recovering the intended behavior and choosing the right scope. Revising one requires identifying what already works and what behavior is failing. Make the smallest effective change when the existing framing is sound; replace the framing when it is the source of the problem.

Evaluate the prompt according to its job. Use real or representative inputs and outputs for narrow task prompts when practical. Review broad skills and durable instructions through varied situations, competing interpretations, boundaries, and conflicts. A finite set of cases cannot prove how a nondeterministic agent will behave in every future context.

Do not turn one failure or successful fix into a universal rule without evidence of a broader pattern. Every instruction should have a clear behavioral purpose grounded in intent, context, scope, or evidence. Remove anything that cannot meet that test.
