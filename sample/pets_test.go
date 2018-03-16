package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/miketonks/swag-validator/sample/server"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)

	code := m.Run()
	os.Exit(code)
}

func TestGetPet(t *testing.T) {

	api := server.SetupAPI()
	router := server.SetupRouter(api)

	resp := getReq(t, router, "/pet/1234")
	assert.Equal(t, 200, resp.Code)

	resp = getReq(t, router, "/pet/ollie")
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "petId: Invalid type. Expected: integer, given: string")
}

func TestPostPet(t *testing.T) {

	api := server.SetupAPI()
	router := server.SetupRouter(api)

	t.Run("POST with a valid payload", func(t *testing.T) {
		validPet := `{"name": "ollie", "dob": "2018-01-01T12:00:00-09:00", "grumpy": true, "uuid": "1c694c09-3210-45d4-be6b-dbd94be1be4f"}`
		resp := postStr(t, router, "/pet", validPet)

		assert.Equal(t, 200, resp.Code)
	})

	requiredFieldMissing := `{"dob": "2018-01-01T12:00:00-09:00", "grumpy": true, "uuid": "1c694c09-3210-45d4-be6b-dbd94be1be4f"}`
	resp := postStr(t, router, "/pet", requiredFieldMissing)
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "name is required")

	invalidDateTime := `{"name": "ollie", "dob": "2018-01-01T12:00:0"}`
	resp = postStr(t, router, "/pet", invalidDateTime)
	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "Does not match format 'date-time'")

	invalidUUID := `{"name": "ollie", "uuid": "foo-bar-baz"}`
	resp = postStr(t, router, "/pet", invalidUUID)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "Does not match format 'uuid'")

	invalidBool := `{"name": "ollie", "grumpy": "very"}`
	resp = postStr(t, router, "/pet", invalidBool)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid type. Expected: boolean, given: string")

	validDate := `{"name": "ollie", "dt": "2018-04-01"}`
	resp = postStr(t, router, "/pet", validDate)
	assert.Equal(t, 200, resp.Code)

	invalidDate := `{"name": "ollie", "dt": "2018-44-01"}`
	resp = postStr(t, router, "/pet", invalidDate)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "month out of range")

	validTime := `{"name": "ollie", "tm": "12:15:00"}`
	resp = postStr(t, router, "/pet", validTime)
	assert.Equal(t, 200, resp.Code)

	invalidTime := `{"name": "ollie", "tm": "12:15:99"}`
	resp = postStr(t, router, "/pet", invalidTime)

	assert.Equal(t, 400, resp.Code)
	assert.Contains(t, resp.Body.String(), "second out of range")
}

func TestNotFound(t *testing.T) {

	api := server.SetupAPI()
	router := server.SetupRouter(api)

	resp := getReq(t, router, "/foo")
	assert.Equal(t, 404, resp.Code)
}

func getReq(t *testing.T, router *gin.Engine, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Error(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func postStr(t *testing.T, router *gin.Engine, path, data string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("POST", path, bytes.NewBuffer([]byte(data)))
	if err != nil {
		t.Error(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}
