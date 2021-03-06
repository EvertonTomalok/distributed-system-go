package cmd

import (
	application "github.com/evertontomalok/distributed-system-go/internal/app"
	userapi "github.com/evertontomalok/distributed-system-go/internal/infra/services/user-api"
	"github.com/evertontomalok/distributed-system-go/internal/ui/workers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runWorkerUserStatusWorker = &cobra.Command{
	Use:   "validate-user-status-worker",
	Short: "Run Worker",
	Run: func(cmd *cobra.Command, args []string) {
		config := application.Configure()
		userapi.UserAdapter = userapi.New(config)

		workers.StartValidateUserStatus(cmd.Context(), config)
	},
}

func init() {
	_ = viper.BindEnv("Kafka.Port", "KAFKA_PORT")
	_ = viper.BindEnv("Kafka.Host", "KAFKA_HOST")
	_ = viper.BindEnv("Mongodb.Host", "MONGODB_HOST")
	_ = viper.BindEnv("UserApi.BaseUrl", "USERAPI_BASEURL")
	rootCmd.AddCommand(runWorkerUserStatusWorker)
}
