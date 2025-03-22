package trackclipboard

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	apiUrl := fmt.Sprintf(API_URL, t.Token)
	data := fmt.Sprintf("chat_id=%s&text=%s", t.ChatID, msg)
	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}

func (t *TelegramChannel) Close() error {
	return nil
}
