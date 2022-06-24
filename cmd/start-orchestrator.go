package cmd

import (
	application "github.com/evertontomalok/distributed-system-go/internal/app"
	event "github.com/evertontomalok/distributed-system-go/internal/core/events"
	mongoDBAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/database/mongodb"
	"github.com/evertontomalok/distributed-system-go/internal/ui/orchestrator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runOrchestrator = &cobra.Command{
	Use:   "orchestrator",
	Short: "Run orchestrator",
	Run: func(cmd *cobra.Command, args []string) {
		config := application.Configure()
		event.EventsAdapter = mongoDBAdapter.New(config)
		application.InitDB(cmd.Context(), config)
		orchestrator.StartOrchestrator(cmd.Context(), config)
	},
}

func init() {
	_ = viper.BindEnv("Kafka.Port", "KAFKA_PORT")
	_ = viper.BindEnv("Kafka.Host", "KAFKA_HOST")
	rootCmd.AddCommand(runOrchestrator)
}
