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

func TestRequestToReservation(t *testing.T) {
	godotenv.Load("../.env.local")
	params := []map[string]string{
		{
			"name":   "create",
			"to":     "/api/v1/reservation",
			"method": "POST",
			"body": `{
				"start_time": "2021-01-01T00:00:00Z",
				"end_time": "2021-01-01T01:00:00Z",
				"timezone": "Asia/Tokyo",
				"content": "test",
				"sub_content": "test"
			}`,
			"expact": "success, create reservation",
		},
		{
			"name":   "read",
			"to":     "/api/v1/reservation",
			"method": "GET",
			"query":  "id=target_id",
			"expact": "success, read reservation",
		},
		{
			"name":   "update",
			"to":     "/api/v1/reservation",
			"method": "PUT",
			"body": `{
				"id": "target_id",
				"content": "test2"
			}`,
			"expact": "success, update reservation",
		},
		{
			"name":   "delete",
			"to":     "/api/v1/reservation",
			"method": "DELETE",
			"query":  "id=target_id",
			"expact": "success, canceled reservation",
		},
		{
			"name":   "re-read",
			"to":     "/api/v1/reservation",
			"method": "GET",
			"query":  "id=target_id",
			"expact": "-1", //削除フラグ
		},
	}

	// mock server
	e := echo.New()
	routes.Endpoint(e, true)

	var id string

	for _, p := range params {
		t.Run(p["method"]+" "+p["to"], func(t *testing.T) {
			var req *http.Request

			// target_idをidに置き換える
			if id != "" {
				p["body"] = strings.ReplaceAll(p["body"], "target_id", id)
				p["query"] = strings.ReplaceAll(p["query"], "target_id", id)
			}

			if p["body"] != "" {
				req = httptest.NewRequest(p["method"], p["to"]+"?"+p["query"], strings.NewReader(p["body"]))
			} else {
				req = httptest.NewRequest(p["method"], p["to"]+"?"+p["query"], nil)
			}

			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			var m map[string]interface{}
			err := json.NewDecoder(rec.Body).Decode(&m)
			assert.NoError(t, err, p["name"])

			data, isThere := m["data"].(map[string]interface{})
			if isThere {
				v, ok := data["id"]
				if ok {
					id = v.(string)
				}
			}

			if p["name"] == "re-read" {
				assert.Equal(t, p["expact"], fmt.Sprintf("%v", data["status"]), p["name"])
			} else {
				assert.Equal(t, p["expact"], m["message"], p["name"])
			}

		})
	}
}
