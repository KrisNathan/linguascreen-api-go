# Quiz

Multiple choice questions about the user's saved vocabs. The server creates the session, keeps the correct answers, and grades every answer. Online only.

## Prerequisites

- User is signed in
- User has at least one course

## Question types

| Type | Shows | Options (1 correct, 3 wrong) | Eligible when |
| --- | --- | --- | --- |
| `fill_blank` | The sentence with the vocab blanked out | `dict_form` of other vocabs | The vocab has an occurrence with `from_part = false` |
| `translation` | `dict_form` and `reading` | `translation` of other vocabs | Always |
| `reading` | `dict_form` | `reading` of other vocabs | `reading` is non-empty and differs from `dict_form` (a kana-only word like `よく` would give the answer away) |

Wrong options:

- 3 random other vocabs from the course
- Never equal to the correct option or to each other (no two "to see", no duplicate readings, no same `dict_form`)
- Not enough distinct candidates for a type → that type is not eligible for this vocab

Each question's type is picked at random from the types its vocab is eligible for. A vocab with no eligible type is skipped (it can still be practiced with flash cards).

`fill_blank` uses the occurrence quizzed least recently (oldest latest `fill_blank` answer, never quizzed first, ties → newest occurrence), so it rotates through the sentences the user read. Options are dictionary forms, so for a conjugated word the blank is `読みます` while the option is `読む`.

## Start

`POST /quiz/start`

### Input

- course id

```json
{
  "course_id": "123"
}
```

### Output

- session_id
- score: so far in this session
- questions: the unanswered ones, in order. The correct option is never included.

```json
{
  "session_id": "...",
  "score": { "correct": 0, "answered": 0, "total": 3 },
  "questions": [
    {
      "id": "...",
      "type": "fill_blank",
      "sentence": "他跑到了火车站。",
      "span": [4, 7],
      "options": ["银行", "火车站", "明天", "问题"]
    },
    {
      "id": "...",
      "type": "translation",
      "dict_form": "読む",
      "reading": "よむ",
      "options": ["to eat", "to read", "to run", "book"]
    },
    {
      "id": "...",
      "type": "reading",
      "dict_form": "銀行",
      "options": ["ぎんこう", "でんしゃ", "がっこう", "えき"]
    }
  ]
}
```

No vocab available returns `questions: []` and no session is created. The client shows "nothing to practice".

### Flow

1. Check the course exists and belongs to the user, otherwise 404
2. Find the course's open quiz session (`closed_at` is null)
   1. Created less than 24 hours ago and has unanswered questions: return it with its unanswered questions
   2. Otherwise (expired, or nothing left after vocabs were unsaved): close it (`closed_at = now`). Unanswered questions are ignored by scheduling.
3. Pick up to `QUIZ_SIZE` vocabs ([Scheduling](scheduling.md)), skipping vocabs with no eligible question type
4. Build one question per vocab (type, options shuffled, `correct_index`)
5. In one transaction: insert the session and its questions
6. Return

A unique index on (course id, kind) where `closed_at` is null prevents two simultaneous starts from creating two open sessions.

Constants:

- `QUIZ_SIZE = 10`

## Submit

`POST /quiz/submit`

One question at a time. No retries of wrong answers within a quiz, so the score stays meaningful.

### Input

- question id
- option_index: the chosen option, 0 to 3

```json
{
  "question_id": "...",
  "option_index": 0
}
```

### Output

- correct
- correct_index
- vocab: `dict_form`, `reading`, `translation` for the feedback screen
- sentence_translation: `fill_blank` only. Shown after answering, it would give the answer away before.
- score: so far in this session
- session_completed: `true` after the last question

```json
{
  "correct": false,
  "correct_index": 1,
  "vocab": { "dict_form": "火车站", "reading": "huǒ chē zhàn", "translation": "train station" },
  "sentence_translation": "He ran to the train station.",
  "score": { "correct": 0, "answered": 1, "total": 3 },
  "session_completed": false
}
```

### Flow

1. Load the question with its session. If it doesn't exist or belongs to another user, 404.
2. Validate `option_index` (0 to 3), otherwise 400
3. Session closed (expired) → 409
4. Already answered → return the stored result unchanged. Retries are harmless and an answer can't be changed.
5. Grade: `correct = option_index == correct_index`. Store `answered_index`, `correct`, `answered_at = now`.
6. If it was the last unanswered question, close the session (`closed_at = now`)
7. Return

Score: `correct` = correct answers, `answered` = answered questions, `total` = questions in the session. Derived, not stored.

## Stored data

See `PracticeSession` and `Question` in the ERD. A session has `kind = quiz`. Options are stored as text, so editing a vocab's translation mid-session doesn't change a question already created.

## Future

- Typed answers instead of multiple choice (painful with Japanese and Chinese keyboards on a phone)
- Smarter wrong options (same part of speech, similar meaning, similar reading)
