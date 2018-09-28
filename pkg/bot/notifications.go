package bot

import (
	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

func MaintenanceFailure(service model.Service) {
	log.Warn("Maintenance took too long", "service", service.Name)
	notify(service, "*MAINTENANCE FAILURE* "+service.Name)
}

func Failure(service model.Service) {
	log.Warn("Service failed", "service", service.Name)
	notify(service, "*FAILURE* "+service.Name)
}

func Recovery(service model.Service) {
	log.Warn("Service recovered", "service", service.Name)
	notify(service, "*RECOVERY* "+service.Name)
}

func notify(service model.Service, text string) {
	users, err := db.ListSubscribedUsers(service)
	if err != nil {
		log.Error("Unable to query users subscribed to service", "service", service, "err", err)
		return
	}

	for _, user := range users {
		send(user, text)
	}
}
