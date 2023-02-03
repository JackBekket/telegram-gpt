package main

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"

	//gpt3 "github.com/PullRequestInc/go-gpt3"

	gogpt "github.com/sashabaranov/go-gpt3"

	//passport "github.com/MoonSHRD/IKY-telegram-bot/artifacts/TGPassport"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var yesNoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
		tgbotapi.NewKeyboardButton("No")),
)

var chooseModelKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("GPT3"),
		tgbotapi.NewKeyboardButton("Codex")),
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Start!")),
)

// to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKey, err = os.ReadFile(".secret")



// type containing all the info about user input
type user struct {
	tgid          int64
	tg_username   string
	dialog_status int64
	gpt_key string
	//gpt_client gpt3.Client
}

type ai_session struct {
	gpt_key		string
	gpt_client	*gogpt.Client
	gpt_model string
}


// main database for dialogs, key (int64) is telegram user id
var userDatabase = make(map[int64]user) // consider to change in persistend data storage?

var sessionDatabase = make(map[int64]ai_session)

var msgTemplates = make(map[string]string)


var myenv map[string]string

// file with settings for enviroment
const envLoc = ".env"

func main() {

	loadEnv()
	ctx := context.Background()
	//pk := myenv["PK"] // load private key from env

	msgTemplates["hello"] = "Hey, this bot is OpenAI chatGPT"
	msgTemplates["case0"] = "Input your openAI API key. It can be created at https://platform.openai.com/account/api-keys"
	msgTemplates["await"] = "Awaiting"
	msgTemplates["case1"] = "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration."



	bot, err := tgbotapi.NewBotAPI(string(tgApiKey))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {
			if _, ok := userDatabase[update.Message.From.ID]; !ok {

				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0,""}
				msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["hello"])
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
				// check for registration
			//	registred := IsAlreadyRegistred(session, update.Message.From.ID)
			/*
				if registred {
					userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 1}
				}
				*/
			
			} else {

				switch userDatabase[update.Message.From.ID].dialog_status {

				//first check for user status, (for a new user status 0 is set automatically), then user reply for the first bot message is logged to a database as name AND user status is updated
				case 0:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["case0"])
						bot.Send(msg)

						tgid := userDatabase[update.Message.From.ID].tgid
						user_name := userDatabase[update.Message.From.ID].tg_username
						fmt.Println(tgid)
						fmt.Println(user_name)

						updateDb.dialog_status = 1
						userDatabase[update.Message.From.ID] = updateDb

					}
					//fallthrough // МЫ ЛЕД ПОД НОГАМИ МАЙОРА!
					// 
				case 1:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						gpt3_m_string := gogpt.GPT3TextDavinci003
						ai_key := update.Message.Text
						ai_client := CreateClient(ai_key)
						userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0,ai_key}
						sessionDatabase[update.Message.From.ID] = ai_session{ai_key,ai_client,gpt3_m_string}
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["case1"])
						msg.ReplyMarkup = chooseModelKeyboard
						bot.Send(msg)
						updateDb.dialog_status = 2
						userDatabase[update.Message.From.ID] = updateDb

					}
					//fallthrough
				case 2:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						if update.Message.Text == "GPT3" {
							// TODO: Write down user choise
							log.Println(update.Message.Text)
							gpt3_m_string := gogpt.GPT3TextDavinci003
							
							log.Println(gpt3_m_string)
							ai_client := sessionDatabase[update.Message.From.ID].gpt_client
							ai_key := sessionDatabase[update.Message.From.ID].gpt_key
							sessionDatabase[update.Message.From.ID] = ai_session{ai_key,ai_client,gpt3_m_string}

							session_model := sessionDatabase[update.Message.From.ID].gpt_model
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "your session model :" + session_model)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "You are successfully connected to chatGPT now try to ask something")
							bot.Send(msg)

							updateDb.dialog_status = 3
							userDatabase[update.Message.From.ID] = updateDb
						} 
						if update.Message.Text == "Codex" {
							// Use codex model
							log.Println(update.Message.Text)
							gpt3_m_string := gogpt.CodexCodeDavinci002
							log.Println(gpt3_m_string)
							ai_client := sessionDatabase[update.Message.From.ID].gpt_client
							ai_key := sessionDatabase[update.Message.From.ID].gpt_key
							sessionDatabase[update.Message.From.ID] = ai_session{ai_key,ai_client,gpt3_m_string}

							session_model := sessionDatabase[update.Message.From.ID].gpt_model
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "your session model :" + session_model)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "You are successfully connected to chatGPT now try to ask something")
							bot.Send(msg)

							updateDb.dialog_status = 3
							userDatabase[update.Message.From.ID] = updateDb
						} 
						if update.Message.Text != "GPT3" && update.Message.Text != "Codex" {
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "type GPT3 or Codex")
							log.Println(update.Message.Text)
							bot.Send(msg)
							updateDb.dialog_status = 2
							userDatabase[update.Message.From.ID] = updateDb
						}
					}

				case 3:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						promt := update.Message.Text
						fmt.Println(promt)
						gpt_model := sessionDatabase[update.Message.From.ID].gpt_model
						log.Println(gpt_model)

						if update.Message.Text == "/restart" {
							updateDb.dialog_status = 4
							userDatabase[update.Message.From.ID] = updateDb
						} else {

						req := CreateComplexRequest(promt,gpt_model)
						c := sessionDatabase[update.Message.From.ID].gpt_client
						resp, err := c.CreateCompletion(ctx, req)
						if err != nil {
							//return
							log.Println("error :", err)
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid,err.Error())
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid,"an error has occured. In order to proceed we need to recreate client and initialize new session")
							bot.Send(msg)
							updateDb.dialog_status = 0
							userDatabase[update.Message.From.ID] = updateDb

						}
						fmt.Println(resp.Choices[0].Text)
						resp_text := resp.Choices[0].Text
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid,resp_text)
						bot.Send(msg)
						updateDb.dialog_status = 3
						userDatabase[update.Message.From.ID] = updateDb
						}
					}

				case 4:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						updateDb.dialog_status = 0
						userDatabase[update.Message.From.ID] = updateDb
					}
				}
			}
		}
	}

} // end of main func

// load enviroment variables from .env file
func loadEnv() {
	var err error
	if myenv, err = godotenv.Read(envLoc); err != nil {
		log.Printf("could not load env from %s: %v", envLoc, err)
	}
}

func CreateClient(AI_apiKey string) (*gogpt.Client){
	client := gogpt.NewClient(AI_apiKey)
	log.Println(client.ListEngines)
	log.Println(client.ListModels)
	log.Println(client.Answers)
	return client
}

func CreateSimpleRequest(input string) (gogpt.CompletionRequest){
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 2048,
		Prompt:    input,
		Echo: true,
	}
	return req
}

// model should be gogpt.GPT3TextDavinci003 or gogpt.CodexCodeDavinci002
func CreateComplexRequest (input string, model string) (gogpt.CompletionRequest) {
	req := gogpt.CompletionRequest{
		Model: model,
		MaxTokens: 2048,
		Prompt: input,
		Echo: true,
	}
	return req
}