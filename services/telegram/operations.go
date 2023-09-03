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

const urlTelegram string = "https://api.telegram.org"

type TelegramInterface interface {
	SendMessage(message string, chatId string, isHtml bool)
	SendRepliedMessage(message string, reply string)
	SendTelegramMessage(message string, paramValues url.Values)
	SendTelegramCallbackQueryResponse(callbackQueryId string)
	SetWebhook(webhookUrl string, token string)
	SendPhotoGet(imgUrl string, chatId string) error
	SendPhoto(imgUrl string, chatId string) error
	GetCache() db.CacheRepositoryInterface
}

type Telegram struct {
	messageType MessageType
	url         string
	serviceUrl  string
	cache       db.CacheRepositoryInterface
	config      config.ConfigInterface
}

func newTextService(config config.ConfigInterface) *Telegram {
	return newTextServiceWithUrl(urlTelegram, config)
}

func newImageService(config config.ConfigInterface) *Telegram {
	return newImageServiceWithUrl(urlTelegram, config)
}

func newTextServiceWithUrl(url string, config config.ConfigInterface) *Telegram {
	t := &Telegram{
		messageType: Text,
		url:         url,
		config:      config,
	}
	t.init()
	return t
}

func newImageServiceWithUrl(url string, config config.ConfigInterface) *Telegram {
	t := &Telegram{
		messageType: Image,
		url:         url,
		config:      config,
	}
	t.init()
	return t
}

func (t *Telegram) init() {
	t.serviceUrl = t.getTelegramUrl()
}

func (t *Telegram) GetCache() db.CacheRepositoryInterface {
	if t.cache == nil {
		var err error
		if t.cache, err = db.NewCacheRepository(); err != nil {
			log.Fatalf("Error creating cache repository: %v", err)
		}
	}
	return t.cache
}

func (t *Telegram) getTelegramUrl() string {
	utils.IfThen(t.config == nil, "Config is nil", "")
	switch t.messageType {
	case Text:
		return utils.IfThen(
			t.config.GetTelegramBotTextToken() == "",
			t.url,
			fmt.Sprintf("%s/%s", t.url, t.config.GetTelegramBotTextToken()),
		)
	case Image:
		return utils.IfThen(
			t.config.GetTelegramBotImageToken() == "",
			t.url,
			fmt.Sprintf("%s/%s", t.url, t.config.GetTelegramBotImageToken()),
		)
	}
	return ""
}

func (t *Telegram) SendMessage(message string, chatId string, isHtml bool) (string, error) {
	params := url.Values{}
	params.Add("chat_id", chatId)
	if isHtml {
		params.Add("parse_mode", "html")
	}
	return t.SendTelegramMessage(message, params)
}

func (t *Telegram) SendRepliedMessage(message string, reply string) (string, error) {
	params := url.Values{}
	params.Add("reply_markup", reply)
	return t.SendTelegramMessage(message, params)
}

func (t *Telegram) SendTelegramMessage(message string, paramValues url.Values) (string, error) {
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

	return utils.IfThen(response != nil, response.Status, ""), err
}

func (t *Telegram) SendTelegramCallbackQueryResponse(callbackQueryId string) (string, error) {
	params := url.Values{}

	params.Add("callback_query_id", callbackQueryId)

	urlMsg := t.serviceUrl + "/answerCallbackQuery?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)

	return utils.IfThen(response != nil, response.Status, ""), err
}

func (t *Telegram) SetWebhook(webhookUrl string, token string) (string, error) {
	params := url.Values{}

	params.Add("url", webhookUrl)
	params.Add("secret_token", token)

	urlMsg := t.serviceUrl + "/setWebhook?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)

	return utils.IfThen(response != nil, response.Status, ""), err
}

func (t *Telegram) SendPhotoGet(imgUrl string, chatId string) (string, error) {
	params := url.Values{}
	params.Add("chat_id", chatId)
	params.Add("photo", imgUrl)

	urlMsg := t.serviceUrl + "/sendPhoto?" + params.Encode()

	fmt.Println("urlMsg", urlMsg)

	response, err := http.Get(urlMsg)

	fmt.Println("err", err)
	fmt.Println("response", response)
	return utils.IfThen(response != nil, response.Status, ""), err
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
	return nil
}
