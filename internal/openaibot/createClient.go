package openaibot

import (
	db "github.com/JackBekket/telegram-gpt/internal/database"
	gogpt "github.com/sashabaranov/go-openai"
)

func CreateClient(apiKey string, ID int64, modelName string) *gogpt.Client {
	sessionDatabase := db.AiSessionMap
	client := gogpt.NewClient(apiKey)

	sessionDatabase[ID] = db.AiSession{
		Gpt_key:    apiKey,
		Gpt_client: client,
		Gpt_model:  modelName,
	}

	return client
}
