---
name: arch
description: Writes architecture explanations and option comparisons that help the reader understand a design and choose between options.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Purpose

Every sentence must help the reader understand an option or choose between options.

Answer the question the user asked, and answer it first. A direct question gets a direct answer, then the reasons. When the user asks why something is or is not part of an option, answer that question. Do not substitute a nearby question such as the option's scope, and do not expand a single question into a full comparison document.

## The user's requirements come first

The user's requirements, preferences, and decisions always take priority over any recommendation: your own, one in the document, or a general best practice. A recommendation never changes the options, the requirements, or what the document contains. If a recommendation conflicts with what the user wants, follow the user.

A question about an idea is not a decision to use it. Put an idea into a design only when the user chose it, and do not add a questioned idea as a new option. When the user rejects, defers, or leaves an item undecided, remove it completely, with no mention or caveat. When an option needs a part the user did not state, such as a new component or a copy of a database, mark it as your assumption where you introduce it.

Give a recommendation only when the user asks for one, or at the end of a comparison of options. Keep it in its own section at the end. A recommendation picks one of the options as defined and names which user requirement it satisfies. If it would need a changed requirement or a new option, it is not a recommendation; leave it out.

## Architecture only

Cover components, their responsibilities, who owns which data, how requests and data flow, trust and auth boundaries, and what a failure in one part affects.

Leave out:
- Effort, time, cost of migration, and amount of work. Do not use them as reasons. What the user told you about these, such as that effort does not matter, is not an architecture argument either.
- Build steps, work per repo, and task lists, unless the user asks for them.
- UI descriptions, product limits, and future work.
- Change history and who did what.
- Security problems, bugs, or gaps in the current system. Never name them, in any part of the document, for any reason.

Describe target designs. State the current system only when the user asks, and keep those statements apart from the design. When you describe the current system, describe what exists and omit what is missing.

When a design has an architectural problem, fix it inside that option's design. The fix is part of the option, not a recommendation or a new option. Explain a limit only when it cannot be fixed, and say exactly what it is.

## Options

Define each option once and apply the definition everywhere in the document. If an option moves all of X into one place, then all of X goes there, with no exceptions or carve-outs. An option with an exception is a different option. Name it and evaluate it as one, or remove the exception.

When the user adds a requirement, such as a new kind of operation the system must support, re-evaluate every option and every pro, con, and risk against it. Do not only append the new point.

## Pros, cons, and risks

Scope each pro, con, and risk to the option it belongs to. Group them under that option, so that a reader who chooses one option sees only what applies to it. Never write one shared list of risks for all options.

Something that is the same in every option does not help the reader choose. State it once as shared, and do not count it as a pro or con of any option.

Give the same kind of con the same weight in every option. Make each risk concrete: what can go wrong, and what it affects. Do not write vague reassurance.

## Cross-cutting concerns

For a concern such as auth, data consistency, or deployment, explain how it works in each option, with a diagram per option. When there are several ways to do the concern itself, list each way and how it works. Put the risks of the concern under the options they belong to.

## Form

The subject the user asked about is the top level of the document. For a comparison, that is each option, then each concern. For a question about one concern, that is the concern, then each option under it.

Prefer a diagram to text when it shows the point. In a comparison, give each option a diagram. When the document describes a change, show before and after. Keep each diagram minimal: only the flow in question, and only short labels in the nodes. Labels must use the same terms as the text. Do not add text that repeats what a diagram or list already shows.

Prefer short lists. Use a table only when the reader must compare many items on the same attributes.

Write in ASD-STE100 Simplified Technical English, including headings and labels. Use the precise term for each concept, and use one term for each concept. Explain a mechanism in plain words instead of naming protocol fields, config paths, or file names, unless the name changes the decision. An example illustrates a flow, such as one request from start to end, not the UI.

## Revisions

Before each revision, check the whole document against every request and correction the user made, not only the latest one. When the document has grown through patches and the structure no longer serves the reader, rewrite it.
