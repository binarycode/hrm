package cmd

import (
	"os"

  "github.com/mdp/qrterminal"
	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/model"
)

var serviceQRCodeCmd = &cobra.Command{
	Use:   "service:qrcode SERVICE_NAME",
	Short: "Service token QR code",
	Long:  "Show service token in QR code format",
}

func init() {
	serviceQRCodeCmd.Args = cobra.ExactArgs(1)

	serviceQRCodeCmd.Run = func(cmd *cobra.Command, args []string) {
		name := args[0]

		service, err := db.GetService(model.Service{Name: name})
		if err != nil {
			log.Fatal("Cannot find service", "name", name)
		}

		qrterminal.Generate(service.Token, qrterminal.L, os.Stdout)
	}
}
