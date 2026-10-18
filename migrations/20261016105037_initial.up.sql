CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  google_sub TEXT UNIQUE NOT NULL,
  email TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);

CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL,

  token_hash CHAR(64) UNIQUE NOT NULL,

  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- revoking all of a user's tokens, and cleaning up expired ones
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);

CREATE TABLE courses (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL,
  base_language TEXT NOT NULL,
  learning_language TEXT NOT NULL,
  daily_target INT NOT NULL DEFAULT 15,

  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
  CONSTRAINT uk_course UNIQUE (user_id, learning_language, base_language),
  CONSTRAINT ck_course_languages CHECK (base_language <> learning_language)
);

CREATE TABLE vocabs (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  course_id UUID NOT NULL,

  dict_form TEXT NOT NULL,
  reading TEXT NOT NULL, -- empty string for languages without a separate reading
  translation TEXT NOT NULL,

  last_saved_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

  CONSTRAINT fk_course_id
    FOREIGN KEY (course_id)
    REFERENCES courses(id)
    ON DELETE CASCADE,
  CONSTRAINT uk_vocab UNIQUE (course_id, dict_form, reading)
);

-- vocab list, ordered by last saved (keyset pagination)
CREATE INDEX idx_vocabs_course_last_saved ON vocabs (course_id, last_saved_at DESC, id DESC);

CREATE TABLE sentences (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  course_id UUID NOT NULL,

  text TEXT NOT NULL,
  text_hash TEXT NOT NULL,
  translation TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

  CONSTRAINT fk_course_id
    FOREIGN KEY (course_id)
    REFERENCES courses(id)
    ON DELETE CASCADE,
  CONSTRAINT uk_sentence UNIQUE (course_id, text_hash)
);

CREATE TABLE occurrences (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  vocab_id UUID NOT NULL,
  sentence_id UUID NOT NULL,

  span_start INT NOT NULL, -- UTF-16 code units
  span_end INT NOT NULL,   -- UTF-16 code units, exclusive
  surface TEXT NOT NULL,
  translation TEXT NOT NULL,
  from_part BOOL NOT NULL DEFAULT false,

  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),

  CONSTRAINT fk_vocab_id
    FOREIGN KEY (vocab_id)
    REFERENCES vocabs(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_sentence_id
    FOREIGN KEY (sentence_id)
    REFERENCES sentences(id)
    ON DELETE CASCADE,
  CONSTRAINT uk_occurrence UNIQUE (vocab_id, sentence_id, span_start),
  CONSTRAINT ck_occurrence_span CHECK (span_start >= 0 AND span_end > span_start)
);

-- uk_occurrence covers lookups by vocab_id, this covers lookups by sentence_id
CREATE INDEX idx_occurrences_sentence_id ON occurrences (sentence_id);

CREATE TABLE practice_sessions (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  course_id UUID NOT NULL,

  kind TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
  closed_at TIMESTAMPTZ, -- completed or expired, null while open

  CONSTRAINT fk_course_id
    FOREIGN KEY (course_id)
    REFERENCES courses(id)
    ON DELETE CASCADE,
  CONSTRAINT ck_practice_session_kind CHECK (kind IN ('quiz', 'flash_card'))
);

-- at most one open session per course and kind
CREATE UNIQUE INDEX uk_practice_sessions_open ON practice_sessions (course_id, kind) WHERE closed_at IS NULL;

CREATE TABLE questions (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  session_id UUID NOT NULL,
  position INT NOT NULL,

  type TEXT NOT NULL,
  vocab_id UUID NOT NULL,
  occurrence_id UUID, -- fill_blank only

  options TEXT[],      -- quiz only
  correct_index INT,   -- quiz only, never sent before answering
  answered_index INT,  -- quiz only, null until answered
  correct BOOL,        -- null until answered
  answered_at TIMESTAMPTZ, -- server time, null until answered

  CONSTRAINT fk_session_id
    FOREIGN KEY (session_id)
    REFERENCES practice_sessions(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_vocab_id
    FOREIGN KEY (vocab_id)
    REFERENCES vocabs(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_occurrence_id
    FOREIGN KEY (occurrence_id)
    REFERENCES occurrences(id)
    ON DELETE CASCADE,
  CONSTRAINT uk_question_position UNIQUE (session_id, position),
  CONSTRAINT ck_question_type CHECK (type IN ('fill_blank', 'translation', 'reading', 'flash_card')),
  CONSTRAINT ck_question_fill_blank CHECK ((type = 'fill_blank') = (occurrence_id IS NOT NULL)),
  CONSTRAINT ck_question_options CHECK (
    (type = 'flash_card' AND options IS NULL AND correct_index IS NULL AND answered_index IS NULL)
    OR (type <> 'flash_card' AND cardinality(options) = 4 AND correct_index BETWEEN 0 AND 3
        AND (answered_index IS NULL OR answered_index BETWEEN 0 AND 3))
  ),
  CONSTRAINT ck_question_answered CHECK ((answered_at IS NULL) = (correct IS NULL))
);

-- scheduling: answered questions per vocab
CREATE INDEX idx_questions_vocab_answered ON questions (vocab_id, answered_at) WHERE answered_at IS NOT NULL;
-- fill_blank rotation: latest quiz per occurrence
CREATE INDEX idx_questions_occurrence_answered ON questions (occurrence_id, answered_at) WHERE occurrence_id IS NOT NULL;
