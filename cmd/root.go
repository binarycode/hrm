package cmd

import (
	"net/url"

	"github.com/mgutz/logxi/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/binarycode/trewoga/pkg/api"
	"github.com/binarycode/trewoga/pkg/bot"
	"github.com/binarycode/trewoga/pkg/db"
	"github.com/binarycode/trewoga/pkg/monitor"
	"github.com/binarycode/trewoga/pkg/server"
)

type Config struct {
	Address string
	DB      string
	Proxy   string
	Token   string
}

var rootCmd = &cobra.Command{
	Use:  "trewoga",
	Long: "Service availability monitoring server",
}

func init() {
	var (
		config     Config
		configFile string
	)

	cobra.EnableCommandSorting = false
	cobra.OnInitialize(func() {
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			viper.AddConfigPath("/etc/trewoga")
			viper.AddConfigPath("$HOME/.trewoga")
			viper.AddConfigPath(".")
			viper.SetConfigName("trewoga")
		}
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Error reading config", "err", err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal("Error parsing config", "err", err)
		}

		db.Open(config.DB)
	})

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		url, err := url.Parse(config.Address)
		if err != nil {
			log.Fatal("Unable to parse address url", "address", config.Address, "err", err)
		}
		https := url.Scheme == "https"

		go api.Start()
		go bot.Start(bot.Config{
			Address: config.Address,
			Proxy:   config.Proxy,
			Token:   config.Token,
			Webhook: https,
		})
		go monitor.Start()

		server.Start(server.Config{
			Host:  url.Hostname(),
			HTTPS: https,
		})
	}

	rootCmd.AddCommand(serviceListCmd)
	rootCmd.AddCommand(serviceAddCmd)
	rootCmd.AddCommand(serviceRemoveCmd)
	rootCmd.AddCommand(userListCmd)
	rootCmd.AddCommand(userRemoveCmd)
	rootCmd.AddCommand(userSubscriptionsCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	defer db.Close()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error running command", "err", err)
	}
}
