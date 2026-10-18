# Scheduling

Decides which vocabs to practice and when. Used by [Quiz](quiz.md) and [Flash Card](flash_card.md) to pick vocabs, and by [View Vocabs](view_vocabs.md) to show progress.

Practice is online only (AI features need a connection anyway), so the server creates every session and grades every answer.

## Source of truth

Answered questions (see `Question` in the ERD) are the only stored facts. Progress is derived from them on every request, never stored.

- A question counts once it has `answered_at`. Unanswered questions (e.g. an abandoned session) are ignored.
- Quiz and flash card answers count the same (self-graded flash cards included)
- `answered_at` is server time

## Ladder

```
LADDER = [1, 3, 7, 14, 30, 60, 120]   # days

last_answered_at = latest answered_at of the vocab
streak           = correct answers after the latest wrong one (all of them if never wrong)
interval         = streak == 0 ? 1 : LADDER[min(streak, 7) - 1]
due_at           = last_answered_at + interval days
state            = no answers → new;  streak >= 5 → known;  otherwise learning
```

| Answer | Result | Streak | Due in |
| --- | --- | --- | --- |
| 1 | correct | 1 | 1 day |
| 2 | correct | 2 | 3 days |
| 3 | correct | 3 | 7 days |
| 4 | wrong | 0 | 1 day |
| 5 | correct | 1 | 1 day |

Changing the rule (e.g. to FSRS) is a code change only. The answer history is replayed through the new rule, no migration.

## Query

Streak and last answer for a set of vocabs, in one query (the ladder is applied in Go):

```sql
WITH answered AS (
  SELECT vocab_id, correct, answered_at
  FROM question
  WHERE answered_at IS NOT NULL AND vocab_id = ANY($1)
),
last_wrong AS (
  SELECT vocab_id, max(answered_at) AS t
  FROM answered
  WHERE NOT correct
  GROUP BY vocab_id
)
SELECT a.vocab_id,
       max(a.answered_at) AS last_answered_at,
       count(*) FILTER (WHERE a.correct AND a.answered_at > coalesce(w.t, '-infinity')) AS streak
FROM answered a
LEFT JOIN last_wrong w USING (vocab_id)
GROUP BY a.vocab_id;
```

Index: `question (vocab_id, answered_at) WHERE answered_at IS NOT NULL`.

Starting a session computes this for every vocab in the course (e.g. 2,000 vocabs × 20 answers = 40k rows, milliseconds). When that gets slow, cache `due_at` on the vocab and update it on every answer. It stays rebuildable from the answers.

## Picking vocabs for a session

Given a course and a size:

1. Compute progress for every vocab in the course
2. Due: `due_at <= now`, most overdue first
3. New: no answers, oldest saved first, up to `NEW_PER_DAY - (vocabs whose first answer was in the last 24 hours)`
4. Take due first, then new, up to the size. Fewer available means a shorter session. None means "nothing to practice".

The caller may skip vocabs it can't use (e.g. the quiz skips vocabs with no eligible question type) and take the next one.

Constants:

- `NEW_PER_DAY = 10`. A rolling 24 hours instead of a calendar day, so no timezone is needed.

## Trust model

- Quiz answers are graded by the server. The correct option is never sent before answering, and a question can only be answered once.
- Flash cards are self-graded and trusted
- Cheating only hurts the cheater's own schedule, so this is enough. Revisit before any feature where answers matter to others or carry value (leaderboards, rewards, teacher dashboards): plausibility checks (minimum time per answer, answers per vocab per day) and rate limits.

## Future

- Cache `due_at` on the vocab
- FSRS instead of the ladder
- `NEW_PER_DAY` per course, user configurable
- Calendar day in the user's timezone instead of a rolling 24 hours
