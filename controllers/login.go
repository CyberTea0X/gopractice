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
)

type LoginInput struct {
	Phone    string `form:"phone" binding:"required"`
	DeviceId *uint  `form:"device_id" binding:"required"`
	Code     string `form:"code" binding:"required"`
}

type LoginOutput struct {
	AccessToken    string   `json:"access_token"`
	RefreshToken   string   `json:"refresh_token"`
	RefreshExpires int64    `json:"refresh_expires"`
	AccessExpires  int64    `json:"access_expires"`
	Roles          []string `json:"role"`
}

// Function that is responsible for user authorization.
//
// In response to a successful authorization request, returns
// access token and refresh token, as well as the time of death of the access token
func Login(db *sql.DB, refreshLifespan, accessLifespan, smsCodeLifespan time.Duration, refreshSecret, accessSecret string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input LoginInput

		if err := c.ShouldBindQuery(&input); err != nil {
			models.ErrJsonResponce(c, http.StatusBadRequest, models.ErrInvalidQuery)
			return
		}

		user, err := models.GetUserByPhone(db, input.Phone)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.Status(http.StatusUnauthorized)
			} else {
				c.Status(http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}

		code, err := models.LastCode(db, user.Id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				models.ErrJsonResponce(c, http.StatusUnauthorized, models.ErrInvalidCode)
				return
			} else {
				c.Status(http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}

		if code.SentAt.Add(smsCodeLifespan).Before(time.Now()) {
			models.ErrJsonResponce(c, http.StatusUnauthorized, models.ErrCodeExpired)
			return
		}

		if err := models.DeleteCode(db, user.Id); err != nil {
			c.Status(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		refreshExpires := time.Now().Add(refreshLifespan)
		refreshClaims := token.NewRefreshClaims(*input.DeviceId, user.Id, user.Roles, refreshExpires)

		refreshId, err := refreshClaims.InsertOrUpdate(db)

		if err != nil {
			log.Println("Error inserting or updating on duplicate key refresh token in the db: ", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		refreshClaims.TokenID = refreshId

		refreshToken, err := refreshClaims.TokenString(refreshSecret)

		if err != nil {
			log.Println("Error generating refresh token: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		accessExpires := time.Now().Add(accessLifespan)
		accessClaims := token.NewAccessClaims(user.Id, user.Roles, accessExpires)
		accessToken, err := accessClaims.TokenString(accessSecret)

		if err != nil {
			log.Println("Error generating access token: ", err.Error())
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, LoginOutput{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			AccessExpires:  accessExpires.Unix(),
			RefreshExpires: refreshExpires.Unix(),
			Roles:          user.Roles,
		})
	}
}
