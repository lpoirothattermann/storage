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

	state, exists := config.GetConfig().States[stateName]
	if exists == false {
		fmt.Printf("Undefined\nYou can create it in your configuration file.\n")

		return
	}

	if state.IsOpen() == true {
		fmt.Printf("Open\n")
	} else {
		fmt.Printf("Close\n")
	}

}
