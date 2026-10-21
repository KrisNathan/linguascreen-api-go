package errs

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
	ErrUpstream     = errors.New("upstream error")
)

// hintError adds a message that is safe to show clients.
type hintError struct {
	err  error
	hint string
}

func (e *hintError) Error() string { return e.err.Error() + ": " + e.hint }
func (e *hintError) Unwrap() error { return e.err }

func WithHint(err error, hint string) error { return &hintError{err: err, hint: hint} }

func Hint(err error) string {
	var h *hintError
	if errors.As(err, &h) {
		return h.hint
	}
	return ""
}
