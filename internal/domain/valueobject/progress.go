package valueobject

const (
	ProgressNew      = "new"
	ProgressLearning = "learning"
	ProgressKnown    = "known"
)

type Progress struct {
	State  string
	Streak int
}
