package cmd

import (
	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

var serviceRemoveCmd = &cobra.Command{
	Use:   "service:remove SERVICE_NAME",
	Short: "Remove service",
	Long:  "Remove specified service",
}

func init() {
	serviceRemoveCmd.Args = cobra.ExactArgs(1)

	serviceRemoveCmd.Run = func(cmd *cobra.Command, args []string) {
		name := args[0]

		service, err := db.GetService(model.Service{Name: name})
		if err != nil {
			log.Fatal("Cannot find service", "name", name)
		}

		if err := db.DestroyService(&service); err != nil {
			log.Fatal("Unable to destroy service", "service", service, "err", err)
		}
	}
}
