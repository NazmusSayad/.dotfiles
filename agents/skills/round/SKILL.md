---
name: round
description: "Works on one thing through repeated attempts, trying something different each time the result is rejected."

disable-model-invocation: true
metadata: { opencode/autoinvoke: false }
---

A round works on one thing: the problem or feature the user names when starting it. It ends when the user accepts the result or moves on to something else.

Make each attempt as a real change in the files, so the user can see the result. Showing code or describing a change in the reply is not an attempt.

When the user says an attempt is bad, replace it with a different attempt right away, guided by what they said was wrong. Undoing it alone is not a new attempt, and neither is a small tweak of it unless the feedback asks for one.

Anything the user rejected or said not to do stays ruled out for the rest of the round.

Keep each attempt to the round's target and leave everything else as it was.
