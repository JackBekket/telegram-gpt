package command

import (
	"log"

	"github.com/JackBekket/telegram-gpt/internal/openaibot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

// Message:	case0 - "Input your openAI API key. It can be created at https://platform.openai.com/accousernamet/api-keys".
//
//	update DialogStatus = 1
func (c *Commander) InputYourAPIKey(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(
		user.ID,
		msgTemplates["case0"],
	)
	c.bot.Send(msg)

	user.DialogStatus = 1
	c.usersDb[chatID] = user
}

// Message: case1 - "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration.".
//
//	update Dialog_Status = 2
func (c *Commander) ChooseModel(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	gptKey := updateMessage.Text
	user := c.usersDb[chatID]
	// I can't validate key at this stage. The only way to validate key is to send test sequence (see case 3)
	// Since this part is oftenly get an usernamecaught exeption, we debug what user input as key. It's bad, I know, but usernametil we got key validation we need this part.
	log.Println("Key promt: ", gptKey)
	user.AiSession.GptKey = gptKey // store key in memory

	msg := tgbotapi.NewMessage(user.ID, msgTemplates["case1"])
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("GPT-3.5")),
		//tgbotapi.NewKeyboardButton("GPT-4"),
		//tgbotapi.NewKeyboardButton("Codex")),
	)
	c.bot.Send(msg)

	user.DialogStatus = 2
	c.usersDb[chatID] = user
}

// Message: "Choose language. If you have different languages then listed, then just send 'Hello' at your desired language".
//
//	update Dialog_Status = 3
func (c *Commander) ModelGPT3DOT5(updateMessage *tgbotapi.Message) {
	// TODO: Write down user choise
	log.Printf("Model selected: %s\n", updateMessage.Text)

	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	modelName := openai.GPT3Dot5Turbo // gpt-3.5
	user.AiSession.GptModel = modelName
	msg := tgbotapi.NewMessage(user.ID, "your session model: "+modelName)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(user.ID, "Choose a language or send 'Hello' in your desired language.")
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("English"),
			tgbotapi.NewKeyboardButton("Russian")),
	)
	c.bot.Send(msg)

	user.DialogStatus = 3
	c.usersDb[chatID] = user
}

// Message: "your session model: Codex".
//
//	update Dialog_Status = 4
func (c *Commander) ModelCodex(updateMessage *tgbotapi.Message) {
	log.Printf("Model selected: %s\n", updateMessage.Text)
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	modelCodex := openai.CodexCodeDavinci002
	user.AiSession.GptModel = modelCodex

	msg := tgbotapi.NewMessage(user.ID, "your session model :"+modelCodex)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(user.ID, msgTemplates["codex_help"])
	msg.ParseMode = "MARKDOWN"
	c.bot.Send(msg)

	user.DialogStatus = 4
	c.usersDb[chatID] = user
}

// ModelGPT and ModelLL codes are the same.
// TODO
func (c *Commander) ModelGPT4(updateMessage *tgbotapi.Message) {
	// TODO: Write down user choise
	log.Printf("Model selected: %s\n", updateMessage.Text)

	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	modelName := openai.GPT4 // gpt-4
	user.AiSession.GptModel = modelName
	msg := tgbotapi.NewMessage(user.ID, "your session model: "+modelName)
	c.bot.Send(msg)

	msg = tgbotapi.NewMessage(user.ID, "Choose a language or send 'Hello' in your desired language.")
	msg.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("English"),
			tgbotapi.NewKeyboardButton("Russian")),
	)
	c.bot.Send(msg)

	user.DialogStatus = 3
	c.usersDb[chatID] = user
}

// update Dialog_Status = 2
func (c *Commander) WrongModel(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(user.ID, "type GPT-3.5")
	c.bot.Send(msg)

	user.DialogStatus = 2
	c.usersDb[chatID] = user
}

// Message: "connecting to openAI"
//
// update update Dialog_Status = 4, for model GPT-3.5
func (c *Commander) ConnectingToOpenAiWithLanguage(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	language := updateMessage.Text
	user := c.usersDb[chatID]

	msg := tgbotapi.NewMessage(user.ID, "connecting to openAI")
	c.bot.Send(msg)

	go openaibot.SetupSequenceWithKey(c.bot, user, language, c.ctx)
}

// Generates an image with the /image command.
//
// Generates and sends text to the user.
//
// update Dialog_Status = 4, for model GPT-3.5
func (c *Commander) DialogSequence(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	user := c.usersDb[chatID]
	switch updateMessage.Command() {
	case "image":
		msg := tgbotapi.NewMessage(user.ID, "Image link generation...")
		c.bot.Send(msg)

		promt := updateMessage.CommandArguments()
		log.Printf("Command /image arg: %s\n", promt)
		go openaibot.StartImageSequence(c.bot, updateMessage, chatID, promt, c.ctx)

	default:
		promt := updateMessage.Text
		go openaibot.StartDialogSequence(c.bot, chatID, promt, c.ctx)
	}
}

// Generates and sends code to the user.
//
// At the moment there is no access to the Codex.
func (c *Commander) CodexSequence(updateMessage *tgbotapi.Message) {
	chatID := updateMessage.From.ID
	promt := updateMessage.Text
	go openaibot.StartCodexSequence(c.bot, chatID, promt, c.ctx)
	//user.DialogStatus = 0
	//userDatabase[chatID] = user
}
