package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

type Config struct {
	Address string
	Proxy   string
	Token   string
	Webhook bool
}

var bot *tgbotapi.BotAPI

func Start(config Config) {
	var err error

	client := &http.Client{}
	if config.Proxy != "" {
		url, err := url.Parse(config.Proxy)
		if err != nil {
			log.Fatal("Unable to parse proxy url", "err", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(url)}
	}

	bot, err = tgbotapi.NewBotAPIWithClient(config.Token, client)
	if (err != nil) || (bot == nil) {
		log.Fatal("Unable to initialize bot", "err", err)
	}
	log.Info("Started bot", "bot", bot.Self.UserName)

	if config.Webhook {
		path := "/webhook/" + config.Token

		_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.Address + path))
		if err != nil {
			log.Fatal("Unable to set webhook", "err", err)
		}

		http.HandleFunc(path, webhook)
	} else {
		updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 60,
		})
		if err != nil {
			log.Fatal("Unable to initalize bot update channel", "err", err)
		}

		for update := range updates {
			process(update.Message)
		}
	}
}

func webhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	process(update.Message)
}

func process(message *tgbotapi.Message) {
	if message == nil {
		return
	}

	from := message.From
	user, err := db.GetUser(model.User{TelegramID: from.ID})
	if err != nil {
		user = model.User{
			TelegramID:   from.ID,
			ChatID:       message.Chat.ID,
			FirstName:    from.FirstName,
			LastName:     from.LastName,
			UserName:     from.UserName,
			LanguageCode: from.LanguageCode,
			IsBot:        from.IsBot,
		}
		err = db.SaveUser(&user)
	}
	if err != nil {
		log.Error("Unable to get user", "user", user, "err", err)
		return
	}

	if message.IsCommand() {
		command(user, message)
	} else if message.Photo != nil {
		for _, photo := range (*message.Photo) {
			token, err := parseQRCode(photo)
			if err != nil {
				send(user, "Cannot parse QR Code")
				return
			}
			subscribe(user, token)
		}
	}
}

func command(user model.User, message *tgbotapi.Message) {
	args := message.CommandArguments()

	switch message.Command() {
	case "help", "?":
		help(user)
	case "start":
		start(user)
	case "status", "*":
		status(user)
	case "subscribe", "s", "+":
		subscribe(user, args)
	case "unsubscribe", "u", "-":
		unsubscribe(user, args)
	}
}

func send(user model.User, text string) {
	message := tgbotapi.NewMessage(user.ChatID, text)
	message.ParseMode = "Markdown"
	bot.Send(message)
}

func parseQRCode(photo *message.PhotoSize) (token string, err error) {
	url, err := bot.GetFileDirectURL(photo.FileID)
	if err != nil {
		log.Error("Unable to get direct URL for a file", "photo", photo, "err", err)
		return
	}

	response, err := http.Get(url)
	if err != nil {
		log.Error("Unable to download file", "url", url, "err", err)
		return
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile("", "photo")
	if err != nil {
		log.Error("Unable to create tempfile", "err", err)
		return
	}
	defer os.Remove(file.Name())

	_, err := io.Copy(out, response.Body)
	if err != nil {
		log.Error("Unable to write the response to file", "err", err)
		return
	}

	err := os.Close(file)
	if err != nil {
		log.Error("Unable to close tempfile", "err", err)
		return
	}
}
