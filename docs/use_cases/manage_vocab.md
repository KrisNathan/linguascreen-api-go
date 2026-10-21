# Manage Vocab

Edit a vocab's translation, or unsave a vocab.

## Prerequisites

- User is signed in
- The vocab exists (see [Save Vocab](save_vocab.md))

## Edit

`PATCH /vocabs/{vocab_id}`

Only the vocab's `translation` is editable (the one shown in the vocab list). `dict_form` and `reading` identify the vocab, so they can't change. A wrong one means unsave and save again.

Occurrence translations (meaning in a specific sentence) are not editable for now.

### Input

```json
{
  "translation": "to read (a book)"
}
```

### Output

The vocab, same shape as a [View Vocabs](view_vocabs.md) list item.

### Flow

1. Load the vocab. If it doesn't exist or belongs to another user, 404.
2. Validate: `translation` non-empty after normalizing (NFC, trim, collapse whitespace) and at most 200 characters, otherwise 400
3. Update and return

Questions already created keep their stored options, so an open quiz isn't affected.

## Unsave

`DELETE /vocabs/{vocab_id}`

Deletes the vocab with its occurrences and questions (`ON DELETE CASCADE`). Its practice history is gone, saving it again starts as new.

### Output

204. Also 204 if it was already deleted, so retries are harmless.

### Flow

1. Load the vocab. Doesn't exist → 204. Belongs to another user → 404.
2. Delete it (occurrences and questions cascade)
3. If it was in an open practice session, the session just has fewer questions. If none are left unanswered, the next submit or start closes it.

Sentences left without occurrences stay. They are harmless and can be cleaned up later.

## Out of scope

- Editing occurrence translations
- Deleting a single occurrence (example sentence) without unsaving the vocab
- Cleaning up sentences without occurrences
