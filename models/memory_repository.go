package models

type MemoryRepository interface {
	NewMemory(userID int, memory MemoryInput) error
	GetMemories(userID int) ([]Memory, error)
	GetWeakestMemories(userID int, limit int) ([]Memory, error)
	UpdateMemoryStrength(memoryID int, newStrength int) error
	UpdateLastReviewedAt(memoryID int) error
	DeleteMemory(memoryID int) error
	GetRandomMemoriesByContentType(contentType string, excludeID int, limit int) ([]Memory, error)
}
