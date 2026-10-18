# Flash Card

Self-graded recall of the user's saved vocabs, without a sentence. Works for every vocab, including vocabs only saved from parts. Online only.

Uses the same session and question tables as [Quiz](quiz.md), with `kind = flash_card`.

## Prerequisites

- User is signed in
- User has at least one course

## Screen

```
┌─────────────────────────────┐        ┌─────────────────────────────┐
│                             │        │            読む              │
│            読む              │  tap   │            よむ              │
│                             │  ───►  │          to read            │
│       [ Show answer ]       │        │   [ Missed ]   [ Got it ]   │
└─────────────────────────────┘        └─────────────────────────────┘
```

1. Front: `dict_form`
2. Tap to flip. Back: `reading` and `translation`
3. User grades themselves: Missed → `correct: false`, Got it → `correct: true`

## Start

`POST /flash-card/start`

### Input

- course id

```json
{
  "course_id": "123"
}
```

### Output

- session_id
- questions: the unanswered ones, in order. The back of the card is included, there is nothing to hide in a self-graded exercise.

```json
{
  "session_id": "...",
  "questions": [
    { "id": "...", "dict_form": "読む", "reading": "よむ", "translation": "to read" }
  ]
}
```

No vocab available returns `questions: []` and no session is created.

### Flow

Same as [Quiz](quiz.md) Start, with `kind = flash_card`, `FLASH_CARD_SIZE` vocabs, no vocab skipped, and one `flash_card` question per vocab (no options).

Constants:

- `FLASH_CARD_SIZE = 20`

## Submit

`POST /flash-card/submit`

### Input

- question id
- correct: the user's self-grade

```json
{
  "question_id": "...",
  "correct": false
}
```

### Output

- retry: when missed, a new question for the same vocab, appended to the end of the session. `null` when correct or out of attempts.
- session_completed: `true` once no unanswered questions are left

```json
{
  "retry": { "id": "...", "dict_form": "読む", "reading": "よむ", "translation": "to read" },
  "session_completed": false
}
```

### Flow

1. Load the question with its session. If it doesn't exist or belongs to another user, 404.
2. Session closed (expired) → 409
3. Already answered → return the stored result unchanged
4. Store `correct`, `answered_at = now`
5. If missed and the vocab has fewer than `MAX_ATTEMPTS` questions in this session, append a new question for it at the end
6. If no unanswered questions are left, close the session
7. Return

Every attempt counts for scheduling. Missed then correct in the same session gives streak 1, due in 1 day.

Constants:

- `MAX_ATTEMPTS = 3` per vocab per session

## Future

- Show one example sentence on the back after flipping (not graded), as a reminder of where the word came from
