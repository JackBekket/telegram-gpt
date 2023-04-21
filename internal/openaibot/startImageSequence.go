package openaibot

import (
	"context"
	"fmt"
	"log"

	"github.com/JackBekket/telegram-gpt/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartImageSequence(
	ID int64,
	promt string,
	ctx context.Context,
	bot *tgbotapi.BotAPI,
	updateMessage *tgbotapi.Message,
) {
	mu.Lock()
	defer mu.Unlock()
	userDatabase := database.UserMap
	sessionDatabase := database.AiSessionMap

	req := createImageRequest(promt)
	c := sessionDatabase[ID].Gpt_client

	resp, err := c.CreateImage(ctx, req)
	if err != nil {
		errorMessage(err, bot, ID, userDatabase)
	} else {

		respUrl := resp.Data[0].URL
		log.Printf("url image: %s\n", respUrl)

		msg1 := tgbotapi.NewMessage(userDatabase[ID].ID, "Done!")
		bot.Send(msg1)

		msg := tgbotapi.NewEditMessageText(
			userDatabase[updateMessage.From.ID].ID,
			updateMessage.MessageID+1,
			fmt.Sprintf("[Result](%s)", respUrl),
		)

		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)

		updatedatabase := userDatabase[ID]
		updatedatabase.Dialog_status = 4
		userDatabase[ID] = updatedatabase
	}
}
