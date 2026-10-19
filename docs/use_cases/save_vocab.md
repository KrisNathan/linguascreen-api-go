# Save Vocab

## Prerequisites

[Analyze Sentence](analyze_sentence.md) use case has been completed, and the user has tapped a token or part to save it.

## Concepts

- Sentence: a normalized sentence in a course. Created when the first vocab is saved from it.
- Vocab: a word in the user's list, identified by (course, `dict_form`, `reading`).
- Occurrence: a vocab inside a sentence, identified by (vocab, sentence, span start). This is what the quiz shows.

A vocab can have many occurrences (same word saved from different sentences), and a sentence can have many occurrences (different words saved from the same sentence).

## Input

All fields are copied from the Analyze Sentence response.

- course id
- sentence: the normalized sentence, unchanged
- sentence_translation: the whole sentence translation
- span: `[start, end)` of the tapped token or part, in runes (Unicode code points)
- from_part: `true` if the user tapped a part instead of a token
- dict_form, reading, translation: of the tapped token or part

`surface` is not sent. The server takes it from `sentence[span]` so it can never disagree with the span.

```json
{
  "course_id": "123",
  "sentence": "よく本を読みます。",
  "sentence_translation": "I often read books.",
  "span": [4, 8],
  "from_part": false,
  "dict_form": "読む",
  "reading": "よむ",
  "translation": "to read"
}
```

## Output

- vocab_created: `false` if the vocab already existed
- occurrence_created: `false` if this exact occurrence was already saved

```json
{
  "vocab_id": "...",
  "occurrence_id": "...",
  "vocab_created": true,
  "occurrence_created": true
}
```

The client uses the flags for feedback:

| vocab_created | occurrence_created | Message |
| --- | --- | --- |
| true | true | Saved |
| false | true | New example added |
| false | false | Already saved |

## Flow

1. Check the course exists and belongs to the user, otherwise 404
2. Validate, otherwise 400
   1. `sentence` is already normalized (normalizing it again changes nothing). Catches clients sending raw text.
   2. `span` is non-empty and within the sentence
   3. `dict_form` and `translation` are non-empty
3. Normalize the vocab fields, so the same word from different LLM calls maps to the same vocab
   1. `dict_form`, `reading`, `translation`: NFC, trim, collapse whitespace into one space
   2. `reading` in Chinese (pinyin): also lowercase, and one space between syllables. Splitting joined syllables (`huǒchē` → `huǒ chē`) is left out, the prompt asks for spaces.
4. In one transaction:
   1. Sentence: insert, or reuse the existing one on conflict (course, text hash). An existing sentence keeps its stored translation.
   2. Vocab: insert, or reuse the existing one on conflict (course, `dict_form`, `reading`). An existing vocab keeps its stored translation, the user may have edited it.
   3. Occurrence: insert, or reuse the existing one on conflict (vocab, sentence, span start)
   4. If the occurrence was created, set the vocab's `last_saved_at` to now (orders the vocab list, see View Vocabs)
5. Return the ids and flags

Every step is insert-or-reuse, so retries and double taps are harmless.

## Stored data

Sentence:

- id, course id, text (normalized), translation, created at
- text hash: SHA-256 of the UTF-8 text, 32 raw bytes (`BYTEA`)
- unique (course id, text hash)

Vocab:

- id, course id, dict form, reading, translation, created at
- translation is the first one saved, user-editable. Shown in the word list.
- last saved at: when the latest occurrence was added. A fact, not a status.
- unique (course id, dict form, reading)
- index (course id, last saved at desc, id desc)

Occurrence:

- id, vocab id, sentence id, span start, span end, surface, created at
- translation: the meaning in this sentence (e.g. 本 is "book" in one sentence, "main" in another). Shown in the quiz.
- from part: occurrences saved from a part are never used for fill-in-the-blank, only recognition. Blanking a piece of a compound (`火车[__]`) is not a fair question.
- unique (vocab id, sentence id, span start)

Not stored: status, due date, streak. These are derived from answered questions (see [Scheduling](scheduling.md)).

## Sentence hash

Computed by the application, not Postgres, so the database only does storage and indexing.

- Computed in one place only: the Sentence domain constructor normalizes, then hashes. The repository stores what it receives.
- `sha256.Sum256([]byte(normalizedText))`
- If the normalization rules change, existing hashes must be recomputed with a one-off script (normalize and re-hash every sentence).

## Out of scope

- Unsaving and editing a vocab: [Manage Vocab](manage_vocab.md)
- Storing explanations with an occurrence: Explain use case.
