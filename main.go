package main

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"

	//passport "github.com/MoonSHRD/IKY-telegram-bot/artifacts/TGPassport"
	//passport "IKY-telegram-bot/artifacts/TGPassport"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var yesNoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
		tgbotapi.NewKeyboardButton("No")),
)

var optionKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("WhoIs")),
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Verify personal wallet")),
)

// to operate the bot, put a text file containing key for your bot acquired from telegram "botfather" to the same directory with this file
var tgApiKey, err = os.ReadFile(".secret")
var bot, error1 = tgbotapi.NewBotAPI(string(tgApiKey))

// type containing all the info about user input
type user struct {
	tgid          int64
	tg_username   string
	dialog_status int64
}


// main database for dialogs, key (int64) is telegram user id
var userDatabase = make(map[int64]user) // consider to change in persistend data storage?

var msgTemplates = make(map[string]string)


var myenv map[string]string

// file with settings for enviroment
const envLoc = ".env"

func main() {

	loadEnv()
	ctx := context.Background()
	//pk := myenv["PK"] // load private key from env

	msgTemplates["hello"] = "Hey, this bot is attaching personal wallets to telegram user & collective wallets to chat id"
	msgTemplates["case0"] = "Open following link in metamask broswer"
	msgTemplates["await"] = "Awaiting for verification"
	msgTemplates["case1"] = "You have successfully authorized your wallet to your account. Now you can use additional functions"
	msgTemplates["who_is"] = "Input wallet address to know it's associated telegram nickname"

	//var baseURL = "http://localhost:3000/"
	//var baseURL = "https://ikytest-gw0gy01is-s0lidarnost.vercel.app/"
	//var baseURL = myenv["BASEURL"]

	bot, err = tgbotapi.NewBotAPI(string(tgApiKey))
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

				userDatabase[update.Message.From.ID] = user{update.Message.Chat.ID, update.Message.Chat.UserName, 0}
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

						/*
						//link := baseURL + tg_id_query + tgid_string + tg_username_query + "@" + user_name
						msg = tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, link)
						bot.Send(msg)
						*/

						updateDb.dialog_status = 1
						userDatabase[update.Message.From.ID] = updateDb

					}
				//	fallthrough // МЫ ЛЕД ПОД НОГАМИ МАЙОРА!
				case 1:
					if updateDb, ok := userDatabase[update.Message.From.ID]; ok {
						msg := tgbotapi.NewMessage(userDatabase[update.Message.From.ID].tgid, msgTemplates["case1"])
						msg.ReplyMarkup = optionKeyboard
						bot.Send(msg)
						updateDb.dialog_status = 2
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

