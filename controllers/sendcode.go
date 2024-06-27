package controllers

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SendCodeInput struct {
	Phone    string `form:"phone" binding:"required"`
	DeviceId *uint  `form:"device_id" binding:"required"`
}

func SendCode(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var input SendCodeInput

		if err := c.ShouldBindQuery(&input); err != nil {
			models.ErrJsonResponce(c, http.StatusBadRequest, models.ErrInvalidQuery)
			return
		}
		log.Println(input.Phone)
		sentAt, err := models.LastCodeSentAt(db, input.Phone)
		if err == nil {
			available := sentAt.Add(time.Minute)
			if !time.Now().After(available) {
				c.JSON(http.StatusConflict, models.ErrorResponce{
					Error:  models.ErrCodeAlreadySent.Error(),
					Status: http.StatusConflict,
					Data:   gin.H{"available_after": available.Unix()},
				})
				return
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		code := rand.Intn(10000) + 1000
		if err = models.SendCode(db, input.Phone, code); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				models.ErrJsonResponce(c, http.StatusBadRequest, models.ErrInvalidPhone)
				return
			}
			c.Status(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"code": fmt.Sprint(code)})
	}
}
