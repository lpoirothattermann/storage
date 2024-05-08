package cmd

import (
	"fmt"

	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "status [state]",
		Short: "Get the status of a state, either open, close or undefined",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:   statusCmdFunc,
	}
	rootCmd.AddCommand(command)
}

func statusCmdFunc(cmd *cobra.Command, args []string) {
	stateName := args[0]

	if config.GetConfig().GetState(stateName).IsOpen() {
		fmt.Printf("Open\n")
	} else {
		fmt.Printf("Close\n")
	}
}
