package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

func unsubscribe(user model.User, update tgbotapi.Update) {
	name := update.Message.CommandArguments()

	if name == "" {
		send(user, "*ERROR* empty service name")
	}

	service, err := db.GetService(model.Service{Name: name})
	if err != nil {
		send(user, "*ERROR* service not found")
		return
	}

	if err = db.Unsubscribe(user, service); err != nil {
		log.Error("Unable to unsubscribe", "user", user, "service", service, "err", err)
		send(user, "*ERROR* unable to unsubscribe")
		return
	}

	send(user, "Successfully unsubscribed from "+escapeMarkdown(service.Name))
}
