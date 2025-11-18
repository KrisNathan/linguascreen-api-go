package repository

import (
	"database/sql"
	"linguascreen/models"
)

type MemoryRepository struct {
	db *sql.DB
}

func NewMemoryRepository(db *sql.DB) *MemoryRepository {
	return &MemoryRepository{db: db}
}

func (r *MemoryRepository) NewMemory(userID int, memory models.MemoryInput) error {
	query := "INSERT INTO memories (user_id, content_type, title, translation, explanation, example_sentence) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := r.db.Exec(query, userID, memory.ContentType, memory.Title, memory.Translation, memory.Explanation, memory.ExampleSentence)
	return err
}

func (r *MemoryRepository) GetMemories(userID int) ([]models.Memory, error) {
	query := "SELECT id, content_type, title, translation, explanation, example_sentence, memory_strength, last_reviewed_at, created_at FROM memories WHERE user_id = ?"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []models.Memory
	for rows.Next() {
		var memory models.Memory
		err := rows.Scan(&memory.ID, &memory.ContentType, &memory.Title, &memory.Translation, &memory.Explanation, &memory.ExampleSentence, &memory.MemoryStrength, &memory.LastReviewedAt, &memory.CreatedAt)
		if err != nil {
			return nil, err
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

func (r *MemoryRepository) GetWeakestMemories(userID int, limit int) ([]models.Memory, error) {
	query := "SELECT id, content_type, title, translation, explanation, example_sentence, memory_strength, last_reviewed_at, created_at FROM memories WHERE user_id = ? ORDER BY memory_strength ASC LIMIT ?"

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []models.Memory
	for rows.Next() {
		var memory models.Memory
		err := rows.Scan(&memory.ID, &memory.ContentType, &memory.Title, &memory.Translation, &memory.Explanation, &memory.ExampleSentence, &memory.MemoryStrength, &memory.LastReviewedAt, &memory.CreatedAt)
		if err != nil {
			return nil, err
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

func (r *MemoryRepository) UpdateMemoryStrength(memoryID int, newStrength int) error {
	query := "UPDATE memories SET memory_strength = ? WHERE id = ?"

	_, err := r.db.Exec(query, newStrength, memoryID)
	return err
}

func (r *MemoryRepository) UpdateLastReviewedAt(memoryID int) error {
	query := "UPDATE memories SET last_reviewed_at = CURRENT_TIMESTAMP WHERE id = ?"

	_, err := r.db.Exec(query, memoryID)
	return err
}

func (r *MemoryRepository) DeleteMemory(memoryID int) error {
	query := "DELETE FROM memories WHERE id = ?"

	_, err := r.db.Exec(query, memoryID)
	return err
}

func (r *MemoryRepository) GetRandomMemoriesByContentType(contentType string, excludeID int, limit int) ([]models.Memory, error) {
	query := "SELECT id, content_type, title, translation, explanation, example_sentence, memory_strength, last_reviewed_at, created_at FROM memories WHERE content_type = ? AND id != ? ORDER BY RAND() LIMIT ?"

	rows, err := r.db.Query(query, contentType, excludeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []models.Memory
	for rows.Next() {
		var memory models.Memory
		err := rows.Scan(&memory.ID, &memory.ContentType, &memory.Title, &memory.Translation, &memory.Explanation, &memory.ExampleSentence, &memory.MemoryStrength, &memory.LastReviewedAt, &memory.CreatedAt)
		if err != nil {
			return nil, err
		}
		memories = append(memories, memory)
	}

	return memories, nil
}
