package command

import (
	"github.com/JackBekket/telegram-gpt/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (c *Commander) AddNewUserToMap(updateMessage *tgbotapi.Message) {
	c.userDb[updateMessage.From.ID] = database.User{
		ID:            updateMessage.Chat.ID,
		Username:      updateMessage.Chat.UserName,
		Dialog_status: 0,
		Gpt_key:       "",
		Admin:         false,
	}

	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, msgTemplates["hello"])
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Start!")),
	)
	c.bot.Send(msg)

	// check for registration
	//	registred := IsAlreadyRegistred(session, updateMessage.From.ID)
	/*
		if registred {
			c.userDb[updateMessage.From.ID] = db.User{updateMessage.Chat.ID, updateMessage.Chat.UserName, 1}
		}
	*/
}
