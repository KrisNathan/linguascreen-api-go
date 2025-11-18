DROP TABLE IF EXISTS memories;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password_hash VARCHAR(256) NOT NULL,
    original_language VARCHAR(64) NOT NULL,
    target_language VARCHAR(64) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users ADD CONSTRAINT username_unique UNIQUE (username);

INSERT INTO users (username, password_hash, original_language, target_language) VALUES
('root', 'root', 'eng', 'jpn');

CREATE TABLE memories (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    content_type VARCHAR(32) NOT NULL,
    title VARCHAR(128) NOT NULL,
    translation VARCHAR(256) NOT NULL,
    explanation VARCHAR(512),
    example_sentence VARCHAR(512),
    memory_strength INT DEFAULT 0,
    last_reviewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
