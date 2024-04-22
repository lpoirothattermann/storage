package cmd

import (
	"fmt"
	"os"

	"github.com/lpoirothattermann/storage/internal/config"
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
