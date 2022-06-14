package cmd

import (
	application "github.com/evertontomalok/distributed-system-go/internal/app"
	event "github.com/evertontomalok/distributed-system-go/internal/domain/events"
	mongoDBAdapter "github.com/evertontomalok/distributed-system-go/internal/infra/mongodb"
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
		orchestrator.StartOrchestrator(cmd.Context(), config)
	},
}

func init() {
	viper.BindEnv("Kafka.Port", "KAFKA_PORT")
	viper.BindEnv("Kafka.Host", "KAFKA_HOST")
	rootCmd.AddCommand(runOrchestrator)
}
