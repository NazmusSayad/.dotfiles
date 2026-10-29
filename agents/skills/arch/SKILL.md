---
name: arch
description: Writes architecture explanations and option comparisons that help the reader understand a design and choose between options.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Purpose

Every sentence must help the reader understand the design or choose between options. Mark every assumption where you introduce it.

Answer the question the user asked, and answer it first. A direct question gets a direct answer, then the reasons. If the question's premise is wrong, say so first. When the user asks why something is or is not part of an option, answer that question. Do not substitute a nearby question such as the option's scope. Add nothing the question does not need, including sections that describe an option beyond the question.

When a question arrives together with a revision request, put the answer in its own section before the document.

## The user's requirements come first

The user's requirements, preferences, and decisions always take priority over any recommendation: your own, one in the document, or a general best practice. A recommendation never changes the options, the requirements, or what the document contains. If a recommendation conflicts with what the user wants, follow the user.

A question about an idea is not a decision to use it. Put an idea the user asked about into a design only when the user chose it, and do not add a questioned idea as a new option. When the user rejects, defers, or leaves an item undecided, remove it completely, with no mention or caveat.

Give a recommendation only when the user asks for one, or at the end of a comparison of options. Keep it in its own section at the end. A recommendation picks one of the options as defined and names which user requirement it satisfies.

## Architecture only

Cover where state lives and who may change it, how data and control move, where trust changes, and what a failure affects. In a comparison, cover what differs between options. State a fact that is the same in every option once, as shared, and only when an option cannot be understood without it.

Leave out:
- Effort, time, cost of migration, and amount of work. Do not use them as reasons, even when the user states them.
- Build steps and task lists, unless the user asks for them.
- UI descriptions, unless the option changes what the client stores, computes, or sends.
- Product limits and future work.
- Change history and who did what.
- Security problems, bugs, or gaps in the current system. Never name them, in any part of the document, for any reason. A risk is a property of an option's target design, never a defect of what exists today.

Describe target designs. Show the current system only when the user asks, or when the document describes a change to it; then once, as a before state, apart from the options. When you describe the current system, describe what exists and omit what is missing.

## Options

Define each option once and apply the definition everywhere in the document. An option's rule has no exceptions. If the user wants an exception, it is a separate option with its own name, pros, cons, and risks.

When an option's definition leaves a choice open, such as where the data lives or where a check runs, or when its design has an architectural problem, fill or fix it inside that option with the smallest addition its definition implies, using only the components the option names. The addition is part of the option, not a recommendation or a new option. If the problem cannot be fixed inside the definition, state it as a con or risk of that option, and say exactly what it is.

When the user adds, changes, or removes a requirement, re-evaluate every option and every pro, con, and risk against it. Do not only append the new point.

## Pros, cons, and risks

Scope each pro, con, and risk to the option it belongs to. Group them under that option, so that a reader who chooses one option sees only what applies to it. Never write one shared list of risks for all options. A risk that exists in every option goes under each option.

A fact that is the same in every option is not a pro or con of any option; a risk that is the same in every option still goes under each option.

Describe the same kind of con with the same detail in every option. Make each risk concrete: what can go wrong, and what it affects.

## Cross-cutting concerns

For a concern such as auth, data consistency, or deployment, explain how it works in each option. Describe the way the option uses. If the definition does not fix the way, choose the one it implies, mark it as an assumption, and describe only that one. List the other ways, each with how it works, only when the user asks what the ways are. Put the risks of the concern under the options they belong to.

## Form

The direct answer comes first. After it, the subject the user asked about is the top level of the document. For a comparison, that is each option, then each concern with each option under it. For a question about one concern, that is the concern, then each option under it.

Prefer a diagram to text when it shows the point. In a comparison, give each option a diagram when it has a structure or flow to show. For a concern, give a diagram per option when its flow differs between options. Keep each diagram minimal: only the flow in question, and only short labels in the nodes. Labels must use the same terms as the text. State each fact once in the document, and do not restate a diagram or list in text.

Prefer short lists. Use a table only when the reader must compare many items on the same attributes.

Write in ASD-STE100 Simplified Technical English, including headings and labels. Use the precise term for each concept, and use one term for each concept. Explain a mechanism in plain words instead of naming protocol fields, config paths, or file names, unless the name changes the decision. An example illustrates a flow, such as one request or one write from start to end, not the UI.

## Revisions

Before each revision, check the whole document against every request and correction the user made, not only the latest one. When the document has grown through patches and the structure no longer serves the reader, rewrite it.
