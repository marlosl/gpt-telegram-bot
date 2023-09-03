package chatgpt

import (
	// "encoding/json"

	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	// "net/url"

	"testing"

	"github.com/marlosl/gpt-telegram-bot/utils/tests"
	"github.com/stretchr/testify/assert"
)

const contentTypeKey = "Content-Type"

func createPostServer(t *testing.T) *httptest.Server {
	ts := tests.CreateTestServer(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)
		t.Logf("RawQuery: %v", r.URL.RawQuery)
		t.Logf("Content-Type: %v", r.Header.Get(contentTypeKey))

		if r.Method == "POST" {
			switch r.URL.Path {
			case "/login-json-html":
				w.Header().Set(contentTypeKey, "text/html")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`<htm><body>Test JSON request with HTML response</body></html>`))
				return
			case "/usersmap":
				// JSON
				if strings.Contains(r.Header.Get(contentTypeKey), "application/json") {
					if r.URL.Query().Get("status") == "500" {
						body, err := io.ReadAll(r.Body)
						if err != nil {
							t.Errorf("Error: could not read post body: %s", err.Error())
						}
						t.Logf("Got query param: status=500 so we're returning the post body as response and a 500 status code. body: %s", string(body))
						w.Header().Set(contentTypeKey, "application/json; charset=utf-8")
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = w.Write(body)
						return
					}

					var users []map[string]interface{}
					jd := json.NewDecoder(r.Body)
					err := jd.Decode(&users)
					w.Header().Set(contentTypeKey, "application/json; charset=utf-8")
					if err != nil {
						t.Logf("Error: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write([]byte(`{ "id": "bad_request", "message": "Unable to read user info" }`))
						return
					}

					// logic check, since we are excepting to reach 1 map records
					if len(users) != 1 {
						t.Log("Error: Excepted count of 1 map records")
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write([]byte(`{ "id": "bad_request", "message": "Expected record count doesn't match" }`))
						return
					}

					w.WriteHeader(http.StatusAccepted)
					_, _ = w.Write([]byte(`{ "message": "Accepted" }`))

					return
				}
			case "/redirect":
				w.Header().Set("Location", "/login")
				w.WriteHeader(http.StatusTemporaryRedirect)
			case "/redirect-with-body":
				body, _ := io.ReadAll(r.Body)
				query := url.Values{}
				query.Add("body", string(body))
				w.Header().Set("Location", "/redirected-with-body?"+query.Encode())
				w.WriteHeader(http.StatusTemporaryRedirect)
			case "/redirected-with-body":
				// body, _ := io.ReadAll(r.Body)
				// assertEqual(t, r.URL.Query().Get("body"), string(body))
				w.WriteHeader(http.StatusOK)
			}
		}
	})

	return ts
}

func TestTelegramSendMessage(t *testing.T) {
	ts := createPostServer(t)
	defer ts.Close()

	assert.NotNil(t, ts)
}
