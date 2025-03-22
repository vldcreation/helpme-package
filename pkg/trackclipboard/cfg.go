package trackclipboard

import "time"

type Config struct {
	App      *APPConfig      `yaml:"app" env:"APP" mapstructure:"app"`
	File     *FileConfig     `yaml:"file" env:"FILE" mapstructure:"file"`
	Telegram *TelegramConfig `yaml:"telegram" env:"TELEGRAM" mapstructure:"telegram"`
}

type APPConfig struct {
	Channel string        `yaml:"channel" env:"APP_CHANNEL" mapstructure:"app_channel"`
	Idle    time.Duration `yaml:"idle" env:"APP_IDLE" mapstructure:"app_idle"`
}
type FileConfig struct {
	Path string `yaml:"path" env:"FILE_PATH" mapstructure:"file_path"`
	Name string `yaml:"name" env:"FILE_NAME" mapstructure:"file_name"`
}

type TelegramConfig struct {
	Token  string `yaml:"token" env:"TOKEN" mapstructure:"telegram_token"`
	ChatID string `yaml:"chat_id" env:"CHAT_ID" mapstructure:"telegram_chat_id"`
}
