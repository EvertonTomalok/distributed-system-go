package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "distributed-system-go",
	Short: "A distributed system using Kafka, postgres and mongodb, simulating a SAGA architecture.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
