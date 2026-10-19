# Analyze Sentence

User selects a whole sentence (sentence detection is done by the frontend). We split it into tokens so the user can tap each word to see its translation and save it.

This replaces selecting individual words, which caused truncated selections (e.g. `走っ` instead of `走った`).

## Input

- course id: determines the learning language and base language
- sentence: as detected by the frontend, may contain OCR line breaks. At most `MAX_SENTENCE_LENGTH` runes (Unicode code points) after normalization.

```json
{
  "course_id": "123",
  "sentence": "よく本を\n読みます。"
}
```

## Output

- sentence: the normalized sentence. The client sends this back unchanged when saving a vocab.
- translation: the whole sentence translated into the base language
- language_mismatch: `true` if the sentence is not in the course's learning language (e.g. English text scanned in the Japanese course). The client warns before saving, since saved vocabs would land in the wrong course.
- tokens: the sentence split into tokens, in order. Concatenating every `surface` must equal `sentence` exactly.
  - surface: the token as seen in the sentence
  - span: `[start, end)` of the surface in `sentence`, in runes (Unicode code points, Go `[]rune` and Dart `String.runes`). Computed by the server, never by the LLM.
  - saveable: `false` for particles, punctuation, symbols and whitespace. Those tokens only have `surface` and `span`.
  - dict_form: dictionary form
  - reading: reading of the dict form (JP: hiragana, CN: pinyin with tone marks e.g. `yín háng`)
  - translation: short dictionary-style meaning that fits this sentence, in the base language (verbs as "to read")
  - parts: optional, when the token is made of smaller words worth learning (e.g. `火车站` → `火车` + `站`). The LLM decides when and how to split. Lets the user tap a sub-word without another request.
    - Each part has `surface`, `span` (absolute, in `sentence`), `dict_form`, `reading`, `translation`
    - Concatenating every part's `surface` must equal the token's `surface`
    - One level deep

```json
{
  "sentence": "よく本を読みます。",
  "translation": "I often read books.",
  "language_mismatch": false,
  "tokens": [
    { "surface": "よく", "span": [0, 2], "saveable": true, "dict_form": "よく", "reading": "よく", "translation": "often" },
    { "surface": "本", "span": [2, 3], "saveable": true, "dict_form": "本", "reading": "ほん", "translation": "book" },
    { "surface": "を", "span": [3, 4], "saveable": false },
    { "surface": "読みます", "span": [4, 8], "saveable": true, "dict_form": "読む", "reading": "よむ", "translation": "to read" },
    { "surface": "。", "span": [8, 9], "saveable": false }
  ]
}
```

Compound token with parts:

```json
{
  "surface": "火车站", "span": [4, 7], "saveable": true,
  "dict_form": "火车站", "reading": "huǒ chē zhàn", "translation": "train station",
  "parts": [
    { "surface": "火车", "span": [4, 6], "dict_form": "火车", "reading": "huǒ chē", "translation": "train" },
    { "surface": "站", "span": [6, 7], "dict_form": "站", "reading": "zhàn", "translation": "station" }
  ]
}
```

## Flow

For now:

1. Normalize sentence
   1. Unicode NFC (`golang.org/x/text/unicode/norm`). Not NFKC, it changes what the user read (e.g. `！` → `!`).
   2. Remove invisible characters: control characters and `unicode.Cf` (zero-width space, BOM), except zero-width joiner (U+200D, used by emoji)
   3. Line breaks: removed between two CJK characters, otherwise replaced with a space
   4. Remove spaces between two CJK characters (ML Kit sometimes inserts them)
   5. Collapse whitespace runs into one space and trim
2. Validate, otherwise 400
   1. Normalized sentence is non-empty
   2. At most `MAX_SENTENCE_LENGTH` runes (OCR can send a whole page, and every token costs money)
3. Check the course exists and belongs to the user, otherwise 404
4. Cache lookup by (`PROMPT_VERSION`, model, learning language, base language, normalized sentence). On hit, return it.
5. Rate limit: at most `LLM_REQUESTS_PER_HOUR` per user, shared with Explain, otherwise 429. Only cache misses count.
6. One LLM call with structured output (JSON schema), temperature 0. Input: sentence, learning language, base language. Output: sentence translation, `language_mismatch` and tokens (`surface`, `saveable`, `dict_form`, `reading`, `translation`, `parts`). No spans.
7. Align and repair (repair instead of reject, so the only failure is the LLM changing the text)
   1. Tokens and parts with an empty `surface` are dropped
   2. Find each token's `surface` in the sentence in order, searching forward from the end of the previous token. This gives each token's `span`. Uncovered characters before the first token, between tokens and after the last token (e.g. a skipped space or punctuation mark) become `saveable: false` tokens.
   3. A saveable token missing `dict_form`, `reading` or `translation` becomes `saveable: false`
   4. Parts whose `surface`s don't concatenate to their token's `surface` are dropped. Otherwise their `span`s are computed within the token.
   5. If a surface can't be found (the LLM altered the text), retry the LLM call once, then return an error (502)
8. Store in cache and return

Constants:

- `MAX_SENTENCE_LENGTH = 300`
- `LLM_REQUESTS_PER_HOUR = 200`
- `PROMPT_VERSION`: bumped whenever the prompt or the output schema changes, so old cached results are no longer served. Changing the model changes the key too.

Cache (Redis, also used by Explain):

- Shared by all users (results contain no personal data)
- Key: `analyze:{PROMPT_VERSION}:{model}:{learning_language}:{base_language}:{sha256 hex of sentence}`. The sentence is hashed to keep keys short.
- Value: the response JSON, stored with a 30 day expiry (`SET ... EX`)
- Redis unavailable: skip the cache and call the LLM. A cache failure never fails the request.

Rate limit (Redis, shared with Explain):

- Fixed one-hour window per user: `ratelimit:llm:{user_id}:{unix hour}`, `INCR`, and set a one-hour expiry when the counter is created
- Checked and incremented only on a cache miss, right before the LLM call
- Over `LLM_REQUESTS_PER_HOUR` → 429 with `Retry-After` (seconds until the window ends)
- Redis unavailable: allow the request and log it (fail open). Acceptable while user numbers are small.

### System prompt

`{{learning_language}}` and `{{base_language}}` are filled from the course. The user message is the normalized sentence only. The JSON schema is passed via structured output, not in the prompt.

```text
You are a language-learning assistant. The user is learning {{learning_language}} and their native language is {{base_language}}. They will send you one sentence in {{learning_language}}, usually extracted by OCR from a book or a screen.

Your job:
1. Translate the whole sentence into {{base_language}}.
2. Split the sentence into tokens so the learner can tap any word to see what it means.
3. Set language_mismatch to true if the sentence is mostly not in {{learning_language}}. Still translate and tokenize it.

Tokens:
- Copy the sentence exactly. The surfaces of all tokens, joined in order, must reproduce the sentence character for character. Do not fix typos or OCR errors, do not convert characters, punctuation or width (keep "！" as "！"), do not add or drop spaces. Punctuation, symbols and spaces are tokens too.
- A token is a unit a learner would look up in a dictionary. Keep an inflected word together with its inflection and auxiliaries as one token (e.g. 読みます, 走った, 食べさせられなかった), not split into grammatical pieces.
- saveable is false for punctuation, symbols, spaces, numbers, and grammatical function words that carry no dictionary meaning on their own (e.g. Japanese particles like を, が, に, Chinese 的, 了). Otherwise true. Only saveable tokens need dict_form, reading and translation.

For each saveable token:
- dict_form: the dictionary (lemma) form, e.g. 読みます → 読む, 走った → 走る. For languages that do not inflect, the word itself.
- reading: the pronunciation of dict_form.
  - Japanese: hiragana. For words written only in kana, repeat dict_form.
  - Chinese: pinyin with tone marks, syllables separated by spaces (e.g. yín háng). Pick the pronunciation that fits this sentence (e.g. 长 is cháng in 很长, zhǎng in 长大).
  - Languages without a separate reading system: empty string.
- translation: a short dictionary-style meaning in {{base_language}} for the sense used in this sentence, 1 to 4 words, no explanations. In English, write verbs as "to ..." (e.g. "to read").

Parts:
- If a token is a compound made of smaller words that are worth learning on their own, add parts (e.g. 火车站 → 火车 + 站, 国際連合 → 国際 + 連合). Otherwise omit parts.
- Do not split inflections or grammatical endings into parts.
- The part surfaces, joined in order, must reproduce the token's surface exactly.
- Each part has its own dict_form, reading and translation, following the rules above. Parts have no parts.

Sentence translation:
- Natural and faithful, in {{base_language}}. If the text looks garbled by OCR, translate your best interpretation, but still copy the token surfaces exactly as given.
```

Out of scope:

- Explanations (grammar, nuance, idioms and set expressions): separate Explain use case.
- Saving: separate Save Vocab use case, which takes `sentence` and the `span` of a token or part.

### Future (custom range selection):

This is not good UX.

Fallback for anything tokens and parts don't cover (e.g. the overlapping `车站` in `火车站`). User long-presses and drags handles at character granularity.

```
POST /analyze/range { course_id, sentence (normalized, as returned by /analyze), span }
-> { surface, span, kind, dict_form, reading, translation }
```

- The range is sent to the LLM marked in the text (`他跑到了【火车】站。`), not as offsets
- `kind`: `word` (saveable), `partial` (range cuts a word, returns a `suggested_surface` e.g. `走った` for `走っ`, user decides), `phrase` (translate only, not saveable)
- Shown as a popup, does not modify the cached analysis

### Future (proper dict lookup, JMdict / CC-CEDICT):

After step 4, look up each saveable token's `dict_form` in the dictionary (matching all listed spellings and readings, not only the headword).

1. If it exists: attach the dictionary entry id, and use the dictionary's reading.
2. If not: keep the LLM output, and mark the token as not in the dictionary.

The API shape stays the same, plus an optional `entry_id` per token.
