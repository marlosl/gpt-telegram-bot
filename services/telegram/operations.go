package telegram

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/marlosl/gpt-telegram-bot/clients/db"
	"github.com/marlosl/gpt-telegram-bot/utils"
	"github.com/marlosl/gpt-telegram-bot/utils/config"
)

type MessageType byte

const (
	Text  MessageType = 1
	Image MessageType = 2
)

type Telegram struct {
	Type       MessageType
	serviceUrl string
	Cache      *db.CacheRepository
}

const urlTelegram string = "https://api.telegram.org"

func NewTextService() *Telegram {
	t := &Telegram{
		Type: Text,
	}
	t.Init()
	return t
}

func NewImageService() *Telegram {
	t := &Telegram{
		Type: Image,
	}
	t.Init()
	return t
}

func (t *Telegram) Init() {
	var err error
	t.serviceUrl = t.GetTelegramUrl()
	t.Cache, err = db.NewCacheRepository()
	if err != nil {
		log.Fatalf("Error creating cache repository: %v", err)
	}
}

func (t *Telegram) GetTelegramUrl() string {
	switch t.Type {
	case Text:
		return fmt.Sprintf("%s/%s", urlTelegram, config.Store.TelegramBotTextToken)
	case Image:
		return fmt.Sprintf("%s/%s", urlTelegram, config.Store.TelegramBotImageToken)
	}
	return ""
}

func (t *Telegram) SendMessage(message string, chatId string, isHtml bool) {
	params := url.Values{}
	params.Add("chat_id", chatId)
	if isHtml {
		params.Add("parse_mode", "html")
	}
	t.SendTelegramMessage(message, params)
}

func (t *Telegram) SendRepliedMessage(message string, reply string) {
	params := url.Values{}
	params.Add("reply_markup", reply)
	t.SendTelegramMessage(message, params)
}

func (t *Telegram) SendTelegramMessage(message string, paramValues url.Values) {
	params := url.Values{}
	if (len(paramValues)) > 0 {
		params = paramValues
	}

	params.Add("text", message)

	urlMsg := t.serviceUrl + "/sendMessage?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)
}

func (t *Telegram) SendTelegramCallbackQueryResponse(callbackQueryId string) {
	params := url.Values{}

	params.Add("callback_query_id", callbackQueryId)

	urlMsg := t.serviceUrl + "/answerCallbackQuery?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)
}

func (t *Telegram) SetWebhook(webhookUrl string, token string) {
	params := url.Values{}

	params.Add("url", webhookUrl)
	params.Add("secret_token", token)

	urlMsg := t.serviceUrl + "/setWebhook?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)
}

func (t *Telegram) SendPhotoGet(imgUrl string, chatId string) error {
	params := url.Values{}
	params.Add("chat_id", chatId)
	params.Add("photo", imgUrl)

	urlMsg := t.serviceUrl + "/sendPhoto?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)
	return err
}

func (t *Telegram) SendPhoto(imgUrl string, chatId string) error {
	params := url.Values{}
	params.Add("chat_id", chatId)
	urlMsg := t.serviceUrl + "/sendPhoto?" + params.Encode()

	fmt.Println("sendPhotoUrl", urlMsg)

	imgFile, err := http.Get(imgUrl)
	if err != nil {
		fmt.Printf("Error getting image: %v\n", err)
		return err
	}
	defer imgFile.Body.Close()
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("photo", imgUrl)
	if err != nil {
		fmt.Printf("Error creating form field: %v\n", err)
		return err
	}

	_, err = io.Copy(fw, imgFile.Body)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return err
	}

	// Close multipart writer.
	writer.Close()
	req, err := http.NewRequest("POST", urlMsg, bytes.NewReader(body.Bytes()))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending image: %v\n", err)
	}

	fmt.Printf("SendPhoto Response: %s\n", utils.SPrintJson(rsp))
	if rsp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with response code: %d\n", rsp.StatusCode)
	}
	fmt.Println("SendPhoto - 9")
	return nil
}
