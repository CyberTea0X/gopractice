package token

import (
	"backend/models"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// While parsing checks if token is valid by checking its sign method and signature
// signature is token specific secret key
func ParseWithClaims(token string, claims jwt.Claims, secret string) (*jwt.Token, error) {
	jwt, result := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if !jwt.Valid {
		return nil, result
	}
	return jwt, result
}

// Extracts bearer token from "Authorization" request header
//
// returns token string or models.ErrNoTokenSpecified
func ExtractBearerToken(c *gin.Context) (string, error) {
	auth := c.Request.Header.Get("Authorization")
	// Authorization usually starts on "Bearer", then comes single whitespace and then token string
	authSplit := strings.Split(auth, " ")
	if (len(authSplit) == 2) && (authSplit[1] != "") {
		return authSplit[1], nil
	} else {
		return "", models.ErrNoTokenSpecified
	}
}
