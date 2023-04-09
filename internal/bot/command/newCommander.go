package command

import (
	"github.com/JackBekket/telegram-gpt/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Commander struct {
	bot         *tgbotapi.BotAPI
	userDb      map[int64]database.User
	aiSessionDb map[int64]database.AiSession
}

func NewCommander(
	bot *tgbotapi.BotAPI,
	userDb map[int64]database.User,
	sessionDB map[int64]database.AiSession,
) *Commander {
	return &Commander{
		bot:         bot,
		userDb:      userDb,
		aiSessionDb: sessionDB,
	}
}
