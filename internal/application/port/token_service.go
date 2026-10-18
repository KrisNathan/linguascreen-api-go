package port

import "time"

type TokenService interface {
	IssueAccess(userID string) (string, error)
	AccessTTL() time.Duration
	Parse(token string) (userID string, err error)
}
