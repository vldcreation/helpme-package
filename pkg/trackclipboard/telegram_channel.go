package trackclipboard

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	API_URL = "https://api.telegram.org/bot%s/sendMessage"
)

type TelegramChannel struct {
	Token  string
	ChatID string
}

func NewTelegramChannel(cfg *TelegramConfig) TrackChannel {
	return &TelegramChannel{
		Token:  cfg.Token,
		ChatID: cfg.ChatID,
	}
}

func (t *TelegramChannel) Send(ctx context.Context, msg string) error {
	query := url.Values{}
	query.Set("chat_id", t.ChatID)
	query.Set("text", msg)
	apiUrl := fmt.Sprintf(API_URL, t.Token)
	_, err := http.PostForm(apiUrl, query)
	return err
}

func (t *TelegramChannel) Close() error {
	return nil
}
