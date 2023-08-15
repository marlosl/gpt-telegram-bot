package telegram

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTelegramUrl(t *testing.T) {
	tests := []struct {
		serviceType MessageType
		expectedURL string
	}{
		{
			serviceType: Text,
			expectedURL: "https://api.telegram.org/your_telegram_bot_text_token_here",
		},
		{
			serviceType: Image,
			expectedURL: "https://api.telegram.org/your_telegram_bot_image_token_here",
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.serviceType), func(t *testing.T) {
			telegram := &Telegram{Type: tt.serviceType}
			assert.Equal(t, tt.expectedURL, telegram.GetTelegramUrl())
		})
	}
}

func TestSendTelegramMessage(t *testing.T) {
	// Mock http server for Telegram API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if endpoint is correct
		if strings.HasSuffix(r.URL.Path, "/sendMessage") {
			// Check parameters
			assert.Equal(t, "1234567890", r.URL.Query().Get("chat_id"))
			assert.Equal(t, "Hello", r.URL.Query().Get("text"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer ts.Close()

	// Override telegram url for testing
	urlTelegram = ts.URL

	tel := NewTextService()
	tel.SendMessage("Hello", "1234567890", false)
}

func TestSendPhoto(t *testing.T) {
	// Mock http server for Telegram API and image fetch
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/image.jpg" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fake-image-data"))
			return
		}

		// Check if endpoint is correct
		if strings.HasSuffix(r.URL.Path, "/sendPhoto") {
			assert.Equal(t, "1234567890", r.URL.Query().Get("chat_id"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	// Override telegram url and imgURL for testing
	urlTelegram = ts.URL
	imgURL := ts.URL + "/image.jpg"

	tel := NewImageService()
	err := tel.SendPhoto(imgURL, "1234567890")

	assert.Nil(t, err)
}

// Testing error scenario
func TestSendPhoto_ErrorScenario(t *testing.T) {
	// Mock http server for Telegram API and image fetch to induce an error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer ts.Close()

	// Override telegram url and imgURL for testing
	urlTelegram = ts.URL
	imgURL := ts.URL + "/image.jpg"

	tel := NewImageService()
	err := tel.SendPhoto(imgURL, "1234567890")

	assert.NotNil(t, err)
}
