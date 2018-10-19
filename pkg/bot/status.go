package bot

import (
	"strings"

	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

func status(user model.User) {
	services, err := db.ListSubscribedServices(user)
	if err != nil {
		log.Error("Unable to get subscribed services", "user", user, "err", err)
		send(user, "*ERROR* unable to get status")
		return
	}

	var b strings.Builder
	for _, service := range services {
		if service.Failure {
			b.WriteString("*FAILURE* ")
		} else if service.Maintenance {
			b.WriteString("*MAINTENANCE* ")
		} else {
			b.WriteString("*OK* ")
		}
		b.WriteString(escapeMarkdown(service.Name))
		b.WriteString("\n")
	}

	if b.Len() == 0 {
		send(user, "No subsribed services")
	} else {
		send(user, b.String())
	}
}
