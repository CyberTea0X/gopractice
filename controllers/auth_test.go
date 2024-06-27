package controllers

import (
	"backend/models/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func AuthTDDTest(t *testing.T) {
	type AuthTest struct {
		request        *http.Request
		expectedStatus int
		name           string
	}
	secret := "secretTokenString"
	AuthTests := make([]AuthTest, 2)
	accessClaims := token.NewAccessClaims(123, []string{"user"}, time.Now().Add(time.Hour))
	accessToken, err := accessClaims.TokenString(secret)

	req, err := http.NewRequest("GET", BasePath+"/auth", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", accessToken)
	AuthTests[0] = AuthTest{req, http.StatusOK, "auth test with token"}
	req, err = http.NewRequest("GET", BasePath+"/auth", nil)
	AuthTests[1] = AuthTest{req, http.StatusUnauthorized, "auth test no token"}
	for _, test := range AuthTests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			engine := gin.New()
			engine.GET(BasePath+"/auth", Auth(secret))
			engine.ServeHTTP(w, test.request)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
