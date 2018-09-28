package cmd

import (
	"strconv"

	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

var userRemoveCmd = &cobra.Command{
	Use:   "user:remove TELEGRAM_ID",
	Short: "Remove user",
	Long:  "Remove specified user",
}

func init() {
	userRemoveCmd.Args = cobra.ExactArgs(1)

	userRemoveCmd.Run = func(cmd *cobra.Command, args []string) {
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

		if err := db.DestroyUser(&user); err != nil {
			log.Fatal("Unable to destroy user", "user", user, "err", err)
		}
	}
}
