# Create Course

## Create

`POST /courses`

### Input

- Base language: the language the user is fluent in (e.g. English)
- Learning language: the language the user wants to learn (e.g. Japanese/Chinese)

Both must be in ISO 639-1 format (e.g. "en", "ja", "zh").

```json
{
  "base_language": "en",
  "learning_language": "ja"
}
```

> Duplicate courses are not allowed. If a course already exists for the given base and learning language, return the existing course ID instead of creating a new one.

### Output

- Course ID: a unique identifier for the course

```json
{
  "course_id": "123",
  "base_language": "en",
  "learning_language": "ja"
}
```

### Flow

1. Validate both languages are supported ISO 639-1 codes and differ from each other, otherwise 400
2. Insert, or return the existing course on conflict (user, learning language, base language)
3. Return

## List

`GET /courses`

The client calls this after signing in (e.g. after a reinstall) to find the user's courses. No courses means the client shows onboarding.

### Output

Ordered by created at, oldest first. Not paginated, a user has only a few courses.

```json
{
  "courses": [
    { "course_id": "123", "base_language": "en", "learning_language": "ja", "created_at": "2026-10-18T10:00:00Z" }
  ]
}
```

## Language names

LLM prompts get language names (e.g. "Japanese"), not codes. The server keeps a small map of supported codes to names, which is also what Create validates against.
