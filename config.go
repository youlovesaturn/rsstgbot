package main

type config struct {
	TelegramToken  string `env:"TOKEN"`
	ChannelID      int64  `env:"CHANNEL"`
	FeedURL        string `env:"FEEDURL"`
	DisablePreview bool   `env:"DISABLE_PREVIEW"`
	Filename       string `env:"FILENAME"`
}
