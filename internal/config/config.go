package config

import (
	"fmt"
	"os"
)

type Config struct {
	TelegramToken string
	OllamaURL     string
	DBPath        string
}

func MustLoad() Config {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		panic(fmt.Sprintf("TELEGRAM_TOKEN is required"))
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "kypidbot.db"
	}

	return Config{
		TelegramToken: token,
		OllamaURL:     ollamaURL,
		DBPath:        dbPath,
	}
}
