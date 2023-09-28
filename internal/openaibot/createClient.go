package openaibot

import (
	gogpt "github.com/sashabaranov/go-openai"
)

func CreateClient(gptKey string) *gogpt.Client {
	return gogpt.NewClient(gptKey)
}
