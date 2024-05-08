package cmd

import (
	"fmt"
	"os"

	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/lpoirothattermann/storage/internal/log"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "health [state]",
		Short: "Check the intergrity for a given state",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:   healthCmdFunc,
	}
	command.Flags().Bool("fix", false, "Try to fix unhealthy state")

	rootCmd.AddCommand(command)
}

func healthCmdFunc(cmd *cobra.Command, args []string) {
	stateName := args[0]
	shouldFix, err := cmd.Flags().GetBool("fix")
	if err != nil {
		log.Critical.Fatalf("Error while parsing fix flag: %v\n", err)
	}

	state := config.GetConfig().GetState(stateName)

	if state.IsOpen() {
		fmt.Printf("State %q is Open.\n\n", stateName)
		healthOpen(state, shouldFix)
	} else {
		fmt.Printf("State %q is Close.\n\n", stateName)
		healthClose(state, shouldFix)
	}
}

func healthOpen(state *config.State, shouldFix bool) {
	checkArchivePresence(state.GetArchivePath())

	if disk.FileOrDirectoryExists(state.GetSymlinkTargetPath()) {
		if shouldFix {
			if err := disk.CreateSymlink(state.GetTemporaryDirectoryPath(), state.GetSymlinkTargetPath(), true); err != nil {
				log.Critical.Fatalf("Error trying to create symlink: %v\n", err)
			}
			fmt.Printf("Decrypted archive OK - Fixed\n")
		} else {
			fmt.Printf("Decrypted archive KO (fixable automatically with --fix). Directory not present at %v\n", state.GetSymlinkTargetPath())
		}
	} else {
		fmt.Printf("Decrypted archive OK\n")
	}

	checkDanglingBackupPresence(state.GetArchiveBackupPath(), shouldFix)
}

func healthClose(state *config.State, shouldFix bool) {
	checkArchivePresence(state.GetArchivePath())
	checkDanglingBackupPresence(state.GetArchiveBackupPath(), shouldFix)
}

func checkArchivePresence(archivePath string) {
	if !disk.FileOrDirectoryExists(archivePath) {
		fmt.Printf("Encrypted archive KO (not fixable automatically). Archive not present at %v\n", archivePath)
	} else {
		fmt.Printf("Encrypted archive OK\n")
	}
}

func checkDanglingBackupPresence(backupPath string, shouldFix bool) {
	if disk.FileOrDirectoryExists(backupPath) {
		if shouldFix {
			if err := os.RemoveAll(backupPath); err != nil {
				log.Critical.Fatalf("Error while deleting dangling backup archive: %v\n", err)
			}
			fmt.Printf("INFO: Dangling archive backup at %v has been deleted.\n", backupPath)
		} else {
			fmt.Printf("\nWARNING: Dangling archive backup at %v .\nCan be deleted with automatically with --fix.\n", backupPath)
		}
	}
}
