package controllers

import (
	"backend/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	// Setup test DB
	engine, db, err := SetupServer("../config_test.toml", "../users.json")
	assert.NoError(t, err)

	// User does not exist, so we need to register them first
	phone := "+1234567890"
	deviceId := uint(12345)
	roles := []string{}
	code := 1234

	userId, err := models.CreateOrUpdateUser(db, phone, roles)
	assert.NoError(t, err)

	err = models.SendCode(db, phone, code)
	assert.NoError(t, err)

	// Now that the user exists, attempt to log in
	loginInput := &LoginInput{
		Phone:    phone,
		DeviceId: &deviceId,
	}
	jsonData, err := json.Marshal(loginInput)
	assert.NoError(t, err)
	req := httptest.NewRequest("GET", BasePath+"/login", bytes.NewBuffer(jsonData))
	query := req.URL.Query()
	query.Add("phone", phone)
	query.Add("device_id", fmt.Sprint(deviceId))
	query.Add("code", fmt.Sprint(code))
	req.URL.RawQuery = query.Encode()
	w := httptest.NewRecorder()

	// Assuming Login is registered under /login path
	engine.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var output LoginOutput
	err = json.Unmarshal(w.Body.Bytes(), &output)
	assert.NoError(t, err)

	// Assertions about the response
	assert.NotEmpty(t, output.AccessToken)
	assert.NotEmpty(t, output.RefreshToken)
	assert.GreaterOrEqual(t, output.RefreshExpires, time.Now().Unix())
	assert.GreaterOrEqual(t, output.AccessExpires, time.Now().Unix())

	_, err = models.LastCode(db, userId)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	// Clean up the test database
	assert.NoError(t, models.CleanDatabase(db))
}
