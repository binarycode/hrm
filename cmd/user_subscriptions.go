package cmd

import (
	"strconv"

	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

var userSubscriptionsCmd = &cobra.Command{
	Use:   "user:subscriptions TELEGRAM_ID",
	Short: "List subscriptions",
	Long:  "List service subscriptions of specified user",
}

func init() {
	userSubscriptionsCmd.Args = cobra.ExactArgs(1)

	userSubscriptionsCmd.Run = func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("Unable to parse TELEGRAM_ID", "err", err)
		}

		if id == 0 {
			log.Fatal("TELEGRAM_ID cannot be zero")
		}

		user, err := db.GetUser(model.User{TelegramID: id})
		if err != nil {
			log.Fatal("Cannot find user", "id", id)
		}

		services, err := db.ListSubscribedServices(user)
		if err != nil {
			log.Fatal("Unable to query services", "err", err)
		}

		for _, service := range services {
			service.Display()
		}
	}
}
