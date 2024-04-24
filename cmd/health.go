package cmd

import (
	"fmt"

	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "health [state]",
		Short: "Check the intergrity for a given state",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:   healthCmdFunc,
	}
	rootCmd.AddCommand(command)
}

func healthCmdFunc(cmd *cobra.Command, args []string) {
	stateName := args[0]

	state, exists := config.GetConfig().States[stateName]
	if exists == false {
		fmt.Printf("Undefined\nYou can create it in your configuration file.\n")

		return
	}

	if state.IsOpen() {
		fmt.Printf("State %q is Open.\n\n", stateName)
		healthOpen(&state)
	} else {
		fmt.Printf("State %q is Close.\n\n", stateName)
		healthClose(&state)
	}
}

func healthOpen(state *config.State) {
	checkArchivePresence(state.GetArchivePath())

	if disk.FileOrDirectoryExists(state.GetSymlinkTargetPath()) == false {
		fmt.Printf("Decrypted archive KO.\nDirectory not present at %v\n", state.GetSymlinkTargetPath())
	} else {
		fmt.Printf("Decrypted archive OK\n")
	}

	checkDanglingBackupPresence(state.GetArchiveBackupPath())
}

func healthClose(state *config.State) {
	checkArchivePresence(state.GetArchivePath())
	checkDanglingBackupPresence(state.GetArchiveBackupPath())
}

func checkArchivePresence(archivePath string) {
	if disk.FileOrDirectoryExists(archivePath) == false {
		fmt.Printf("Encrypted archive KO.\nArchive not present at %v\n", archivePath)
	} else {
		fmt.Printf("Encrypted archive OK\n")
	}
}

func checkDanglingBackupPresence(backupPath string) {
	if disk.FileOrDirectoryExists(backupPath) == true {
		fmt.Printf("\nWARNING: Dangling archive backup at %v\n", backupPath)
	}
}
