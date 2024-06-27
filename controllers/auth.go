package controllers

import (
	"backend/models"
	"backend/models/token"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Checks if token specified in the Authorization header is valid
func Auth(accessSecret string) func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken, err := token.ExtractBearerToken(c)
		if err != nil {
			models.ErrJsonResponce(c, http.StatusUnauthorized, models.ErrNoTokenSpecified)
			return
		}

		_, err = token.AccessFromString(accessToken, accessSecret)

		if errors.Is(err, jwt.ErrTokenExpired) {
			models.ErrJsonResponce(c, http.StatusUnauthorized, models.ErrTokenExpired)
			return
		}

		if err != nil {
			models.ErrJsonResponce(c, http.StatusUnauthorized, models.ErrInvalidToken)
			return
		}

		c.Status(http.StatusOK)
	}
}
