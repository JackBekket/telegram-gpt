package openaibot

import (
	"log"

	gogpt "github.com/sashabaranov/go-openai"
	//"github.com/JackBekket/telegram-gpt/internal/bot/env"
)

func CreateClient(gptKey string) *gogpt.Client {
	return gogpt.NewClient(gptKey)
	
}

//
func CreateDefConfig(gptKey string, baseURL string) *gogpt.Client {
	cfg := gogpt.DefaultConfig(gptKey)
	cfg.BaseURL = baseURL
	return gogpt.NewClientWithConfig(cfg)	
}


// create anonymouse client with baseURL
func CreateCustomBaseConfig(baseURL string) *gogpt.Client {
	cfg := gogpt.DefaultConfig("")
	cfg.BaseURL = baseURL
	return gogpt.NewClientWithConfig(cfg)	
}


func CreateLocalhostClient() *gogpt.Client {
	cfg := gogpt.DefaultConfig("")
	cfg.BaseURL = "http://127.0.0.1:8080"
	return gogpt.NewClientWithConfig(cfg)
}

func CreateLocalhostClientWithCheck(lpwd string,user_promt string) *gogpt.Client {
	if (lpwd == user_promt) {
		log.Println(lpwd)
		log.Println(user_promt)
		log.Println("creating localhost client")
		cfg := gogpt.DefaultConfig(user_promt)
		cfg.BaseURL = "http://127.0.0.1:8080"
		return gogpt.NewClientWithConfig(cfg)
	} else {
		log.Println("creating connection to open-ai")
		return CreateClient(user_promt)
	}

}
