package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/JackBekket/telegram-gpt/internal/bot/command"
	db "github.com/JackBekket/telegram-gpt/internal/database"
	aibot "github.com/JackBekket/telegram-gpt/internal/openaibot"
	"github.com/joho/godotenv"

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
	comm := command.NewCommander(bot, userDatabase, sessionDatabase)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	//whenever bot gets a new message, check for user id in the database happens, if it's a new user, the entry in the database is created.
	for update := range updates {

		if update.Message == nil {
			continue
		}
		if _, ok := userDatabase[update.Message.From.ID]; !ok {

			fmt.Println("ID: ", update.Message.From.ID)
			fmt.Println("username: ", update.Message.From.FirstName)
			fmt.Println("username: ", update.Message.From.UserName)
			user_id := update.Message.From.ID
			admin := false

			// if admin then get key from env
			if user_id == minty_id {
				admin = comm.AddAdminToMap(minty_k, update.Message)
			}

			if user_id == a_id {
				admin = comm.AddAdminToMap(ak, update.Message)

			}
			if user_id == ox_id {
				admin = comm.AddAdminToMap(ox_key, update.Message)

			}
			if !admin {

				comm.AddNewUserToMap(update.Message)
			}

		} else {

			switch userDatabase[update.Message.From.ID].Dialog_status {

			//first check for user status, (for a new user status 0 is set automatically), then user reply for the first bot message is logged to a database as name AND user status is updated
			case 0:
				if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
					// update Dialog_status = 1
					comm.InputYourAPIKey(update.Message, &updateDb)
				}
			case 1:
				if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
					// update Dialog_status = 2
					comm.ChooseModel(update.Message, &updateDb)
				}
			case 2:
				if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
					switch update.Message.Text {
					case "GPT-3.5":
						// update Dialog_status = 3
						comm.ModelGPT3DOT5(update.Message, &updateDb)
					case "Codex":
						// update Dialog_status = 4
						comm.ModelCodex(update.Message, &updateDb)
					// case "GPT-4":
					// 	comm.ModelGPT4(update.Message, &updateDb)
					default:
						comm.WrongModel(update.Message, &updateDb)
					}
					// Can't use commands until connected to gpt chat
					// update.Message.Command() != ""
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

} // end of main func

// load enviroment variables from .env file
func loadEnv() {
	var err error
	if myenv, err = godotenv.Read(envLoc); err != nil {
		log.Printf("could not load env from %s: %v", envLoc, err)
	}
}
