# AGENTS.md

## Workflow

- Propose a solution first. Make no edits until the user accepts.
- If a new file does not fit an existing folder, name the target folder in the proposal.

## Skills

- /find-docs: before using any library API.
- /asd-ste100: for all responses to the user, docs, LLM prompts, and error messages.

## Architecture

Target: DDD with CQRS-lite (modeled on wild-workouts-go-ddd-example).

The current code does not follow this yet. It still has structure from the prototype. Write new code to the target. Move old code to the target only when the task touches it.

- `domain/`: aggregates, entities, value objects, repo interfaces. No infra imports.
- `domain/repo/`: takes and returns aggregates only.
- `application/query/`: read models and view types. Reads may join across aggregates in SQL.
- `application/usecase/`: one use case per file.
- `infra/`: adapters. `infra/repo/Pg*Repo` implements both the repo and the read model.
- Each folder holds one kind of file. `infra/http/handler/` holds handlers only (no DTOs or helpers).
- Display data goes into a view, never onto an aggregate.

## Code style

- Keep it simple.
- Comment only what the code cannot say: why, gotchas, magic values. No doc comments that repeat the signature.
