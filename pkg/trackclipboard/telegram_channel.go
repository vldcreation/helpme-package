package trackclipboard

import "context"

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
	return nil
}

func (t *TelegramChannel) Close() error {
	return nil
}
