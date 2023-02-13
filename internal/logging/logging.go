package logging

import (
	"os"

	as "github.com/palavrapasse/aspirador/pkg"
)

const (
	telegramBotTokenEnvKey = "telegram_token"
	telegramChatIdEnvKey   = "telegram_chat_id"
)

const loggingFilePathEnvKey = "logging_fp"

var (
	telegramBotToken = os.Getenv(telegramBotTokenEnvKey)
	telegramChatId   = os.Getenv(telegramChatIdEnvKey)
)

var loggingFilePath = os.Getenv(loggingFilePathEnvKey)

var Aspirador as.Aspirador

func CreateAspiradorClients() []as.Client {

	consoleClient := as.NewConsoleClient()

	telegramClient := as.NewTelegramClient(telegramBotToken, telegramChatId, as.WARNING, as.ERROR)

	fileClient, err := as.NewFileClient(loggingFilePath)

	if err != nil {
		return []as.Client{consoleClient, telegramClient}
	}

	return []as.Client{consoleClient, telegramClient, fileClient}
}
