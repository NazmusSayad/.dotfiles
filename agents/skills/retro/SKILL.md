---
name: retro
description: Audits the user's history and reports issues.

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

## Purpose

The audit is a report that helps the user understand what keeps going wrong between them and their agents. It is not a fix. Do not edit instructions, skills, settings, hooks, or project files, and do not offer to apply changes; the only file you write is the report. A finding may say what would help, but its job is to explain the problem and its cause clearly enough that the user can decide what to do.

The most important findings are repeated ones: mistakes the agent made again and again, and problems the user's prompting caused again and again. One-off incidents matter only when they were costly.

## Reading the history

Session transcripts are JSONL files in `~/.claude/projects/<project>/<session>.jsonl`. Subagent transcripts live in `<session>/subagents/` and can be skipped unless a finding depends on what a subagent did. `~/.claude/history.jsonl` lists every prompt the user typed with its `sessionId`; prompts whose session has no transcript left show what the user asked but not how the agent responded.

Most `user` entries are not the user's prompts. Only entries with `origin.kind == "human"` are typed prompts. Entries with `isMeta` are skill expansions or system context, `tool_result` blocks are tool output, and `[Request interrupted by user...]` marks an interruption. Tool results with `is_error` include tool calls the user rejected, often with their reason. `AskUserQuestion` results contain the user's answers. `isCompactSummary` marks a context compaction. Skill loads appear as `Skill` tool calls and as `isMeta` entries starting with `Base directory for this skill:`.

If the user names a date range or project, analyze only the sessions in that scope. Otherwise analyze the whole history.

The history is far too large to read in one context, and most of its size is tool output. Condense each session to the user's prompts, interruptions, rejections, answers, and the agent's text and tool calls, keeping tool output only as its size and any errors. Then split the sessions into batches and analyze them with parallel subagents, giving each its session list and the path to this file. The analysis is complete when every session with a human prompt has been assigned to exactly one batch and every batch has returned findings with evidence, not summaries of what each session was about.

## What to look for

The user's corrections are the strongest evidence of a problem: interruptions, rejected tool calls, repeated requests, "no", "I said", reverted work, and frustration. For each one, read what happened before it and decide the cause:

- **Prompting:** the prompt was ambiguous, underspecified, or missing context the user supplied later.
- **Unwritten preference:** the user had to restate something that was never written down.
- **Weak instruction:** an existing instruction was unclear, conflicted with another, or was not loaded when needed.
- **Agent failure:** the agent failed despite clear instructions, either by producing a wrong result or through one of the behavior problems below.

This attribution is the core of the audit. Do not blame the user for an agent failure or excuse the agent for ignoring a clear instruction.

Behavior problems are about how the agent works with the user rather than whether a result was correct: doing more or less than asked, acting when the user was discussing, going deep on side investigations instead of delivering what was asked, asking questions it could answer itself or deciding things that were the user's to decide, continuing after an interruption without adjusting, looping or giving up, agreeing instead of pushing back, and claiming completion without verifying. Keep these apart from agent mistakes, which are wrong results such as broken code, wrong files, or incorrect facts.

Some problems never produce a correction, so also look for:

- **Navigation:** the agent spent many tool calls finding a file or fact that a pointer in an instruction file would have made obvious.
- **Missing checks:** the agent made a mistake that a linter, type check, test, or CI job could have caught. Check whether the repo already has one that went unused, and note repos with no guardrail at all.
- **Information access:** the agent guessed or asked the user for information it could have read itself with access to logs, a service, or a CLI.
- **Repeated workflows:** requests, pasted context, or follow-ups like "check now" that recur across sessions and could become a skill, command, or durable instruction.

Tooling problems also deserve their own attention:

- **Failing tools:** tool calls, CLIs, MCP servers, hooks, or permission rules that error, hang, or block repeatedly, and what the agent did in response.
- **Tool misuse:** the agent using the wrong tool for the job, misusing a tool's arguments, or working around a tool instead of using it.
- **Tool economy:** tool calls with very large output, repeated identical calls, or tooling that wastes context.

Communication problems are failures even when the agent's actions were right: replies the user could not follow, answers that missed the question, conclusions buried in detail, jargon, or excessive length. Signals include "what?", the user rephrasing the same question, and the user acting on a misunderstanding of the reply.

Skill loading problems show up in which skills loaded and when:

- **Missed loads:** an automatically invoked skill whose subject came up but that did not load. Manually invoked skills, marked `disable-model-invocation: true`, never load on their own and do not count.
- **Wrong loads:** a skill that loaded for work outside its subject.
- **Unused skills:** skills that never loaded in the analyzed history.

Typos and casual wording matter only when they caused a misunderstanding.

Cite each finding with the session file, timestamp, and a short quote. Count how often it recurs and across how many sessions and projects. A pattern seen once is an anecdote, and the report should say so. Never quote secrets, credentials, or tokens from transcripts.

## Compare with existing instructions

Read `~/.claude/CLAUDE.md`, the `CLAUDE.md` files of projects in the history, `~/.claude/settings.json`, and the skills in `~/.claude/skills`. For each recurring problem, note whether an instruction already covers it. A covered problem that keeps happening says something different from an uncovered one. Also note instructions that are never relevant, that would not change the agent's default behavior, or that led to behavior the user corrected.

## Report

Write the report to `~/.claude/audits/<YYYY-MM-DD>.md`. Start with the coverage: the sessions, prompts, projects, and date range analyzed, plus anything skipped.

Then include these sections, each ordered by how much the problem costs the user, weighing frequency and impact:

1. **Repeated agent mistakes:** wrong results the agent keeps producing, how often, where, and whether an existing instruction should have prevented them.
2. **Behavior issues:** how the agent keeps working against the user's intent, how often, where, and whether an existing instruction should have prevented it.
3. **Repeated prompting problems:** what in the user's prompts keeps causing misunderstandings or rework, how often, with examples of the prompt and what it led to.
4. **Communication issues:** the kinds of replies that confused the user or missed the point, with examples.
5. **Tooling issues:** failing tools, tool misuse, and tool economy, with how often each occurred and what it cost.
6. **Skill loading problems:** missed loads, wrong loads, and unused skills.
7. **Environment gaps:** navigation, missing checks, information access, and repeated workflows.
8. **Instruction issues:** existing instructions that are weak, conflicting, unused, or counterproductive.
9. **Rework hotspots:** the sessions that cost the most in back-and-forth, retries, and compactions, and what went wrong in each.

If earlier reports exist in `~/.claude/audits/`, end with **Change since the last audit:** for each problem in the most recent earlier report whose scope overlaps this one, say whether it is resolved, still happening, or worse, and mark which problems in this report are new.

In the chat, summarize the few most important repeated problems and link the report.
