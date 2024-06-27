package controllers

import (
	"backend/models"
	"backend/models/token"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const refreshPath = "/refresh"

// Generates refresh token for testing purposes
func generateTestRefresh(t *testing.T, refresh string, router *gin.Engine) *RefreshOutput {
	req, err := http.NewRequest("GET", BasePath+refreshPath, nil)
	assert.NoError(t, err)
	q := req.URL.Query()
	q.Add("token", refresh)
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	rawBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	output := new(RefreshOutput)
	err = json.Unmarshal(rawBody, output)
	assert.NoError(t, err, "Failed to parse refresh token output from body")
	return output
}

func TestRefreshSucceds(t *testing.T) {

	config, err := models.ParseConfig("../config_test.toml")
	assert.NoError(t, err)
	db, err := models.SetupDatabase(&config.Database)
	assert.NoError(t, err)
	err = models.MigrateDatabase(db)
	assert.NoError(t, err)

	router := gin.Default()
	router.GET(BasePath+"/refresh", Refresh(
		db,
		config.Tokens.Refresh.Lifespan(),
		config.Tokens.Access.Lifespan(),
		config.Tokens.Refresh.Secret,
		config.Tokens.Access.Secret,
	))

	defer func() {
		err := models.CleanDatabase(db)
		assert.NoError(t, err)
	}()

	// User does not exist, so we need to register them first
	phone := "+1234567890"
	deviceId := uint(12345)
	roles := []string{}

	userId, err := models.CreateOrUpdateUser(db, phone, roles)
	assert.NoError(t, err)

	refreshExpires := time.Now().Add(config.Tokens.Refresh.Lifespan())
	refreshClaims := token.NewRefreshClaims(deviceId, userId, roles, refreshExpires)

	tokenId, err := refreshClaims.InsertOrUpdate(db)
	assert.NoError(t, err)
	refreshClaims.TokenID = tokenId
	refreshToken, err := refreshClaims.TokenString(config.Tokens.Refresh.Secret)
	assert.NoError(t, err)

	refreshed := generateTestRefresh(t, refreshToken, router)
	r := generateTestRefresh(t, refreshed.RefreshToken, router)
	assert.NotEmpty(t, r.AccessToken)
	assert.NotEmpty(t, r.RefreshToken)
	assert.NotEmpty(t, r.RefreshExpires)
	assert.NotEmpty(t, r.AccessExpires)
}
