package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

func subscribe(user model.User, update tgbotapi.Update) {
	token := update.Message.CommandArguments()

	if token == "" {
		send(user, "*ERROR* empty token")
	}

	service, err := db.GetService(model.Service{Token: token})
	if err != nil {
		send(user, "*ERROR* service not found")
		return
	}

	if err = db.Subscribe(user, service); err != nil {
		log.Error("Unable to subscribe", "user", user, "service", service, "err", err)
		send(user, "*ERROR* unable to subscribe")
		return
	}

	send(user, "Successfully subscribed to "+escapeMarkdown(service.Name))
}
