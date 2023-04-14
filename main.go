package main

import (
	"context"
	"fmt"
	"log"

	"github.com/JackBekket/telegram-gpt/internal/bot/command"
	"github.com/JackBekket/telegram-gpt/internal/bot/env"
	db "github.com/JackBekket/telegram-gpt/internal/database"

	//passport "github.com/MoonSHRD/IKY-telegram-bot/artifacts/TGPassport"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Panicf("could not load env from: %v", err)
	}

	token, err := env.LoadTGToken()
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("tg token missing: %v\n", err)
	}

	ctx := context.Background()

	// init database and commander
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

			fmt.Printf("ID: %v\nusername: %s\n",
				update.Message.From.ID,
				update.Message.From.UserName,
			)
			userID := update.Message.From.ID
			// Dialog_status for Admins = 2, other users = 0
			comm.CheckAdmin(userID, update.Message)

		} else {

			switch userDatabase[update.Message.From.ID].Dialog_status {

			// first check for user status, (for a new user status 0 is set automatically),
			// then user reply for the first bot message is logged to a database as name AND user status is updated
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
				}
			case 3:
				//  update Dialog_Status = 4, for model GPT-3.5
				comm.ConnectingToOpenAiWithLanguage(update.Message, ctx)
			case 4:
				// update update Dialog_Status = 4, for model GPT-3.5
				comm.DialogSequence(update.Message, ctx)
			case 5:
				comm.CodexSequence(update.Message, ctx)
			}

		}

	}

} // end of main func
