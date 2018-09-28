package cmd

import (
	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
)

var serviceListCmd = &cobra.Command{
	Use:   "service:list",
	Short: "List services",
	Long:  "List all registered services",
}

func init() {
	serviceListCmd.Run = func(cmd *cobra.Command, args []string) {
		services, err := db.ListServices()
		if err != nil {
			log.Fatal("Unable to query services", "err", err)
		}

		for _, service := range services {
			service.Display()
		}
	}
}
