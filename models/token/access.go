package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	UserId int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}

func NewAccessClaims(userId int64, roles []string, expiresAt time.Time) *AccessClaims {
	t := new(AccessClaims)
	t.UserId = userId
	t.ExpiresAt = jwt.NewNumericDate(expiresAt)
	t.Roles = roles
	return t
}

// Parses token from token string
func AccessFromString(token string, secret string) (*AccessClaims, error) {
	t, err := ParseWithClaims(token, &AccessClaims{}, secret)
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(*AccessClaims); ok {
		return claims, nil
	}
	return nil, errors.New("unknown claims type, cannot proceed")
}

// Encodes RefreshToken model into a JWT string
func (c *AccessClaims) TokenString(secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token_string, err := token.SignedString([]byte(secret))

	return token_string, err
}
