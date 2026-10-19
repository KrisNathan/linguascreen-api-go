package valueobject

import "time"

const (
	ProgressNew      = "new"
	ProgressLearning = "learning"
	ProgressKnown    = "known"
)

type Progress struct {
	State  string
	Streak int
	DueAt  *time.Time
}
