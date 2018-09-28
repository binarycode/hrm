package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/binarycode/trewoga/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version",
	Long:  "Get version",
}

func init() {
	versionCmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	}
}
