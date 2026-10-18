package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const accessTTL = 15 * time.Minute

type JWTTokenService struct {
	signingKID string
	keys       map[string][]byte
}

// NewTokenService signs with keys[signingKID] and verifies with any key in keys.
func NewTokenService(signingKID string, keys map[string][]byte) (*JWTTokenService, error) {
	if _, ok := keys[signingKID]; !ok {
		return nil, fmt.Errorf("signing kid %q not found in keys", signingKID)
	}
	return &JWTTokenService{signingKID: signingKID, keys: keys}, nil
}

// ParseKeys parses "kid:secret,kid:secret" into a key map.
func ParseKeys(s string) (map[string][]byte, error) {
	keys := map[string][]byte{}
	for pair := range strings.SplitSeq(s, ",") {
		kid, secret, ok := strings.Cut(strings.TrimSpace(pair), ":")
		if !ok || kid == "" || secret == "" {
			return nil, errors.New("invalid keys, expected kid:secret,kid:secret")
		}
		keys[kid] = []byte(secret)
	}
	return keys, nil
}

func (s *JWTTokenService) IssueAccess(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
	})
	token.Header["kid"] = s.signingKID
	return token.SignedString(s.keys[s.signingKID])
}

func (*JWTTokenService) AccessTTL() time.Duration {
	return accessTTL
}

func (s *JWTTokenService) Parse(tokenString string) (string, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		key, ok := s.keys[kid]
		if !ok {
			return nil, fmt.Errorf("unknown key id %q", kid)
		}
		return key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return "", err
	}
	return tok.Claims.(*jwt.RegisteredClaims).Subject, nil
}
