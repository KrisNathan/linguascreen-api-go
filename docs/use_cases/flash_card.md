# Flash Card

Self-graded recall of the user's saved vocabs, without a sentence. Works for every vocab, including vocabs only saved from parts. Online only.

No sessions: each self-grade is stored as a `FlashCardReview`, the client keeps the deck. The server only picks the vocabs and records reviews.

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

- cards: vocabs to practice, in order. The back of the card is included, there is nothing to hide in a self-graded exercise.

```json
{
  "cards": [
    { "vocab_id": "...", "dict_form": "読む", "reading": "よむ", "translation": "to read" }
  ]
}
```

No vocab available returns `cards: []`.

### Flow

1. Check the course exists and belongs to the user, otherwise 404
2. Pick up to `FLASH_CARD_SIZE` vocabs ([Scheduling](scheduling.md)), no vocab skipped
3. Return

Constants:

- `FLASH_CARD_SIZE = 20`

## Submit

`POST /flash-card/submit`

### Input

- vocab id
- correct: the user's self-grade

```json
{
  "vocab_id": "...",
  "correct": false
}
```

### Output

`204 No Content`

### Flow

1. Load the vocab. If it doesn't exist or belongs to another user, 404.
2. Insert a `FlashCardReview` with `correct` and `reviewed_at = now`

Every review updates the vocab's `xp` and `last_reviewed_at` like a quiz answer ([Scheduling](scheduling.md)). Re-showing missed cards within the same sitting is up to the client.

## Future

- Show one example sentence on the back after flipping (not graded), as a reminder of where the word came from
