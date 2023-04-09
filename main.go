package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	db "github.com/JackBekket/telegram-gpt/internal/database"
	aibot "github.com/JackBekket/telegram-gpt/internal/openaibot"
	"github.com/joho/godotenv"
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

var msgTemplates = make(map[string]string)

var myenv map[string]string

// file with settings for enviroment
const envLoc = ".env"

func main() {

	loadEnv()

	// constants from env
	ak := myenv["ADMIN_KEY"] //
	a_id_s := myenv["ADMIN_ID"]
	a_id, err := strconv.ParseInt(a_id_s, 0, 64)
	if err != nil {
		fmt.Println("error: ", err)
	}

	minty_k := myenv["MINTY_KEY"]
	minty_id_s := myenv["MINTY_ID"]
	minty_id, err := strconv.ParseInt(minty_id_s, 0, 64)
	if err != nil {
		fmt.Println("error: ", err)
	}

	ox_key := myenv["OK_KEY"]
	ox_id_s := myenv["OK_ID"]
	ox_id, err := strconv.ParseInt(ox_id_s, 0, 64)
	if err != nil {
		fmt.Println("error: ", err)
	}

	tg_key_env := myenv["TG_KEY"]

	ctx := context.Background()

	msgTemplates["hello"] = "Hey, this bot is OpenAI chatGPT. This is open beta, so I'm sustaining it at my laptop, so bot will be restarted oftenly"
	msgTemplates["case0"] = "Input your openAI API key. It can be created at https://platform.openai.com/account/api-keys"
	msgTemplates["await"] = "Awaiting"
	msgTemplates["case1"] = "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration."
	msgTemplates["codex_help"] = "``` # describe your task in comments like this or put your lines of code you need to autocomplete ```"

	/*
		bot, err := tgbotapi.NewBotAPI(string(tgApiKey)[:len(string(tgApiKey))-1])
		if err != nil {
			log.Panic(err)
		}
	*/

	bot, err := tgbotapi.NewBotAPI(tg_key_env)
	if err != nil {
		log.Fatalf("tg token: %v\n", err)
	}

	// init database
	userDatabase := db.UserMap
	sessionDatabase := db.AiSessionMap

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message != nil {
			if _, ok := userDatabase[update.Message.From.ID]; !ok {

				fmt.Println("ID: ", update.Message.From.ID)
				fmt.Println("username: ", update.Message.From.FirstName)
				fmt.Println("username: ", update.Message.From.UserName)
				user_id := update.Message.From.ID
				admin := false

				// if admin then get key from env
				if user_id == minty_id {
					ai_key := minty_k
					userDatabase[update.Message.From.ID] = db.User{update.Message.Chat.ID, update.Message.Chat.UserName, 2, ai_key}
					fmt.Println("minty authorized")
					admin = true
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "authorized")
					bot.Send(msg)
					msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["case1"])
					msg.ReplyMarkup = chooseModelKeyboard
					bot.Send(msg)
				}

				if user_id == a_id {
					ai_key := ak
					userDatabase[update.Message.From.ID] = db.User{update.Message.Chat.ID, update.Message.Chat.UserName, 2, ai_key}
					fmt.Println("a authorized")
					admin = true
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "authorized")
					bot.Send(msg)
					msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["case1"])
					msg.ReplyMarkup = chooseModelKeyboard
					bot.Send(msg)

				}
				if user_id == ox_id {
					ai_key := ox_key
					userDatabase[update.Message.From.ID] = db.User{update.Message.Chat.ID, update.Message.Chat.UserName, 2, ai_key}
					fmt.Println("ox authorized")
					admin = true
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "authorized")
					bot.Send(msg)
					msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["case1"])
					msg.ReplyMarkup = chooseModelKeyboard
					bot.Send(msg)

				} else if admin == false {

					userDatabase[update.Message.From.ID] = db.User{update.Message.Chat.ID, update.Message.Chat.UserName, 0, ""}
					msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["hello"])
					msg.ReplyMarkup = mainKeyboard
					bot.Send(msg)

					// check for registration
					//	registred := IsAlreadyRegistred(session, update.Message.From.ID)
					/*
						if registred {
							userDatabase[update.Message.From.ID] = db.User{update.Message.Chat.ID, update.Message.Chat.UserName, 1}
						}
					*/
				}

			} else {

				switch userDatabase[update.Message.From.ID].Dialog_status {

				//first check for user status, (for a new user status 0 is set automatically), then user reply for the first bot message is logged to a database as name AND user status is updated
				case 0:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["case0"])
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)

						ID := userDatabase[update.Message.From.ID].ID
						user_name := userDatabase[update.Message.From.ID].Username
						fmt.Println(ID)
						fmt.Println(user_name)

						updateDb.Dialog_status = 1
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
						userDatabase[update.Message.From.ID] = db.User{update.Message.Chat.ID, update.Message.Chat.UserName, 0, ai_key}
						// I can't validate key at this stage. The only way to validate key is to send test sequence (see case 3)
						// Since this part is oftenly get an uncaught exeption, we debug what user input as key. It's bad, I know, but until we got key validation we need this part.
						log.Println("key promt: ", ai_key)
						updateDb.Gpt_key = ai_key // store key in memory

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["case1"])
						msg.ReplyMarkup = chooseModelKeyboard
						bot.Send(msg)
						updateDb.Dialog_status = 2
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
							ai_client := sessionDatabase[update.Message.From.ID].Gpt_client
							ai_key := sessionDatabase[update.Message.From.ID].Gpt_key
							sessionDatabase[update.Message.From.ID] = db.AiSession{ai_key, ai_client, model_name}

							session_model := sessionDatabase[update.Message.From.ID].Gpt_model
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "your session model :"+session_model)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language")
							msg.ReplyMarkup = languageKeyboard
							bot.Send(msg)

							updateDb.Dialog_status = 3
							userDatabase[update.Message.From.ID] = updateDb
						}
						if update.Message.Text == "Codex" {
							// Use codex model
							log.Println(update.Message.Text)
							gpt3_m_string := gogpt.CodexCodeDavinci002
							log.Println(gpt3_m_string)
							ai_client := sessionDatabase[update.Message.From.ID].Gpt_client
							ai_key := sessionDatabase[update.Message.From.ID].Gpt_key
							sessionDatabase[update.Message.From.ID] = db.AiSession{ai_key, ai_client, gpt3_m_string}

							session_model := sessionDatabase[update.Message.From.ID].Gpt_model
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "your session model :"+session_model)
							msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
							bot.Send(msg)
							msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, msgTemplates["codex_help"])
							msg.ParseMode = "MARKDOWN"
							bot.Send(msg)
							//msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "Choose language. Note that dataset used for training models in languages different from english may be *CENSORED*. This is problem with dataset, not model itself")
							//msg.ReplyMarkup = languageKeyboard
							//bot.Send(msg)

							updateDb.Dialog_status = 4
							userDatabase[update.Message.From.ID] = updateDb
						}
						/*
							if update.Message.Text == "GPT-4" {
								log.Printf("buttom: %v\n", update.Message.Text)
								model_name := gogpt.GPT4 // gpt-4
								log.Printf("modelName: %v\n", model_name)

								ai_client := sessionDatabase[update.Message.From.ID].Gpt_client
								ai_key := sessionDatabase[update.Message.From.ID].Gpt_key
								sessionDatabase[update.Message.From.ID] = ai_session{ai_key, ai_client, model_name}

								session_model := sessionDatabase[update.Message.From.ID].Gpt_model
								msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "your session model: "+session_model)
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
								bot.Send(msg)
								msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language")
								msg.ReplyMarkup = languageKeyboard
								bot.Send(msg)

								updateDb.Dialog_status = 3
								userDatabase[update.Message.From.ID] = updateDb
							}
						*/
						// Can't use commands until connected to gpt chat
						if update.Message.Text != "GPT-3.5" && update.Message.Text != "Codex" && update.Message.Command() != "" {
							msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "type GPT-3.5 or Codex")
							log.Println(update.Message.Text)
							bot.Send(msg)
							updateDb.Dialog_status = 2
							userDatabase[update.Message.From.ID] = updateDb
						}
					}

				case 3:

					if update.Message.Text == "eng" {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "connecting to openAI")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)
						ai_key := userDatabase[update.Message.From.ID].Gpt_key
						un := userDatabase[update.Message.From.ID].Username
						ai_model := sessionDatabase[update.Message.From.ID].Gpt_model
						go aibot.SetupSequenceWithKey(update.Message.From.ID, bot, ai_key, un, ai_model, update.Message.Text, ctx)

					}

					if update.Message.Text == "ru" {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "connecting to openAI")
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
						bot.Send(msg)
						ai_key := userDatabase[update.Message.From.ID].Gpt_key
						un := userDatabase[update.Message.From.ID].Username
						ai_model := sessionDatabase[update.Message.From.ID].Gpt_model
						go aibot.SetupSequenceWithKey(update.Message.From.ID, bot, ai_key, un, ai_model, update.Message.Text, ctx)

					}

				case 4:

					if update.Message.Command() == "image" {

						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].ID, "Image link generation...")
						bot.Send(msg)

						promt := update.Message.CommandArguments()
						log.Printf("Command /image arg: %s\n", promt)

						go aibot.StartImageSequence(update.Message.From.ID, promt, ctx, bot, &update)

					} else {

						promt := update.Message.Text
						fmt.Println(promt)
						gpt_model := sessionDatabase[update.Message.From.ID].Gpt_model
						log.Println(gpt_model)

						go aibot.StartDialogSequence(promt, update.Message.From.ID, ctx, bot)
					}

				case 5:

					promt := update.Message.Text
					fmt.Println(promt)
					gpt_model := sessionDatabase[update.Message.From.ID].Gpt_model
					log.Println(gpt_model)

					go aibot.StartCodexSequence(promt, update.Message.From.ID, ctx, bot)
					//updateDb.Dialog_status = 0
					//userDatabase[update.Message.From.ID] = updateDb

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
