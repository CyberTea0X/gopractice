package models

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJson      = errors.New("invalid json")
	ErrNoTokenSpecified = errors.New("no token specified in Authorization header")
	ErrTokenExpired     = errors.New("token has expired")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenNotExists   = errors.New("token does not exist")
	ErrInvalidQuery     = errors.New("invalid query")
	ErrCodeAlreadySent  = errors.New("sms code already sent, wait before requesting another code")
	ErrInvalidCode      = errors.New("sms code is invalid")
	ErrCodeExpired      = errors.New("sms code has expired")
	ErrInvalidPhone     = errors.New("user with this phone does not exist")
)

func ErrJsonResponce(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponce{Error: err.Error(), Status: status})
}

type ErrorResponce struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
}
