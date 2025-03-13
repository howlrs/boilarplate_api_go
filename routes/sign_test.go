package routes_test

import (
	"backend/routes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// サインインとJWTトークンを使ったプライベートリクエストのテスト
func TestSinginAndPrivateRequest(t *testing.T) {
	godotenv.Load("../.env.local")

	params := []map[string]string{
		{
			"name":         "test_signin",
			"url":          "/api/v1/public/signin",
			"method":       "POST",
			"body":         `{"email": "test_signin", "password": "test_signin"}`,
			"expactStatus": "200",
		},
		{
			"name":         "test_private_health",
			"url":          "/api/v1/private/health",
			"method":       "GET",
			"body":         "",
			"expactStatus": "200",
		},
	}

	// echo モックアップサーバ
	e := echo.New()
	isTest := true
	routes.Endpoint(e, isTest)

	var jwtToken string

	for _, param := range params {
		// create request
		var req *http.Request
		if jwtToken != "" {
			req = httptest.NewRequest(param["method"], param["url"], strings.NewReader(param["body"]))
			req.Header.Set("Authorization", "Bearer "+jwtToken)
		} else if param["body"] != "" {
			req = httptest.NewRequest(param["method"], param["url"], strings.NewReader(param["body"]))
		} else {
			req = httptest.NewRequest(param["method"], param["url"], nil)
		}

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var result map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&result)
		assert.NoError(t, err, param["name"])

		data, isThere := result["data"].(map[string]interface{})
		if isThere {
			v, ok := data["token"]
			if ok {
				jwtToken = v.(string)
			}
		}

		// status code
		if param["expactStatus"] != "" {
			assert.Equal(t, param["expactStatus"], fmt.Sprintf("%v", rec.Code), param["name"])
		}
	}
}
