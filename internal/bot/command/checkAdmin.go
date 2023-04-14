package command

import (
	"fmt"

	"github.com/JackBekket/telegram-gpt/internal/bot/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Updates "dialog_status" in the database. Admins - 2, other users - 0.
//
// Loads the key from env into the database.
func (c *Commander) CheckAdmin(id int64, updateMessage *tgbotapi.Message) {

	adminID := env.LoadAdminsID()
	adminAiKey := env.LoadAdminsAiKey()

	switch id {
	case adminID["ADMIN_ID"]:

		_, ok := adminAiKey["ADMIN_KEY"]
		if !ok {
			c.errorParseEnv("ADMIN_KEY", updateMessage)
			return
		}
		c.AddAdminToMap(adminAiKey["ADMIN_KEY"], updateMessage)
	case adminID["MINTY_ID"]:
		key, ok := adminAiKey["MINTY_KEY"]
		if !ok {
			c.errorParseEnv("MINTY_KEY", updateMessage)
			return
		}
		if key == "" {
			msg := tgbotapi.NewMessage(
				updateMessage.From.ID,
				"env \"%s\" is empty.",
			)
			c.bot.Send(msg)
			return
		}
		c.AddAdminToMap(adminAiKey["MINTY_KEY"], updateMessage)
	case adminID["OK_ID"]:
		key, ok := adminAiKey["OK_KEY"]
		if !ok {
			c.errorParseEnv("OK_KEY", updateMessage)
			return
		}
		if key == "" {
			msg := tgbotapi.NewMessage(
				updateMessage.From.ID,
				"env \"%s\" is empty.",
			)
			c.bot.Send(msg)
			return
		}
		c.AddAdminToMap(adminAiKey["OK_KEY"], updateMessage)
	case adminID["MURS_ID"]:
		key, ok := adminAiKey["MURS_KEY"]
		if !ok {
			c.errorParseEnv("MURS_KEY", updateMessage)
			return
		}
		if key == "" {
			msg := tgbotapi.NewMessage(
				updateMessage.From.ID,
				"env \"%s\" is empty.",
			)
			c.bot.Send(msg)
			return
		}
		c.AddAdminToMap(adminAiKey["MURS_KEY"], updateMessage)
	default:
		c.AddNewUserToMap(updateMessage)
	}

}

func (c *Commander) errorParseEnv(adminEnv string, updateMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(
		updateMessage.From.ID,
		fmt.Sprintf("env \"%s\" is missing.", adminEnv),
	)
	c.bot.Send(msg)

	// Directs to case 0
	c.AddNewUserToMap(updateMessage)
}
