# View Vocabs

Two endpoints: the vocab list (lazy loaded) and a single vocab with all its occurrences.

## Prerequisites

- User is signed in
- User has at least one course

## List

### Input

- course id
- limit: page size, default 30, max 100
- cursor: from the previous page's `next_cursor`, omitted for the first page

```
GET /courses/{course_id}/vocabs?limit=30&cursor=<opaque>
```

### Output

- items: vocabs ordered by last saved, newest first
  - progress (see [Scheduling](scheduling.md))
    - state: `new` (no answers) | `learning` | `known`
    - streak: correct answers in a row since the last miss
- next_cursor: `null` on the last page

```json
{
  "items": [
    {
      "id": "...",
      "dict_form": "読む",
      "reading": "よむ",
      "translation": "to read",
      "occurrence_count": 3,
      "last_saved_at": "2026-10-18T10:00:00Z",
      "progress": { "state": "learning", "streak": 2 }
    }
  ],
  "next_cursor": "..."
}
```

No saved vocabs returns `items: []` and `next_cursor: null`, not an error.

### Order

Last saved: when the latest occurrence was added (see Save Vocab). Saving a word again from a new sentence moves it to the top.

Not ordered by:

- Created: saving a word again from a new sentence would leave it buried
- Any update: the whole list would reshuffle after every quiz or edit
- Progress (weakest): duplicates the quiz. Possible later as a filter.

### Pagination

Cursor (keyset) pagination, not offset. New vocabs saved while scrolling don't shift the pages, so no item is shown twice, and deep pages stay fast.

- Sort key: `(last_saved_at DESC, id DESC)`. `id` breaks ties between identical timestamps.
- Cursor: `(last_saved_at, id)` of the last item, encoded as an opaque base64url string. The client never parses it.
- Next page: `WHERE course_id = ? AND (last_saved_at, id) < (cursor.last_saved_at, cursor.id)`
- Index: `(course_id, last_saved_at DESC, id DESC)`
- A vocab saved again mid-scroll jumps to the top. The cursor never repeats or skips the other rows.

### Flow

1. Check the course exists and belongs to the user, otherwise 404
2. Validate `limit` (1 to 100) and decode `cursor`, otherwise 400
3. Fetch `limit + 1` vocabs. If the extra row exists there is a next page: drop it, and `next_cursor` is built from the last returned row.
4. Count occurrences for the returned vocabs only, in one grouped query

## Single vocab

### Input

```
GET /vocabs/{vocab_id}
```

### Output

The vocab (same fields as a list item) plus all its occurrences, newest first. Not paginated, a vocab rarely has more than a handful.

```json
{
  "id": "...",
  "dict_form": "読む",
  "reading": "よむ",
  "translation": "to read",
  "last_saved_at": "2026-10-18T10:00:00Z",
  "progress": { "state": "learning", "streak": 2 },
  "occurrences": [
    {
      "id": "...",
      "sentence": "よく本を読みます。",
      "sentence_translation": "I often read books.",
      "surface": "読みます",
      "span": [4, 8],
      "translation": "to read",
      "from_part": false,
      "created_at": "2026-10-18T10:00:00Z"
    }
  ]
}
```

### Flow

1. Load the vocab. If it doesn't exist or its course doesn't belong to the user, 404 (same response for both, so ids can't be probed).
2. Load its occurrences with their sentences

## Out of scope

- Search (dict form, reading, translation) and filters (due, known). With search the cursor still works, it only adds a `WHERE` condition.
- Editing and deleting a vocab: [Manage Vocab](manage_vocab.md)
