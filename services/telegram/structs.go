package telegram

import (
	"context"
	"time"
)

type RecordType interface {
	Package | Event | ChatStatus
}

type Package struct {
	ID          string `json:"id,omitempty" bson:"_id,omitempty"`
	Package     string `json:"package"`
	Description string `json:"description"`
	Delivered   bool   `json:"delivered"`
}

type Event struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	Package    string    `json:"package"`
	LastUpdate time.Time `json:"lastUpdate"`
	Checksum   string    `json:"checksum"`
}

type ChatStatus struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId        int64     `json:"chatId" bson:"chatId"`
	Status        string    `json:"status"`
	LastParameter string    `json:"lastParameter" bson:"lastParameter"`
	LastUpdate    time.Time `json:"lastUpdate" bson:"lastUpdate"`
}

type From struct {
	ID           int64   `json:"id,omitempty"`
	IsBot        bool    `json:"is_bot"`
	FirstName    string  `json:"first_name"`
	LanguageCode *string `json:"language_code"`
	UserName     *string `json:"username"`
}

type Chat struct {
	ID        int64  `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
}

type Entity struct {
	OffSet int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type InlineKeyboard struct {
	Buttons [1][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type Message struct {
	MessageId    int64           `json:"message_id,omitempty"`
	From         *From           `json:"from,omitempty"`
	Chat         *Chat           `json:"chat,omitempty"`
	Date         int64           `json:"date,omitempty"`
	Text         string          `json:"text"`
	Entities     *[]Entity       `json:"entities,omitempty"`
	ReplayMarkup *InlineKeyboard `json:"reply_markup,omitempty"`
}

type CallbackQuery struct {
	ID           string   `json:"id"`
	From         *From    `json:"from,omitempty"`
	Message      *Message `json:"message,omitempty"`
	ChatInstance string   `json:"chat_instance"`
	Data         string   `json:"data"`
}

type WebhookMessage struct {
	UpdateId      int64          `json:"update_id,omitempty"`
	Message       *Message       `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type MsgParams struct {
	FirstArg        string
	SecondArg       string
	ThirdArg        string
	LastStatus      string
	PkgNumber       string
	PkgDescription  string
	ChatId          int64
	ChatMsg         string
	Ctx             context.Context
	Status          *ChatStatus
	CallbackQueryId string
}

type ImageMessage struct {
	ChatId   string `json:"chatId"`
	ImageUrl string `json:"imageUrl"`
}
