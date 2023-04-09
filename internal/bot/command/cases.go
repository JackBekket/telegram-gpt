package command

import (
	"context"
	"fmt"
	"log"

	"github.com/JackBekket/telegram-gpt/internal/database"
	"github.com/JackBekket/telegram-gpt/internal/openaibot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

// Message:	case0 - "Input your openAI API key. It can be created at https://platform.openai.com/accousernamet/api-keys".
//
//	update Dialog_status = 1
func (c *Commander) InputYourAPIKey(updateMessage *tgbotapi.Message, updateDb *database.User) {
	msg := tgbotapi.NewMessage(
		c.userDb[updateMessage.From.ID].ID,
		msgTemplates["case0"],
	)
	c.bot.Send(msg)

	log.Printf(
		"New user: id: %v\n\t\t\t\tusername: %s\n",
		c.userDb[updateMessage.From.ID].ID,
		c.userDb[updateMessage.From.ID].Username,
	)

	updateDb.Dialog_status = 1
	c.userDb[updateMessage.From.ID] = *updateDb
}

// Message: case1 - "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration.".
//
//	update Dialog_Status = 2
func (c *Commander) ChooseModel(updateMessage *tgbotapi.Message, updateDb *database.User) {
	aiKey := updateMessage.Text
	c.userDb[updateMessage.From.ID] = database.User{
		ID:            updateMessage.Chat.ID,
		Username:      updateMessage.Chat.UserName,
		Dialog_status: 0,
		Gpt_key:       aiKey,
	}
	// I can't validate key at this stage. The only way to validate key is to send test sequence (see case 3)
	// Since this part is oftenly get an usernamecaught exeption, we debug what user input as key. It's bad, I know, but usernametil we got key validation we need this part.
	log.Println("Key promt: ", aiKey)
	updateDb.Gpt_key = aiKey // store key in memory

	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, msgTemplates["case1"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("GPT-3.5"),
			//tgbotapi.NewKeyboardButton("GPT-4"),
			tgbotapi.NewKeyboardButton("Codex")),
	)
	c.bot.Send(msg)

	updateDb.Dialog_status = 2
	c.userDb[updateMessage.From.ID] = *updateDb
}

// Message: "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language".
//
//	update Dialog_Status = 3
func (c *Commander) ModelGPT3DOT5(updateMessage *tgbotapi.Message, updateDb *database.User) {
	// TODO: Write down user choise
	log.Printf("Model selected: %s\n", updateMessage.Text)

	modelName := openai.GPT3Dot5Turbo // gpt-3.5
	client := c.aiSessionDb[updateMessage.From.ID].Gpt_client
	key := c.aiSessionDb[updateMessage.From.ID].Gpt_key
	c.aiSessionDb[updateMessage.From.ID] = database.AiSession{
		Gpt_key:    key,
		Gpt_client: client,
		Gpt_model:  modelName,
	}

	sessionModel := c.aiSessionDb[updateMessage.From.ID].Gpt_model
	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "your session model :"+sessionModel)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language")
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("eng"),
			tgbotapi.NewKeyboardButton("ru")),
	)
	c.bot.Send(msg)

	updateDb.Dialog_status = 3
	c.userDb[updateMessage.From.ID] = *updateDb
}

// Message: "your session model: Codex".
//
//	update Dialog_Status = 4
func (c *Commander) ModelCodex(updateMessage *tgbotapi.Message, updateDb *database.User) {
	log.Printf("Model selected: %s\n", updateMessage.Text)

	codex := openai.CodexCodeDavinci002
	client := c.aiSessionDb[updateMessage.From.ID].Gpt_client
	key := c.aiSessionDb[updateMessage.From.ID].Gpt_key
	c.aiSessionDb[updateMessage.From.ID] = database.AiSession{
		Gpt_key:    key,
		Gpt_client: client,
		Gpt_model:  codex,
	}

	sessionModel := c.aiSessionDb[updateMessage.From.ID].Gpt_model
	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "your session model :"+sessionModel)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, msgTemplates["codex_help"])
	msg.ParseMode = "MARKDOWN"
	c.bot.Send(msg)
	//msg = tgbotapi.NewMessage(
	// 	c.userDb[updateMessage.From.ID].ID,
	// 	"Choose language. Note that dataset used for training models in languages different from english may be *CENSORED*. This is problem with dataset, not model itself",
	// )
	//msg.ReplyMarkup = languageKeyboard
	//bot.Send(msg)

	updateDb.Dialog_status = 4
	c.userDb[updateMessage.From.ID] = *updateDb
}

// ModelGPT and ModelLL codes are the same.
// TODO
func (c *Commander) ModelGPT4(updateMessage *tgbotapi.Message, updateDb *database.User) {
	// TODO: Write down user choise
	log.Printf("Model selected: %s\n", updateMessage.Text)

	modelName := openai.GPT4 // ModelGPT3DOT5 and ModeGPT4 code are the same except for this line.
	client := c.aiSessionDb[updateMessage.From.ID].Gpt_client
	key := c.aiSessionDb[updateMessage.From.ID].Gpt_key
	c.aiSessionDb[updateMessage.From.ID] = database.AiSession{
		Gpt_key:    key,
		Gpt_client: client,
		Gpt_model:  modelName,
	}

	sessionModel := c.aiSessionDb[updateMessage.From.ID].Gpt_model
	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "your session model :"+sessionModel)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language")
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("eng"),
			tgbotapi.NewKeyboardButton("ru")),
	)
	c.bot.Send(msg)

	updateDb.Dialog_status = 3
	c.userDb[updateMessage.From.ID] = *updateDb
}

// update Dialog_Status = 2
func (c *Commander) WrongModel(updateMessage *tgbotapi.Message, updateDb *database.User) {
	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "type GPT-3.5 or Codex")
	log.Println(updateMessage.Text)
	c.bot.Send(msg)
	updateDb.Dialog_status = 2
	c.userDb[updateMessage.From.ID] = *updateDb
}

// Message: "connecting to openAI"
//
// update update Dialog_Status = 4, for model GPT-3.5
func (c *Commander) ConnectingToOpenAiWithLanguage(
	updateMessage *tgbotapi.Message,
	ctx context.Context,

) {
	msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "connecting to openAI")
	c.bot.Send(msg)

	aikey := c.userDb[updateMessage.From.ID].Gpt_key
	username := c.userDb[updateMessage.From.ID].Username
	aiModel := c.aiSessionDb[updateMessage.From.ID].Gpt_model
	language := updateMessage.Text
	go openaibot.SetupSequenceWithKey(updateMessage.From.ID, c.bot, aikey, username, aiModel, language, ctx)
}

// Generates an image with the /image command.
//
// Generates and sends text to the user.
//
// update Dialog_Status = 4, for model GPT-3.5
func (c *Commander) DialogSequence(
	updateMessage *tgbotapi.Message,
	ctx context.Context,
) {
	switch updateMessage.Command() {
	case "image":

		msg := tgbotapi.NewMessage(c.userDb[updateMessage.From.ID].ID, "Image link generation...")
		c.bot.Send(msg)

		promt := updateMessage.CommandArguments()
		log.Printf("Command /image arg: %s\n", promt)
		go openaibot.StartImageSequence(updateMessage.From.ID, promt, ctx, c.bot, updateMessage)

	default:
		promt := updateMessage.Text
		fmt.Println(promt)
		gpt_model := c.aiSessionDb[updateMessage.From.ID].Gpt_model
		log.Println(gpt_model)

		go openaibot.StartDialogSequence(promt, updateMessage.From.ID, ctx, c.bot)
	}
}

// Generates and sends code to the user.
//
// At the moment there is no access to the Codex.
func (c *Commander) CodexSequence(
	updateMessage *tgbotapi.Message,
	ctx context.Context,
) {
	promt := updateMessage.Text
	fmt.Println(promt)
	gpt_model := c.aiSessionDb[updateMessage.From.ID].Gpt_model
	log.Println(gpt_model)

	go openaibot.StartCodexSequence(promt, updateMessage.From.ID, ctx, c.bot)
	//updateDb.Dialog_status = 0
	//userDatabase[updateMessage.From.ID] = updateDb
}
