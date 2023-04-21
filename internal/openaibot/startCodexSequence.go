package openaibot

import (
	"context"
	"fmt"
	"log"

	db "github.com/JackBekket/telegram-gpt/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartCodexSequence(promt string, ID int64, ctx context.Context, bot *tgbotapi.BotAPI) {
	mu.Lock()
	defer mu.Unlock()
	userDatabase := db.UserMap
	sessionDatabase := db.AiSessionMap
	log.Printf(
		"GPT model: %s,\npromt: %s\n",
		sessionDatabase[ID].Gpt_model,
		promt,
	)

	req := createCodexRequest(promt)
	c := sessionDatabase[ID].Gpt_client

	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {

		errorMessage(err, bot, ID, userDatabase)

	} else {
		fmt.Println(resp.Choices[0].Text)
		resp_text := resp.Choices[0].Text
		msg := tgbotapi.NewMessage(userDatabase[ID].ID, resp_text)
		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)
		updateDb := userDatabase[ID]
		updateDb.Dialog_status = 5
		userDatabase[ID] = updateDb
	}

}
