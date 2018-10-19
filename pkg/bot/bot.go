package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

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

	if os.Getenv("DEBUG") != "" {
		bot.Debug = true
	}

	if config.Webhook {
		path := "/webhook/" + config.Token

		_, err = bot.SetWebhook(tgbotapi.NewWebhook(config.Address + path))
		if err != nil {
			log.Fatal("Unable to set webhook", "err", err)
		}

		http.HandleFunc(path, webhook)
	} else {
		_, err = bot.RemoveWebhook()
		if err != nil {
			log.Fatal("Unable to remove webhook", "err", err)
		}

		updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
			Offset:  0,
			Timeout: 60,
		})
		if err != nil {
			log.Fatal("Unable to initalize bot update channel", "err", err)
		}

		for update := range updates {
			process(update)
		}
	}
}

func webhook(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	process(update)
}

func process(update tgbotapi.Update) {
	if (update.Message == nil) || (!update.Message.IsCommand()) {
		return
	}

	from := update.Message.From

	user, err := db.GetUser(model.User{TelegramID: from.ID})
	if err != nil {
		user = model.User{
			TelegramID:   from.ID,
			ChatID:       update.Message.Chat.ID,
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

	switch update.Message.Command() {
	case "help", "?":
		help(user)
	case "start":
		start(user)
	case "status", "*":
		status(user)
	case "subscribe", "s", "+":
		subscribe(user, update)
	case "unsubscribe", "u", "-":
		unsubscribe(user, update)
	case "version", "v":
		version(user)
	}
}

func getUser(update tgbotapi.Update) (user model.User, err error) {
	from := update.Message.From

	user, err = db.GetUser(model.User{TelegramID: from.ID})
	if err == nil {
		return
	}

	user = model.User{
		TelegramID:   from.ID,
		ChatID:       update.Message.Chat.ID,
		FirstName:    from.FirstName,
		LastName:     from.LastName,
		UserName:     from.UserName,
		LanguageCode: from.LanguageCode,
		IsBot:        from.IsBot,
	}
	err = db.SaveUser(&user)
	return
}

func send(user model.User, text string) {
	message := tgbotapi.NewMessage(user.ChatID, text)
	message.ParseMode = "Markdown"
	bot.Send(message)
}

func escapeMarkdown(text string) string {
	text = strings.Replace(text, "*", "\\*", -1)
	text = strings.Replace(text, "_", "\\_", -1)
	text = strings.Replace(text, "[", "\\[", -1)
	text = strings.Replace(text, "]", "\\]", -1)
	text = strings.Replace(text, "`", "\\`", -1)

	return text
}
