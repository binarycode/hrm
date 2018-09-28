package monitor

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mgutz/logxi/v1"

	"github.com/binarycode/trewoga/pkg/bot"
	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

func Start() {
	known := []string{""}

	for {
		unknown, err := db.ListScopedServices(func(db *gorm.DB) *gorm.DB {
			return db.Not("Token", known)
		})
		if err != nil {
			log.Fatal("Unable to query unknown services", "err", err)
		}

		for _, service := range unknown {
			go monitor(service)
			known = append(known, service.Token)
		}

		time.Sleep(10 * time.Second)
	}
}

func monitor(service model.Service) {
	for {
		time.Sleep(service.FailureTimeout)

		service, err := db.GetService(model.Service{Token: service.Token})
		if gorm.IsRecordNotFoundError(err) {
			return
		} else if err != nil {
			log.Fatal("Unable to find service", "err", err)
		}

		if service.Maintenance {
			if time.Now().After(service.MaintenanceFailureAt) {
				bot.MaintenanceFailure(service)
				service.Maintenance = false
			}
		} else {
			if time.Now().After(service.PingAt.Add(service.FailureTimeout)) {
				if !service.Failure {
					bot.Failure(service)
					service.Failure = true
				}

				service.Recovering = false
			}

			if service.Failure && service.Recovering && time.Now().After(service.RecoveryAt) {
				bot.Recovery(service)
				service.Failure = false
			}
		}

		if err = db.SaveService(&service); err != nil {
			log.Fatal("Unable to save service", "err", err)
		}
	}
}
