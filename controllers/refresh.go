package controllers

import (
	"backend/models"
	"backend/models/token"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshOutput struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	RefreshExpires int64  `json:"refresh_expires"`
	AccessExpires  int64  `json:"access_expires"`
}

func Refresh(db *sql.DB, refreshLifespan, accessLifespan time.Duration, refreshSecret, accessSecret string) func(*gin.Context) {
	return func(c *gin.Context) {
		inputToken := c.Query("token")
		if inputToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrNoTokenSpecified})
			log.Println("BAAD")
			return
		}

		refreshClaims, err := token.RefreshFromString(inputToken, refreshSecret)

		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrTokenExpired})
			return
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidToken})
			log.Println("REEALLY BAAD")
			return
		}

		exists, err := refreshClaims.Exists(db)
		if err != nil {
			log.Println("Error while checking if refresh token exists: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		if exists == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrTokenNotExists})
			log.Println("SOOOOO BAAD")
			return
		}

		refreshExpires := time.Now().Add(refreshLifespan)
		refreshClaims.ExpiresAt = jwt.NewNumericDate(refreshExpires)

		_, err = refreshClaims.Update(db, uint64(refreshExpires.Unix()))

		if err != nil {
			log.Println("Error updating refresh token identifier in the database: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		refreshToken, err := refreshClaims.TokenString(refreshSecret)

		if err != nil {
			log.Println("Error generating refresh token from old refresh token: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		accessExpires := time.Now().Add(accessLifespan)
		accessClaims := token.NewAccessClaims(refreshClaims.UserID, refreshClaims.Roles, accessExpires)
		accessToken, err := accessClaims.TokenString(accessSecret)

		if err != nil {
			log.Println("Error generating access token from old refresh token: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, RefreshOutput{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			RefreshExpires: refreshExpires.Unix(),
			AccessExpires:  accessExpires.Unix(),
		})
	}

}
