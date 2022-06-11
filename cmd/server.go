package cmd

import (
	application "github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/ui/rest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "Run http server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		config := application.Configure()

		application.InitDB(ctx, config)
		application.InitKafka(ctx, config)
		rest.RunServer(ctx, config)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	viper.BindEnv("HOST")
	viper.BindEnv("PORT")
}
