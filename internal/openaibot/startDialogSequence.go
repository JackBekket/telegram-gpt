package openaibot

import (
	"context"
	"fmt"
	"log"

	"github.com/JackBekket/telegram-gpt/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Notifies the user that an error occurred while creating the request.
// "An error has occured. In order to proceed we need to recreate client and initialize new session"
func errorMessage(err error, bot *tgbotapi.BotAPI, ID int64, db map[int64]database.User) {

	log.Println("error :", err)
	msg := tgbotapi.NewMessage(db[ID].ID, err.Error())
	bot.Send(msg)
	msg = tgbotapi.NewMessage(db[ID].ID, "an error has occured. In order to proceed we need to recreate client and initialize new session")
	bot.Send(msg)
	updateDb := db[ID]
	updateDb.Dialog_status = 0
	db[ID] = updateDb
}

func StartDialogSequence(promt string, ID int64, ctx context.Context, bot *tgbotapi.BotAPI) {
	mu.Lock()
	defer mu.Unlock()
	userDatabase := database.UserMap
	sessionDatabase := database.AiSessionMap

	gpt_model := sessionDatabase[ID].Gpt_model
	log.Printf(
		"GPT model: %s,\npromt: %s\n",
		gpt_model,
		promt,
	)

	req := createComplexChatRequest(promt, gpt_model)
	c := sessionDatabase[ID].Gpt_client

	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		errorMessage(err, bot, ID, userDatabase)
	} else {
		fmt.Println(resp.Choices[0].Message.Content)
		resp_text := resp.Choices[0].Message.Content
		msg := tgbotapi.NewMessage(userDatabase[ID].ID, resp_text)
		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)
		updateDb := userDatabase[ID]
		updateDb.Dialog_status = 4
		userDatabase[ID] = updateDb
	}

}
