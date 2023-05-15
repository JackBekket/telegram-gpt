package openaibot

import (
	db "github.com/JackBekket/telegram-gpt/internal/database"
	gogpt "github.com/sashabaranov/go-openai"
)

func CreateClient(gptKey string, chatID int64, modelName string) *gogpt.Client {

	client := gogpt.NewClient(gptKey)
	db.UsersMap[chatID] = db.User{
		AiSession: db.AiSession{
			GptKey:    gptKey,
			GptClient: *client,
			GptModel:  modelName,
		},
	}

	return client
}
