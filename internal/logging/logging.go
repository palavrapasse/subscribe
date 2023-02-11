package logging

import (
	"os"

	as "github.com/palavrapasse/aspirador/pkg"
)

const (
	telegramBotTokenEnvKey = "telegram_token"
	telegramChatIdEnvKey   = "telegram_chat_id"
)

var (
	telegramBotToken = os.Getenv(telegramBotTokenEnvKey)
	telegramChatId   = os.Getenv(telegramChatIdEnvKey)
)

var Aspirador as.Aspirador

func CreateAspiradorClients() []as.Client {

	consoleClient := as.NewConsoleClient()

	telegramClient := as.NewTelegramClient(telegramBotToken, telegramChatId, as.WARNING, as.ERROR)

	return []as.Client{consoleClient, telegramClient}
}
