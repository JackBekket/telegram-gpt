package command

var msgTemplates = map[string]string{
	"hello":      "Hey, this bot is OpenAI chatGPT. This is open beta, so I'm sustaining it at my laptop, so bot will be restarted oftenly",
	"case0":      "Input your openAI API key. It can be created at https://platform.openai.com/account/api-keys",
	"await":      "Awaiting",
	"case1":      "Choose model to use. GPT3 is for text-based tasks, Codex for codegeneration.",
	"codex_help": "``` # describe your task in comments like this or put your lines of code you need to autocomplete ```",
}
