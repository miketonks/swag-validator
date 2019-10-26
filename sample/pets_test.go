package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
	assert.Contains(t, resp.Body.String(), "Invalid type. Expected: integer, given: string")
}

func TestPostPet(t *testing.T) {

	api := server.SetupAPI()
	router := server.SetupRouter(api)

	t.Run("POST with a valid payload", func(t *testing.T) {
		validPet := `{"name": "ollie", "dob": "2018-01-01T12:00:00-09:00", "grumpy": true, "uuid": "1c694c09-3210-45d4-be6b-dbd94be1be4f"}`
		resp := postStr(t, router, "/pet", validPet)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("POST with a missing field", func(t *testing.T) {
		requiredFieldMissing := `{"dob": "2018-01-01T12:00:00-09:00", "grumpy": true, "uuid": "1c694c09-3210-45d4-be6b-dbd94be1be4f"}`
		resp := postStr(t, router, "/pet", requiredFieldMissing)
		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "name is required")
	})

	t.Run("POST with invalid date-time", func(t *testing.T) {
		invalidDateTime := `{"name": "ollie", "dob": "2018-01-01T12:00:0"}`
		resp := postStr(t, router, "/pet", invalidDateTime)
		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "Field does not match format 'date-time'")
	})

	t.Run("POST with invalid uuid", func(t *testing.T) {
		invalidUUID := `{"name": "ollie", "uuid": "foo-bar-baz"}`
		resp := postStr(t, router, "/pet", invalidUUID)
		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "Field does not match format 'uuid'")
	})

	t.Run("POST with invalid bool", func(t *testing.T) {
		invalidBool := `{"name": "ollie", "grumpy": "very"}`
		resp := postStr(t, router, "/pet", invalidBool)
		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "Invalid type. Expected: boolean, given: string")
	})

	t.Run("POST with valid date", func(t *testing.T) {
		validDate := `{"name": "ollie", "dt": "2018-04-01"}`
		resp := postStr(t, router, "/pet", validDate)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("POST with invalid date", func(t *testing.T) {
		invalidDate := `{"name": "ollie", "dt": "2018-44-01"}`
		resp := postStr(t, router, "/pet", invalidDate)
		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "month out of range")
	})

	t.Run("POST with valid time", func(t *testing.T) {
		validTime := `{"name": "ollie", "tm": "12:15:00"}`
		resp := postStr(t, router, "/pet", validTime)
		assert.Equal(t, 200, resp.Code)
	})

	t.Run("POST with invalid time", func(t *testing.T) {
		invalidTime := `{"name": "ollie", "tm": "12:15:99"}`
		resp := postStr(t, router, "/pet", invalidTime)
		assert.Equal(t, 400, resp.Code)
		assert.Contains(t, resp.Body.String(), "second out of range")
	})
}

func TestUploadFile(t *testing.T) {

	api := server.SetupAPI()
	router := server.SetupRouter(api)

	t.Run("POST with a valid payload", func(t *testing.T) {
		extraParams := map[string]string{"name": "test"}
		resp := postFile(t, router, "/upload", extraParams, "upfile", "pets_test.go")
		assert.Equal(t, 200, resp.Code)
	})
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

// Creates a new file upload http request with optional extra params
func postFile(t *testing.T, router *gin.Engine, path string, params map[string]string, paramName, filePath string) *httptest.ResponseRecorder {
	file, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
		return nil
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		t.Error(err)
		return nil
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Error(err)
		return nil
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		t.Error(err)
		return nil
	}

	req, err := http.NewRequest("POST", path, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		t.Error(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}
