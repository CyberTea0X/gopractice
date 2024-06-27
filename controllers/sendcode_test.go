package controllers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendCode(t *testing.T) {
	engine, db, err := SetupServer("../config_test.toml", "../users.json")
	assert.NoError(t, err)
	defer func() {
		if err := models.CleanDatabase(db); err != nil {
			t.Fatal(err)
		}
	}()
	deviceId := uint(1)
	phone := "+1234567890"
	roles := []string{"user"}
	_, err = models.CreateOrUpdateUser(db, phone, roles)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", BasePath+"/sendcode", nil)
	assert.NoError(t, err)
	query := req.URL.Query()
	query.Set("phone", phone)
	query.Set("device_id", fmt.Sprint(deviceId))
	req.URL.RawQuery = query.Encode()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, w.Body)

	var output map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &output)
	assert.NoError(t, err)

	assert.NotEmpty(t, output["code"])
	// Clean up the test database
	assert.NoError(t, models.CleanDatabase(db))
}
