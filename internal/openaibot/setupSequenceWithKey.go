package openaibot

import (
	"context"
	"log"

	"github.com/JackBekket/telegram-gpt/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gogpt "github.com/sashabaranov/go-openai"
)

func SetupSequenceWithKey(
	ID int64,
	bot *tgbotapi.BotAPI,
	aiKey, username, model, language string,
	ctx context.Context,
) {
	userDatabase := database.UserMap
	sessionDatabase := database.AiSessionMap

	client := CreateClient(aiKey, ID, model) // creating client (but we don't know if it works)

	if language == "eng" {
		probe, err := tryLanguage(client, model, 0, ctx)
		if err != nil {
			errorMessage(err, bot, ID, userDatabase)
		} else {
			probeAnswer(probe, username, aiKey, ID, client, model, bot, userDatabase, sessionDatabase)
		}
	}

	if language == "ru" {
		probe, err := tryLanguage(client, model, 1, ctx)
		if err != nil {
			errorMessage(err, bot, ID, userDatabase)
		} else {
			probeAnswer(probe, username, aiKey, ID, client, model, bot, userDatabase, sessionDatabase)

		}
	}

}

func tryLanguage(client *gogpt.Client, model string, language int, ctx context.Context) (string, error) {
	var languagepromt string
	if language == 0 {
		languagepromt = "Hi, Do you speak english?"
	}
	if language == 1 {
		languagepromt = "Привет, ты говоришь по русски?"
	}
	log.Printf("Language: %v\n", languagepromt)

	req := createComplexChatRequest(languagepromt, model)
	log.Printf("request: %v\n", req)

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	} else {
		answer := resp.Choices[0].Message.Content
		return answer, nil
	}
}

func probeAnswer(
	probe, username, aiKey string,
	ID int64,
	client *gogpt.Client,
	model string,
	bot *tgbotapi.BotAPI,
	userDatabase map[int64]database.User,
	sessionDatabase map[int64]database.AiSession,
) {

	msg := tgbotapi.NewMessage(userDatabase[ID].ID, probe)
	bot.Send(msg)

	userDatabase[ID] = database.User{
		ID:            ID,
		Username:      username,
		Dialog_status: 0,
		Gpt_key:       aiKey}

	sessionDatabase[ID] = database.AiSession{
		Gpt_key:    aiKey,
		Gpt_client: client,
		Gpt_model:  model,
	}

	updateDb := userDatabase[ID]
	updateDb.Dialog_status = 4
	userDatabase[ID] = updateDb
}
