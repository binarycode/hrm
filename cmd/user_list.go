package cmd

import (
	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
)

var userListCmd = &cobra.Command{
	Use:   "user:list",
	Short: "List users",
	Long:  "List all registered users",
}

func init() {
	userListCmd.Run = func(cmd *cobra.Command, args []string) {
		users, err := db.ListUsers()
		if err != nil {
			log.Fatal("Unable to query users", "err", err)
		}

		for _, user := range users {
			user.Display()
		}
	}
}
