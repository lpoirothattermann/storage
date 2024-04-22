package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print Storage version",
		Run:   versionCmdFunc,
	})
}

func versionCmdFunc(cmd *cobra.Command, args []string) {
	fmt.Println("0.0")
}
