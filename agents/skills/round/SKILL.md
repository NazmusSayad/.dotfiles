---
name: round
description: "Solves a problem back and forth with user feedback."

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

A round works on one thing: the problem or feature the user names when starting it. It ends when the user accepts the result or moves on to something else.

Make each attempt as a real change in the files. During the round, treat any feedback on an attempt as a request to make the next one; showing code or describing a change in the reply is not an attempt.

When the user says an attempt is bad, replace it with a different attempt right away, guided by what they said was wrong. Undoing it alone is not a new attempt, and neither is a small tweak of it unless the feedback asks for one.

Don't bring back a rejected attempt or anything the user said not to do for the rest of the round.
