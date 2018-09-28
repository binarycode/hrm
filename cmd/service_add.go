package cmd

import (
	"time"

	"github.com/mgutz/logxi/v1"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

var serviceAddCmd = &cobra.Command{
	Use:   "service:add SERVICE_NAME",
	Short: "Register service",
	Long:  "Register new service with specified timeouts",
}

func init() {
	var failure, recovery, maintenance int

	serviceAddCmd.Flags().IntVarP(&failure, "failure", "f", 30, "failure timeout in seconds")
	serviceAddCmd.Flags().IntVarP(&recovery, "recovery", "r", 120, "recovery interval in seconds")
	serviceAddCmd.Flags().IntVarP(&maintenance, "maintenance", "m", 20, "maintenance timeout in minutes")
	serviceAddCmd.Args = cobra.ExactArgs(1)

	serviceAddCmd.Run = func(cmd *cobra.Command, args []string) {
		service := model.Service{
			Name:               args[0],
			Failure:            false,
			Recovering:         false,
			Maintenance:        false,
			FailureTimeout:     time.Second * time.Duration(failure),
			RecoveryInterval:   time.Second * time.Duration(recovery),
			MaintenanceTimeout: time.Minute * time.Duration(maintenance),
		}

		var min time.Duration

		min, _ = time.ParseDuration("5s")
		if service.FailureTimeout < min {
			log.Fatal("Failure timeout is too small", "value", service.FailureTimeout, "min", min)
		}

		min, _ = time.ParseDuration("1s")
		if service.RecoveryInterval < min {
			log.Fatal("Recovery interval is too small", "value", service.RecoveryInterval, "min", min)
		}

		min, _ = time.ParseDuration("1m")
		if service.MaintenanceTimeout < min {
			log.Fatal("Maintenance timeout is too small", "value", service.MaintenanceTimeout, "min", min)
		}

		_, err := db.GetService(model.Service{Name: service.Name})
		if err == nil {
			log.Fatal("Service already exists", "service", service)
		}

		token, err := uuid.NewV4()
		if err != nil {
			log.Fatal("Unable to generate service token", "err", err)
		}
		service.Token = token.String()

		if err := db.SaveService(&service); err != nil {
			log.Fatal("Unable to create service", "service", service, "err", err)
		}

		service.Display()
	}
}
