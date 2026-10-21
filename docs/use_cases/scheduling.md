# Scheduling

Decides which vocabs to practice. Used by [Quiz](quiz.md) and [Flash Card](flash_card.md) to pick vocabs.

Practice is online only (AI features need a connection anyway), so the server creates every session and grades every answer.

## Score

Each vocab stores `xp` and `correct_streak` (both start at 0) and `last_reviewed_at` (null until the first answer, which means the vocab is new). All three are updated on every answer:

- Correct: `xp += 1`, `correct_streak += 1`
- Wrong: `xp = max(0, xp - WRONG_PENALTY)`, `correct_streak = 0`
- `last_reviewed_at = now` (server time)

Quiz answers and flash card reviews count the same (self-graded flash cards included). Unanswered questions (e.g. an abandoned session) change nothing.

## Priority

```
priority = days_since(coalesce(last_reviewed_at, created_at)) / (1 + xp)
```

A new vocab counts its age from `created_at`. Highest priority first. The more xp a vocab has, the less each day of age counts, so it takes longer to come back.

- A new vocab (xp 0) has priority equal to its age in days
- A vocab with xp 4 needs 5 times as long to reach the same priority
- A vocab just answered has priority near 0, so it doesn't repeat right away
- A wrong answer lowers xp, so the priority jumps up and it comes back soon
- Priority keeps growing with time, so even a strong vocab eventually comes back

## Picking vocabs for a session

Given a course and a size, take the vocabs with the highest priority:

```sql
SELECT ... FROM vocabs
WHERE course_id = $1
ORDER BY extract(epoch FROM now() - coalesce(last_reviewed_at, created_at)) / 86400 / (1 + xp) DESC
LIMIT $2
```

There is no due cutoff, so a session is always full while the course has enough vocabs. Fewer vocabs than the size means a shorter session. No vocabs means "nothing to practice".

The caller may skip vocabs it can't use (e.g. the quiz skips vocabs with no eligible question type). It asks for extra (e.g. `size * 2`) and takes the first that fit.

Constants:

- `WRONG_PENALTY = 2`. Larger than the +1 for a correct answer, so a hard vocab doesn't drift away.

## State

Shown in [View Vocabs](view_vocabs.md), derived from the stored fields:

```
state = last_reviewed_at is null → new;  xp >= KNOWN_XP → known;  otherwise learning
```

- `KNOWN_XP = 5`. Uses xp rather than the streak so the state agrees with the priority, and one miss doesn't demote a well-known vocab. `last_reviewed_at` is still needed for `new`, since xp can drop back to 0.
- streak: `correct_streak`, shown only as a number

## Trust model

- Quiz answers are graded by the server. The correct option is never sent before answering, and a question can only be answered once.
- Flash cards are self-graded and trusted
- Cheating only hurts the cheater's own schedule, so this is enough. Revisit before any feature where answers matter to others or carry value (leaderboards, rewards, teacher dashboards): plausibility checks (minimum time per answer, answers per vocab per day) and rate limits.

## Future

- A due cutoff, so "nothing to practice" can happen (e.g. `last_reviewed_at < now - 1 day`)
- Faster growth than linear (e.g. divide by `2^min(xp, 7)`)
- FSRS instead of xp
