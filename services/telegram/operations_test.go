package telegram

import (
	// "encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	// "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/marlosl/gpt-telegram-bot/utils/config"
	"github.com/marlosl/gpt-telegram-bot/utils/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodPatch   = "PATCH"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

func getTestDataPath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, ".testdata")
}

func createGetServer(t *testing.T) *httptest.Server {
	var attempt int32
	var sequence int32
	var lastRequest time.Time
	ts := tests.CreateTestServer(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Method: %v", r.Method)
		t.Logf("Path: %v", r.URL.Path)

		if r.Method == MethodGet {
			switch r.URL.Path {
			case "/":
				_, _ = w.Write([]byte("TestGet: text response"))
			case "/sendMessage":
				_, _ = w.Write([]byte(""))
			case "/answerCallbackQuery":
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"TestGet": "JSON response"}`))
			case "/setWebhook":
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte("TestGet: Invalid JSON"))
			case "/sendPhoto":
				_, _ = w.Write([]byte("TestGet: text response with size > 30"))
			case "/long-json":
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"TestGet": "JSON response with size > 30"}`))
			case "/mypage":
				w.WriteHeader(http.StatusBadRequest)
			case "/mypage2":
				_, _ = w.Write([]byte("TestGet: text response from mypage2"))
			case "/set-retrycount-test":
				attp := atomic.AddInt32(&attempt, 1)
				if attp <= 4 {
					time.Sleep(time.Second * 6)
				}
				_, _ = w.Write([]byte("TestClientRetry page"))
			case "/set-retrywaittime-test":
				// Returns time.Duration since last request here
				// or 0 for the very first request
				if atomic.LoadInt32(&attempt) == 0 {
					lastRequest = time.Now()
					_, _ = fmt.Fprint(w, "0")
				} else {
					now := time.Now()
					sinceLastRequest := now.Sub(lastRequest)
					lastRequest = now
					_, _ = fmt.Fprintf(w, "%d", uint64(sinceLastRequest))
				}
				atomic.AddInt32(&attempt, 1)

			case "/set-timeout-test-with-sequence":
				seq := atomic.AddInt32(&sequence, 1)
				time.Sleep(time.Second * 2)
				_, _ = fmt.Fprintf(w, "%d", seq)
			case "/set-timeout-test":
				time.Sleep(time.Second * 6)
				_, _ = w.Write([]byte("TestClientTimeout page"))
			case "/my-image.png":
				fileBytes, _ := os.ReadFile(filepath.Join(getTestDataPath(), "test-img.png"))
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
				_, _ = w.Write(fileBytes)
			case "/get-method-payload-test":
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("Error: could not read get body: %s", err.Error())
				}
				_, _ = w.Write(body)
			case "/host-header":
				_, _ = w.Write([]byte(r.Host))
			}

			switch {
			case strings.HasPrefix(r.URL.Path, "/v1/users/sample@sample.com/100002"):
				if strings.HasSuffix(r.URL.Path, "details") {
					_, _ = w.Write([]byte("TestGetPathParams: text response: " + r.URL.String()))
				} else {
					_, _ = w.Write([]byte("TestPathParamURLInput: text response: " + r.URL.String()))
				}
			}

		}
	})

	return ts
}

func createTelegramTextService() *Telegram {
	cfgMock := new(config.ConfigMock)
	cfgMock.On("GetTelegramBotTextToken", mock.Anything).Return("")

	telegram := NewTextService(cfgMock)
	return telegram
}

func TestTelegramSendMessage(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	assert.NotNil(t, ts)

	telegram := createTelegramTextService()
	assert.NotNil(t, telegram)

	assert.NotEmpty(t, telegram.getTelegramUrl())

	status, err := telegram.SendMessage("TestSendMessage", "1234567890", false)
	assert.Nil(t, err)
	assert.Equal(t, "200 OK", status)
}

func TestTelegramSendRepliedMessage(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	assert.NotNil(t, ts)

	telegram := createTelegramTextService()
	assert.NotNil(t, telegram)

	assert.NotEmpty(t, telegram.getTelegramUrl())

	status, err := telegram.SendRepliedMessage("TestSendMessage", "1234567890")
	assert.Nil(t, err)
	assert.Equal(t, "200 OK", status)
}

func TestTelegramSendTelegramCallbackQueryResponse(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	assert.NotNil(t, ts)

	telegram := createTelegramTextService()
	assert.NotNil(t, telegram)

	assert.NotEmpty(t, telegram.getTelegramUrl())

	status, err := telegram.SendTelegramCallbackQueryResponse("TestSendMessage")
	assert.Nil(t, err)
	assert.Equal(t, "200 OK", status)
}

func TestTelegramSetWebhook(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	assert.NotNil(t, ts)

	telegram := createTelegramTextService()
	assert.NotNil(t, telegram)

	assert.NotEmpty(t, telegram.getTelegramUrl())

	status, err := telegram.SetWebhook("https://webhook.com/status", "1234567890")
	assert.Nil(t, err)
	assert.Equal(t, "200 OK", status)
}

func TestTelegramSendPhotoGet(t *testing.T) {
	ts := createGetServer(t)
	defer ts.Close()

	assert.NotNil(t, ts)

	telegram := createTelegramTextService()
	assert.NotNil(t, telegram)

	assert.NotEmpty(t, telegram.getTelegramUrl())

	status, err := telegram.SendPhotoGet("http://image/image.png", "1234567890")
	assert.Nil(t, err)
	assert.Equal(t, "200 OK", status)
}

// func createPostServer(t *testing.T) *httptest.Server {
// 	ts := tests.CreateTestServer(func(w http.ResponseWriter, r *http.Request) {
// 		t.Logf("Method: %v", r.Method)
// 		t.Logf("Path: %v", r.URL.Path)
// 		t.Logf("RawQuery: %v", r.URL.RawQuery)
// 		t.Logf("Content-Type: %v", r.Header.Get(hdrContentTypeKey))

// 		if r.Method == MethodPost {
// 			handleLoginEndpoint(t, w, r)

// 			handleUsersEndpoint(t, w, r)
// 			switch r.URL.Path {
// 			case "/login-json-html":
// 				w.Header().Set(hdrContentTypeKey, "text/html")
// 				w.WriteHeader(http.StatusOK)
// 				_, _ = w.Write([]byte(`<htm><body>Test JSON request with HTML response</body></html>`))
// 				return
// 			case "/usersmap":
// 				// JSON
// 				if IsJSONType(r.Header.Get(hdrContentTypeKey)) {
// 					if r.URL.Query().Get("status") == "500" {
// 						body, err := io.ReadAll(r.Body)
// 						if err != nil {
// 							t.Errorf("Error: could not read post body: %s", err.Error())
// 						}
// 						t.Logf("Got query param: status=500 so we're returning the post body as response and a 500 status code. body: %s", string(body))
// 						w.Header().Set(hdrContentTypeKey, "application/json; charset=utf-8")
// 						w.WriteHeader(http.StatusInternalServerError)
// 						_, _ = w.Write(body)
// 						return
// 					}

// 					var users []map[string]interface{}
// 					jd := json.NewDecoder(r.Body)
// 					err := jd.Decode(&users)
// 					w.Header().Set(hdrContentTypeKey, "application/json; charset=utf-8")
// 					if err != nil {
// 						t.Logf("Error: %v", err)
// 						w.WriteHeader(http.StatusBadRequest)
// 						_, _ = w.Write([]byte(`{ "id": "bad_request", "message": "Unable to read user info" }`))
// 						return
// 					}

// 					// logic check, since we are excepting to reach 1 map records
// 					if len(users) != 1 {
// 						t.Log("Error: Excepted count of 1 map records")
// 						w.WriteHeader(http.StatusBadRequest)
// 						_, _ = w.Write([]byte(`{ "id": "bad_request", "message": "Expected record count doesn't match" }`))
// 						return
// 					}

// 					w.WriteHeader(http.StatusAccepted)
// 					_, _ = w.Write([]byte(`{ "message": "Accepted" }`))

// 					return
// 				}
// 			case "/redirect":
// 				w.Header().Set(hdrLocationKey, "/login")
// 				w.WriteHeader(http.StatusTemporaryRedirect)
// 			case "/redirect-with-body":
// 				body, _ := io.ReadAll(r.Body)
// 				query := url.Values{}
// 				query.Add("body", string(body))
// 				w.Header().Set(hdrLocationKey, "/redirected-with-body?"+query.Encode())
// 				w.WriteHeader(http.StatusTemporaryRedirect)
// 			case "/redirected-with-body":
// 				body, _ := io.ReadAll(r.Body)
// 				assertEqual(t, r.URL.Query().Get("body"), string(body))
// 				w.WriteHeader(http.StatusOK)
// 			}
// 		}
// 	})

// 	return ts
// }

// func createFormPostServer(t *testing.T) *httptest.Server {
// 	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
// 		t.Logf("Method: %v", r.Method)
// 		t.Logf("Path: %v", r.URL.Path)
// 		t.Logf("Content-Type: %v", r.Header.Get(hdrContentTypeKey))

// 		if r.Method == MethodPost {
// 			_ = r.ParseMultipartForm(10e6)

// 			if r.URL.Path == "/profile" {
// 				t.Logf("FirstName: %v", r.FormValue("first_name"))
// 				t.Logf("LastName: %v", r.FormValue("last_name"))
// 				t.Logf("City: %v", r.FormValue("city"))
// 				t.Logf("Zip Code: %v", r.FormValue("zip_code"))

// 				_, _ = w.Write([]byte("Success"))
// 				return
// 			} else if r.URL.Path == "/search" {
// 				formEncodedData := r.Form.Encode()
// 				t.Logf("Received Form Encoded values: %v", formEncodedData)

// 				assertEqual(t, true, strings.Contains(formEncodedData, "search_criteria=pencil"))
// 				assertEqual(t, true, strings.Contains(formEncodedData, "search_criteria=glass"))

// 				_, _ = w.Write([]byte("Success"))
// 				return
// 			} else if r.URL.Path == "/upload" {
// 				t.Logf("FirstName: %v", r.FormValue("first_name"))
// 				t.Logf("LastName: %v", r.FormValue("last_name"))

// 				targetPath := filepath.Join(getTestDataPath(), "upload")
// 				_ = os.MkdirAll(targetPath, 0700)

// 				for _, fhdrs := range r.MultipartForm.File {
// 					for _, hdr := range fhdrs {
// 						t.Logf("Name: %v", hdr.Filename)
// 						t.Logf("Header: %v", hdr.Header)
// 						dotPos := strings.LastIndex(hdr.Filename, ".")

// 						fname := fmt.Sprintf("%s-%v%s", hdr.Filename[:dotPos], time.Now().Unix(), hdr.Filename[dotPos:])
// 						t.Logf("Write name: %v", fname)

// 						infile, _ := hdr.Open()
// 						f, err := os.OpenFile(filepath.Join(targetPath, fname), os.O_WRONLY|os.O_CREATE, 0666)
// 						if err != nil {
// 							t.Logf("Error: %v", err)
// 							return
// 						}
// 						defer func() {
// 							_ = f.Close()
// 						}()
// 						_, _ = io.Copy(f, infile)

// 						_, _ = w.Write([]byte(fmt.Sprintf("File: %v, uploaded as: %v\n", hdr.Filename, fname)))
// 					}
// 				}

// 				return
// 			}
// 		}
// 	})

// 	return ts
// }

// func createFormPatchServer(t *testing.T) *httptest.Server {
// 	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
// 		t.Logf("Method: %v", r.Method)
// 		t.Logf("Path: %v", r.URL.Path)
// 		t.Logf("Content-Type: %v", r.Header.Get(hdrContentTypeKey))

// 		if r.Method == MethodPatch {
// 			_ = r.ParseMultipartForm(10e6)

// 			if r.URL.Path == "/upload" {
// 				t.Logf("FirstName: %v", r.FormValue("first_name"))
// 				t.Logf("LastName: %v", r.FormValue("last_name"))

// 				targetPath := filepath.Join(getTestDataPath(), "upload")
// 				_ = os.MkdirAll(targetPath, 0700)

// 				for _, fhdrs := range r.MultipartForm.File {
// 					for _, hdr := range fhdrs {
// 						t.Logf("Name: %v", hdr.Filename)
// 						t.Logf("Header: %v", hdr.Header)
// 						dotPos := strings.LastIndex(hdr.Filename, ".")

// 						fname := fmt.Sprintf("%s-%v%s", hdr.Filename[:dotPos], time.Now().Unix(), hdr.Filename[dotPos:])
// 						t.Logf("Write name: %v", fname)

// 						infile, _ := hdr.Open()
// 						f, err := os.OpenFile(filepath.Join(targetPath, fname), os.O_WRONLY|os.O_CREATE, 0666)
// 						if err != nil {
// 							t.Logf("Error: %v", err)
// 							return
// 						}
// 						defer func() {
// 							_ = f.Close()
// 						}()
// 						_, _ = io.Copy(f, infile)

// 						_, _ = w.Write([]byte(fmt.Sprintf("File: %v, uploaded as: %v\n", hdr.Filename, fname)))
// 					}
// 				}

// 				return
// 			}
// 		}
// 	})

// 	return ts
// }

// func createFilePostServer(t *testing.T) *httptest.Server {
// 	ts := createTestServer(func(w http.ResponseWriter, r *http.Request) {
// 		t.Logf("Method: %v", r.Method)
// 		t.Logf("Path: %v", r.URL.Path)
// 		t.Logf("Content-Type: %v", r.Header.Get(hdrContentTypeKey))

// 		if r.Method != MethodPost {
// 			t.Log("createPostServer:: Not a Post request")
// 			w.WriteHeader(http.StatusBadRequest)
// 			fmt.Fprint(w, http.StatusText(http.StatusBadRequest))
// 			return
// 		}

// 		targetPath := filepath.Join(getTestDataPath(), "upload-large")
// 		_ = os.MkdirAll(targetPath, 0700)
// 		defer cleanupFiles(targetPath)

// 		switch r.URL.Path {
// 		case "/upload":
// 			f, err := os.OpenFile(filepath.Join(targetPath, "large-file.png"),
// 				os.O_WRONLY|os.O_CREATE, 0666)
// 			if err != nil {
// 				t.Logf("Error: %v", err)
// 				return
// 			}
// 			defer func() {
// 				_ = f.Close()
// 			}()
// 			size, _ := io.Copy(f, r.Body)

// 			fmt.Fprintf(w, "File Uploaded successfully, file size: %v", size)
// 		case "/set-reset-multipart-readers-test":
// 			w.Header().Set(hdrContentTypeKey, "application/json; charset=utf-8")
// 			w.WriteHeader(http.StatusInternalServerError)
// 			_, _ = fmt.Fprintf(w, `{ "message": "error" }`)
// 		}
// 	})

// 	return ts
// }
