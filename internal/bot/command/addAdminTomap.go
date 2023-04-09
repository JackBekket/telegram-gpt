package command

import (
	"fmt"

	"github.com/JackBekket/telegram-gpt/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commander) AddAdminToMap(
	adminKey string,
	updateMessage *tgbotapi.Message,

) bool {
	c.userDb[updateMessage.From.ID] = database.User{
		ID:            updateMessage.Chat.ID,
		Username:      updateMessage.Chat.UserName,
		Dialog_status: 2,
		Gpt_key:       adminKey,
		Admin:         true,
	}

	fmt.Printf("%s authorized\n", updateMessage.Chat.UserName)

	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "authorized")
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, msgTemplates["case1"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("GPT-3.5"),
			//tgbotapi.NewKeyboardButton("GPT-4"),
			tgbotapi.NewKeyboardButton("Codex")),
	)
	c.bot.Send(msg)
	return true
}
