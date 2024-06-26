package cmd

import (
	"os"

	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/disk"
	"github.com/lpoirothattermann/storage/internal/log"
	"github.com/spf13/cobra"
)

func init() {
	command := &cobra.Command{
		Use:   "open [state]",
		Short: "Open state",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:   openCmdFunc,
	}
	command.Flags().BoolP("force", "f", false, "Force state opening, lose data if state is not closed")

	rootCmd.AddCommand(command)
}

func openCmdFunc(cmd *cobra.Command, args []string) {
	forceCommand, _ := cmd.Flags().GetBool("force")
	stateName := args[0]

	state := config.GetConfig().GetState(stateName)

	file, err := os.Open(state.GetArchivePath())
	if err != nil {
		log.Critical.Fatalf("Error while opening archive: %v\n", err)
	}

	bundleReader, err := bundler.NewReader(file, state.GetIdentity())
	if err != nil {
		log.Critical.Fatalf("Error while opening bundle reader: %v\n", err)
	}

	tmpDirectoryPath := state.GetTemporaryDirectoryPath()
	if disk.FileOrDirectoryExists(tmpDirectoryPath) {
		if forceCommand {
			if err := os.RemoveAll(tmpDirectoryPath); err != nil {
				log.Critical.Fatalf("Error while removing old temporary directory: %v\n", err)
			}
		} else {
			log.Critical.Fatalf("State is already open.")
		}
	}

	if err := disk.UnbundleToPath(bundleReader, tmpDirectoryPath); err != nil {
		log.Critical.Fatalf("Error while writing archive on disk: %v\n", err)
	}

	if err := disk.CreateSymlink(tmpDirectoryPath, state.GetSymlinkTargetPath(), true); err != nil {
		log.Critical.Fatalf("Error while creating symlink: %v\n", err)
	}
}
