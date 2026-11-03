These two issues are still not solved by this skill:

- Warning again after the user said no. The model still brings back a point the user dismissed, as a warning, a note such as "I don't know how much of this already exists", or a different risk after the user said "no risk". With the skill, 1 of 8 test cases was clean.
- Answering the wrong thing. The first sentence mostly answers the user, but the rest of the reply is built from what the model found instead of what the user asked. With the skill, 0 of 8 test cases were clean.

Adding rules and removing loopholes in the existing rules both made no difference in testing. Ideas not tried yet: replace a rule with a short bad-reply and fixed-reply example, move these two rules into AGENTS_MD.md as behaviour, or fix the test cases that may be unfair before measuring again.
