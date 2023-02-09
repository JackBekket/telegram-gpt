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

var languageKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("eng"),
		tgbotapi.NewKeyboardButton("ru")),
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

	msgTemplates["hello"] = "Hey, this bot is OpenAI chatGPT. This is open beta, so I'm sustaining it at my laptop, so bot will be restarted oftenly"
	msgTemplates["case0"] = "Input your openAI API key. It can be created at https://platform.openai.com/account/api-keys"
	msgTemplates["await"] = "Awaiting"
	msgTemplates["case1"] = "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration."
	msgTemplates["codex_help"] = "``` # describe your task in comments like this or put your lines of code you need to autocomplete ```"


 
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
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
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
						//gpt3_m_string := gogpt.GPT3TextDavinci003
						ai_key := update.Message.Text
						//tg_username := update.Message.Chat.UserName
						//ai_client := CreateClient(ai_key)
						userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0,ai_key}
						// I can't validate key at this stage. The only way to validate key is to send test sequence (see case 3)
						// Since this part is oftenly get an uncaught exeption, we debug what user input as key. It's bad, I know, but until we got key validation we need this part.
						log.Println("key promt: ", ai_key)	
						updateDb.gpt_key = ai_key	// store key in memory
						
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
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "Choose language. Note that dataset used for training models in languages different from english may be *CENSORED*. This is problem with dataset, not model itself.")
							msg.ReplyMarkup = languageKeyboard
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
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["codex_help"])
							msg.ParseMode = "MARKDOWN"
							bot.Send(msg)
							//msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "Choose language. Note that dataset used for training models in languages different from english may be *CENSORED*. This is problem with dataset, not model itself")
							//msg.ReplyMarkup = languageKeyboard
							//bot.Send(msg)

							updateDb.dialog_status = 4
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
					
						if update.Message.Text == "eng" {
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "connecting to openAI")
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							ai_key := userDatabase[update.Message.From.ID].gpt_key
							un := userDatabase[update.Message.From.ID].tg_username
							ai_model := sessionDatabase[update.Message.From.ID].gpt_model
							go SetupSequenceWithKey(update.Message.From.ID,bot,ai_key,un,ai_model,update.Message.Text,ctx)

							
						}

						if update.Message.Text == "ru" {
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "connecting to openAI")
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							ai_key := userDatabase[update.Message.From.ID].gpt_key
							un := userDatabase[update.Message.From.ID].tg_username
							ai_model := sessionDatabase[update.Message.From.ID].gpt_model
							go SetupSequenceWithKey(update.Message.From.ID,bot,ai_key,un,ai_model,update.Message.Text,ctx)

							
						}

					
				case 4:
					//	if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						promt := update.Message.Text
						fmt.Println(promt)
						gpt_model := sessionDatabase[update.Message.From.ID].gpt_model
						log.Println(gpt_model)

						go StartDialogSequence(promt, update.Message.From.ID, ctx, *bot)


				case 5:
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

func SetupSequenceWithKey(tgid_ int64,bot *tgbotapi.BotAPI, ai_key string, tg_username string, gpt3_m_string_ string, language_ string, ctx context.Context) {
	/*
	msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid, msgTemplates["case0"])	// input key
	bot.Send(msg)
	*/

	ai_client := CreateClient(ai_key,tgid_,gpt3_m_string_)	// creating client (but we don't know if it works)
	//ai_model := 

	if language_ == "eng" {
		probe, err := TryLanguage(ai_client, gpt3_m_string_, 0,ctx)
		if err != nil {
			log.Println("error :", err)
							msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid,err.Error())
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[tgid_].tgid,"an error has occured. In order to proceed we need to recreate client and initialize new session")
							bot.Send(msg)
							updateDb := userDatabase[tgid_]
							updateDb.dialog_status = 0
							userDatabase[tgid_] = updateDb
		} else {
		log.Println(probe)
		msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid,probe)
		bot.Send(msg)

		userDatabase[tgid_] = user{tgid_, tg_username, 0,ai_key}
		sessionDatabase[tgid_] = ai_session{ai_key,ai_client,gpt3_m_string_}

		updateDb := userDatabase[tgid_]
		updateDb.dialog_status = 4
		userDatabase[tgid_] = updateDb
		}
	}

	if language_ == "ru" {
		probe, err := TryLanguage(ai_client, gpt3_m_string_, 1,ctx)
		if err != nil {
			log.Println("error :", err)
							msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid,err.Error())
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[tgid_].tgid,"an error has occured. In order to proceed we need to recreate client and initialize new session")
							bot.Send(msg)
							updateDb := userDatabase[tgid_]
							updateDb.dialog_status = 0
							userDatabase[tgid_] = updateDb
		} else {
		log.Println(probe)
		msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid,probe)
		bot.Send(msg)

		userDatabase[tgid_] = user{tgid_, tg_username, 0,ai_key}
		sessionDatabase[tgid_] = ai_session{ai_key,ai_client,gpt3_m_string_}

		updateDb := userDatabase[tgid_]
		updateDb.dialog_status = 4
		userDatabase[tgid_] = updateDb
		
		}
	}

}

func TryLanguage(client_ *gogpt.Client, model string, language int, ctx context.Context) (string, error){
	var language_promt string
	if language == 0 {
		language_promt = "Hi, Do you speak english?"
	}
	if language == 1 {
		language_promt = "Привет, ты говоришь по русски?"
	}
	log.Println(language_promt)
	req := CreateComplexRequest(language_promt, model)
	log.Println(req)
	resp, err := client_.CreateCompletion(ctx,req)
	if err != nil {
		return "nil",err
	} else {
		//return resp,nil
		answer := resp.Choices[0].Text
		return answer,err
	}
}

func StartDialogSequence(promt string, tgid int64, ctx context.Context, bot tgbotapi.BotAPI) {
	
	fmt.Println(promt)
	gpt_model := sessionDatabase[tgid].gpt_model
	log.Println(gpt_model)



	req := CreateComplexRequest(promt,gpt_model)
	c := sessionDatabase[tgid].gpt_client
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		//return
		log.Println("error :", err)
		msg := tgbotapi.NewMessage(userDatabase[tgid].tgid,err.Error())
		bot.Send(msg)
		msg = tgbotapi.NewMessage(userDatabase[tgid].tgid,"an error has occured. In order to proceed we need to recreate client and initialize new session")
		bot.Send(msg)
		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 0
		userDatabase[tgid] = updateDb

	} else {
	fmt.Println(resp.Choices[0].Text)
	resp_text := resp.Choices[0].Text
	msg := tgbotapi.NewMessage(userDatabase[tgid].tgid,resp_text)
	bot.Send(msg)
	updateDb := userDatabase[tgid]
	updateDb.dialog_status = 4
	userDatabase[tgid] = updateDb
	}
	
}

func CreateClient(AI_apiKey string, tgid int64,gpt3_m_string_ string) (*gogpt.Client){
	client := gogpt.NewClient(AI_apiKey)
	log.Println(client.ListEngines)
	log.Println(client.ListModels)
	log.Println(client.Answers)
	sessionDatabase[tgid] = ai_session{AI_apiKey,client,gpt3_m_string_}
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