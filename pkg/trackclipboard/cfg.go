package trackclipboard

import "time"

type Config struct {
	Channel  string
	Idle     time.Duration
	File     *FileConfig
	Telegram *TelegramConfig
}

type FileConfig struct {
	Path string
	Name string
}

type TelegramConfig struct {
	Token  string
	ChatID string
}
