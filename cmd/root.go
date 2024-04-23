package cmd

import (
	"fmt"
	"log"

	"github.com/lpoirothattermann/storage/internal/config"
	logInternal "github.com/lpoirothattermann/storage/internal/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "storage",
	Short: fmt.Sprintf(
		`
Storage manages encryption states of your directories using Age encryption (github.com/FiloSottile/age).
Config file: %q
		`,
		config.GetConfigFilePath(),
	),
}

func init() {}

func Execute() {
	if err := logInternal.Initialization(config.GetConfig().LogFilePath); err != nil {
		log.Fatalf("Error while initializing log: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		logInternal.Critical.Fatalf("Error: %v\n", err)
	}

	if err := logInternal.Close(); err != nil {
		logInternal.Critical.Fatalf("Error while closing log file: %v", err)
	}
}
