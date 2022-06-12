package cmd

import (
	application "github.com/evertontomalok/distributed-system-go/internal/app"
	"github.com/evertontomalok/distributed-system-go/internal/ui/workers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runWorkerBalanceWorker = &cobra.Command{
	Use:   "validate-balance-worker",
	Short: "Run worker",
	Run: func(cmd *cobra.Command, args []string) {
		config := application.Configure()

		workers.StartValidateBalance(cmd.Context(), config)
	},
}

func init() {
	viper.BindEnv("Kafka.Port", "KAFKA_PORT")
	viper.BindEnv("Kafka.Host", "KAFKA_HOST")
	rootCmd.AddCommand(runWorkerBalanceWorker)
}
