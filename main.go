package main

import (
	"context"
	"fmt"
	"log"

	"os"

	gogpt "github.com/sashabaranov/go-openai"

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
		tgbotapi.NewKeyboardButton("GPT-3.5"),
		//tgbotapi.NewKeyboardButton("GPT-4"),
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
	gpt_key       string
	//gpt_client gpt3.Client
}

type ai_session struct {
	gpt_key    string
	gpt_client *gogpt.Client
	gpt_model  string
}

// main database for dialogs, key (int64) is telegram user id
var userDatabase = make(map[int64]user) // consider to change in persistend data storage?

var sessionDatabase = make(map[int64]ai_session)

var msgTemplates = make(map[string]string)

var myenv map[string]string

// file with settings for enviroment
const envLoc = ".env"

func main() {

	ctx := context.Background()

	msgTemplates["hello"] = "Hey, this bot is OpenAI chatGPT. This is open beta, so I'm sustaining it at my laptop, so bot will be restarted oftenly"
	msgTemplates["case0"] = "Input your openAI API key. It can be created at https://platform.openai.com/account/api-keys"
	msgTemplates["await"] = "Awaiting"
	msgTemplates["case1"] = "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration."
	msgTemplates["codex_help"] = "``` # describe your task in comments like this or put your lines of code you need to autocomplete ```"

	bot, err := tgbotapi.NewBotAPI(string(tgApiKey)[:len(string(tgApiKey))-1])
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

				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0, ""}
				msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["hello"])
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
				fmt.Println("tgid: ",update.Message.From.ID)
				fmt.Println("username: ",update.Message.From.String)
				fmt.Println("username: ",update.Message.From.UserName)
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
						userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0, ai_key}
						// I can't validate key at this stage. The only way to validate key is to send test sequence (see case 3)
						// Since this part is oftenly get an uncaught exeption, we debug what user input as key. It's bad, I know, but until we got key validation we need this part.
						log.Println("key promt: ", ai_key)
						updateDb.gpt_key = ai_key // store key in memory

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["case1"])
						msg.ReplyMarkup = chooseModelKeyboard
						bot.Send(msg)
						updateDb.dialog_status = 2
						userDatabase[update.Message.From.ID] = updateDb

					}
					//fallthrough
				case 2:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {

						if update.Message.Text == "GPT-3.5" {
							// TODO: Write down user choise
							log.Println(update.Message.Text)
							//gpt3_m_string := gogpt.GPT3TextDavinci003

							model_name := gogpt.GPT3Dot5Turbo // gpt-3.5
							
							log.Println(model_name)
							ai_client := sessionDatabase[update.Message.From.ID].gpt_client
							ai_key := sessionDatabase[update.Message.From.ID].gpt_key
							sessionDatabase[update.Message.From.ID] = ai_session{ai_key, ai_client, model_name}

							session_model := sessionDatabase[update.Message.From.ID].gpt_model
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "your session model :"+session_model)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language")
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
							sessionDatabase[update.Message.From.ID] = ai_session{ai_key, ai_client, gpt3_m_string}

							session_model := sessionDatabase[update.Message.From.ID].gpt_model
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "your session model :"+session_model)
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
						/*
							if update.Message.Text == "GPT-4" {
								log.Printf("buttom: %v\n", update.Message.Text)
								model_name := gogpt.GPT4 // gpt-4
								log.Printf("modelName: %v\n", model_name)

								ai_client := sessionDatabase[update.Message.From.ID].gpt_client
								ai_key := sessionDatabase[update.Message.From.ID].gpt_key
								sessionDatabase[update.Message.From.ID] = ai_session{ai_key, ai_client, model_name}

								session_model := sessionDatabase[update.Message.From.ID].gpt_model
								msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "your session model: "+session_model)
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
								bot.Send(msg)
								msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language")
								msg.ReplyMarkup = languageKeyboard
								bot.Send(msg)

								updateDb.dialog_status = 3
								userDatabase[update.Message.From.ID] = updateDb
							}
						*/
						// Can't use commands until connected to gpt chat
						if update.Message.Text != "GPT-3.5" && update.Message.Text != "Codex" && update.Message.Command() != "" {
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "type GPT-3.5 or Codex")
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
						go SetupSequenceWithKey(update.Message.From.ID, bot, ai_key, un, ai_model, update.Message.Text, ctx)

					}

					if update.Message.Text == "ru" {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "connecting to openAI")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)
						ai_key := userDatabase[update.Message.From.ID].gpt_key
						un := userDatabase[update.Message.From.ID].tg_username
						ai_model := sessionDatabase[update.Message.From.ID].gpt_model
						go SetupSequenceWithKey(update.Message.From.ID, bot, ai_key, un, ai_model, update.Message.Text, ctx)

					}

				case 4:

					if update.Message.Command() == "image" {

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, "Image link generation...")
						bot.Send(msg)

						promt := update.Message.CommandArguments()
						log.Printf("Command /image arg: %s\n", promt)

						go StartImageSequence(update.Message.From.ID, promt, ctx, bot, &update)

					} else {

						promt := update.Message.Text
						fmt.Println(promt)
						gpt_model := sessionDatabase[update.Message.From.ID].gpt_model
						log.Println(gpt_model)

						go StartDialogSequence(promt, update.Message.From.ID, ctx, *bot)
					}

				case 5:

					promt := update.Message.Text
					fmt.Println(promt)
					gpt_model := sessionDatabase[update.Message.From.ID].gpt_model
					log.Println(gpt_model)

					go StartCodexSequence(promt, update.Message.From.ID, ctx, *bot)
					//updateDb.dialog_status = 0
					//userDatabase[update.Message.From.ID] = updateDb

				}
			}
		}
	}

} // end of main func

func SetupSequenceWithKey(tgid_ int64, bot *tgbotapi.BotAPI, ai_key string, tg_username string, model_name string, language_ string, ctx context.Context) {
	/*
		msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid, msgTemplates["case0"])	// input key
		bot.Send(msg)
	*/

	ai_client := CreateClient(ai_key, tgid_, model_name) // creating client (but we don't know if it works)
	//ai_model :=

	if language_ == "eng" {
		probe, err := TryLanguage(ai_client, model_name, 0, ctx)
		if err != nil {
			log.Println("error :", err)
			msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid, err.Error())
			bot.Send(msg)
			msg = tgbotapi.NewMessage(userDatabase[tgid_].tgid, "an error has occured. In order to proceed we need to recreate client and initialize new session")
			bot.Send(msg)
			updateDb := userDatabase[tgid_]
			updateDb.dialog_status = 0
			userDatabase[tgid_] = updateDb
		} else {
			log.Println(probe)
			msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid, probe)
			bot.Send(msg)

			userDatabase[tgid_] = user{tgid_, tg_username, 0, ai_key}
			sessionDatabase[tgid_] = ai_session{ai_key, ai_client, model_name}

			updateDb := userDatabase[tgid_]
			updateDb.dialog_status = 4
			userDatabase[tgid_] = updateDb
		}
	}

	if language_ == "ru" {
		probe, err := TryLanguage(ai_client, model_name, 1, ctx)
		if err != nil {
			log.Println("error :", err)
			msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid, err.Error())
			bot.Send(msg)
			msg = tgbotapi.NewMessage(userDatabase[tgid_].tgid, "an error has occured. In order to proceed we need to recreate client and initialize new session")
			bot.Send(msg)
			updateDb := userDatabase[tgid_]
			updateDb.dialog_status = 0
			userDatabase[tgid_] = updateDb
		} else {
			log.Println(probe)
			msg := tgbotapi.NewMessage(userDatabase[tgid_].tgid, probe)
			bot.Send(msg)

			userDatabase[tgid_] = user{tgid_, tg_username, 0, ai_key}
			sessionDatabase[tgid_] = ai_session{ai_key, ai_client, model_name}

			updateDb := userDatabase[tgid_]
			updateDb.dialog_status = 4
			userDatabase[tgid_] = updateDb

		}
	}

}

func TryLanguage(client_ *gogpt.Client, model string, language int, ctx context.Context) (string, error) {
	var language_promt string
	if language == 0 {
		language_promt = "Hi, Do you speak english?"
	}
	if language == 1 {
		language_promt = "Привет, ты говоришь по русски?"
	}
	log.Println(language_promt)
	req := CreateComplexChatRequest(language_promt, model)
	log.Println(req)
	resp, err := client_.CreateChatCompletion(ctx, req)
	if err != nil {
		return "nil", err
	} else {
		//return resp,nil
		answer := resp.Choices[0].Message.Content
		return answer, err
	}
}

func StartDialogSequence(promt string, tgid int64, ctx context.Context, bot tgbotapi.BotAPI) {

	fmt.Println(promt)
	gpt_model := sessionDatabase[tgid].gpt_model
	log.Printf("GPT model: %s\n", gpt_model)

	req := CreateComplexChatRequest(promt, gpt_model)
	c := sessionDatabase[tgid].gpt_client
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		//return
		log.Println("error :", err)
		msg := tgbotapi.NewMessage(userDatabase[tgid].tgid, err.Error())
		bot.Send(msg)
		msg = tgbotapi.NewMessage(userDatabase[tgid].tgid, "an error has occured. In order to proceed we need to recreate client and initialize new session")
		bot.Send(msg)
		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 0
		userDatabase[tgid] = updateDb

	} else {
		fmt.Println(resp.Choices[0].Message.Content)
		resp_text := resp.Choices[0].Message.Content
		msg := tgbotapi.NewMessage(userDatabase[tgid].tgid, resp_text)
		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)
		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 4
		userDatabase[tgid] = updateDb
	}

}

func StartCodexSequence(promt string, tgid int64, ctx context.Context, bot tgbotapi.BotAPI) {

	gpt_model := sessionDatabase[tgid].gpt_model
	log.Println(gpt_model)

	req := CreateCodexRequest(promt)
	c := sessionDatabase[tgid].gpt_client
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		//return
		log.Println("error :", err)
		msg := tgbotapi.NewMessage(userDatabase[tgid].tgid, err.Error())
		bot.Send(msg)
		msg = tgbotapi.NewMessage(userDatabase[tgid].tgid, "an error has occured. In order to proceed we need to recreate client and initialize new session")
		bot.Send(msg)
		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 0
		userDatabase[tgid] = updateDb

	} else {
		fmt.Println(resp.Choices[0].Text)
		resp_text := resp.Choices[0].Text
		msg := tgbotapi.NewMessage(userDatabase[tgid].tgid, resp_text)
		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)
		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 5
		userDatabase[tgid] = updateDb
	}

}

func StartImageSequence(tgid int64, promt string, ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	req := CreateImageRequest(promt)
	c := sessionDatabase[tgid].gpt_client

	resp, err := c.CreateImage(ctx, req)
	if err != nil {

		log.Println("error :", err)
		msg := tgbotapi.NewMessage(userDatabase[tgid].tgid, err.Error())
		bot.Send(msg)
		msg = tgbotapi.NewMessage(userDatabase[tgid].tgid, "an error has occured. In order to proceed we need to recreate client and initialize new session")
		bot.Send(msg)
		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 0
		userDatabase[tgid] = updateDb

	} else {

		respUrl := resp.Data[0].URL
		log.Printf("url image: %s\n", respUrl)

		msg1 := tgbotapi.NewMessage(userDatabase[tgid].tgid, "Done!")
		bot.Send(msg1)

		msg := tgbotapi.NewEditMessageText(
			userDatabase[update.Message.From.ID].tgid,
			update.Message.MessageID+1,
			fmt.Sprintf("[Result](%s)", respUrl),
		)

		msg.ParseMode = "MARKDOWN"
		bot.Send(msg)

		updateDb := userDatabase[tgid]
		updateDb.dialog_status = 4
		userDatabase[tgid] = updateDb
	}

}

func CreateClient(AI_apiKey string, tgid int64, model_name string) *gogpt.Client {
	client := gogpt.NewClient(AI_apiKey)

	sessionDatabase[tgid] = ai_session{AI_apiKey, client, model_name}
	return client
}

/*
// used for GPT-3
func CreateSimpleTextRequest(input string) (gogpt.CompletionRequest){
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		MaxTokens: 2048,
		Prompt:    input,
		Echo: true,
	}
	return req
}
*/

// GPT-3.5
func CreateSimpleChatRequest(input string) gogpt.ChatCompletionRequest {
	req := gogpt.ChatCompletionRequest{
		Model:     gogpt.GPT3Dot5Turbo,
		MaxTokens: 3000,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role:    gogpt.ChatMessageRoleUser,
				Content: input,
			},
		}}
	return req
}

/*
// model should be gogpt.GPT3TextDavinci003 or gogpt.CodexCodeDavinci002
// WARN -- deprecated!
func CreateComplexRequest (input string, model_name string) (gogpt.CompletionRequest) {
	req := gogpt.CompletionRequest{
		Model: model_name,
		MaxTokens: 2048,
		Prompt: input,
		Echo: true,
	}
	return req
}
*/

func CreateComplexChatRequest(input string, model_name string) gogpt.ChatCompletionRequest {
	req := gogpt.ChatCompletionRequest{
		Model:     model_name,
		MaxTokens: 3000,
		Messages: []gogpt.ChatCompletionMessage{
			{
				Role:    gogpt.ChatMessageRoleUser,
				Content: input,
			},
		}}
	return req
}

// for code generation
func CreateCodexRequest(input string) gogpt.CompletionRequest {
	req := gogpt.CompletionRequest{
		Model:     gogpt.CodexCodeDavinci002,
		MaxTokens: 6000,
		Prompt:    input,
		Echo:      true,
	}
	return req
}

func CreateImageRequest(input string) gogpt.ImageRequest {
	return gogpt.ImageRequest{
		Prompt:         input,
		Size:           gogpt.CreateImageSize1024x1024,
		ResponseFormat: gogpt.CreateImageResponseFormatURL,
		N:              1,
	}
}
