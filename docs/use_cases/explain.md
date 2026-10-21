# Explain

This is very simple: send sentence and LLM explains it in the base language.

The LLM can also explain grammar, particles, idioms, set expressions and other things that are not in the translation from [analyze_sentence](analyze_sentence.md).

## Prerequisites

[Analyze Sentence](analyze_sentence.md) use case has been completed, and the user has tapped Explain.

## Input

- course id: determines the learning language and base language
- sentence: the normalized sentence from Analyze Sentence, unchanged
- sentence_translation: from Analyze Sentence, so the explanation stays consistent with the translation the user already saw

```json
{
  "course_id": "123",
  "sentence": "気をつけてください。",
  "sentence_translation": "Please be careful."
}
```

## Output

- explanation: markdown, in the base language

```json
{
  "explanation": "## Meaning\nA polite request asking someone to be careful.\n\n## Grammar\n..."
}
```

Expected latency is roughly 3 to 8 seconds (about 300 to 500 output tokens), instant on a cache hit. The client shows a loading state.

## Flow

1. Check the course exists and belongs to the user, otherwise 404
2. Validate, otherwise 400
   1. `sentence` is already normalized (normalizing it again changes nothing)
   2. `sentence` is non-empty and at most `MAX_SENTENCE_LENGTH` runes (same limit as Analyze Sentence)
3. Cache lookup by (`PROMPT_VERSION`, model, learning language, base language, sentence). On hit, return it.
4. Rate limit: shares `LLM_REQUESTS_PER_HOUR` with Analyze Sentence, otherwise 429. Only cache misses count.
5. LLM call with the request context, so a client disconnect cancels the LLM call. LLM failure → 502.
6. Store in cache and return

Cache and rate limit work the same as in [Analyze Sentence](analyze_sentence.md) (Redis, 30 day expiry, a Redis failure never fails the request), with key `explain:{PROMPT_VERSION}:{model}:{learning_language}:{base_language}:{sha256 hex of sentence}`. `PROMPT_VERSION` is Explain's own, bumped whenever its prompt changes.

## Future (streaming)

If latency turns out to be a problem, stream the explanation with Server-Sent Events (Gin `c.Stream` + `c.SSEvent`, openai-go `Chat.Completions.NewStreaming`). The prompt and cache stay the same. Things to handle then: errors after the first chunk can't change the HTTP status (send an `error` event), chunks sent as JSON (markdown newlines break raw SSE data lines), `WriteTimeout`, and proxies buffering the response.

## System Prompt

`{{learning_language}}` and `{{base_language}}` are filled from the course (language names, e.g. "Japanese", not codes). The user message contains the sentence and its translation.

```text
You are a friendly language tutor. The learner is studying {{learning_language}} and wants explanations in {{base_language}}.
They will send a sentence in {{learning_language}} (usually extracted by OCR from a book or a screen) and its translation.

Explain the sentence in {{base_language}}, in Markdown, using only the sections that apply:

## Meaning
What the sentence means here, and how it differs from a word-for-word reading. Stay consistent with the given translation unless it is wrong; if so, say so briefly.

## Grammar
Each grammar point that matters: conjugations, particles, sentence patterns. Name the form (e.g. "て-form + ください"), what it does, and the dictionary form of inflected words.

## Expressions
Idioms or set phrases (e.g. 気をつける), with their meaning as a whole.

## Nuance
Politeness level, tone, register, or anything a learner might misread.

Rules:
- Write readings for words beyond beginner level (Japanese: hiragana in parentheses, Chinese: pinyin with tone marks).
- Be concise: short paragraphs and bullet points, no introduction or closing remarks.
- If the text looks garbled by OCR, say what it most likely says before explaining.
```

User message:

```text
Sentence: {{sentence}}
Translation: {{sentence_translation}}
```

## Out of scope

- Explaining a single token or part (focus span). Whole sentence only for now.
- Follow-up questions (chat)
- Storing explanations with a saved occurrence. Tapping Explain again hits the cache.
