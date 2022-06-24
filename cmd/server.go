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
		application.ConfigureFlags()
		application.InitDB(ctx, config)
		application.InitKafka(ctx, config)
		rest.RunServer(ctx, config)
	},
}

func init() {
	_ = viper.BindEnv("Host", "HOST")
	_ = viper.BindEnv("Port", "PORT")
	_ = viper.BindEnv("Kafka.Port", "KAFKA_PORT")
	_ = viper.BindEnv("Kafka.Host", "KAFKA_HOST")
	_ = viper.BindEnv("Postgres.Host", "POSTGRES_HOST")
	rootCmd.AddCommand(runCmd)
}
