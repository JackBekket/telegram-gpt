package database

import gogpt "github.com/sashabaranov/go-openai"

type AiSession struct {
	Gpt_key    string
	Gpt_client *gogpt.Client
	Gpt_model  string
}

var AiSessionMap = make(map[int64]AiSession)

// func NewAiSessionDataBase() *map[int64]*AiSession {
// 	db := make(map[int64]*AiSession)
// 	return &db
// }
