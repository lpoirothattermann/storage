package cmd

import (
	"log"
	"os"

	"github.com/lpoirothattermann/storage/internal/bundler"
	"github.com/lpoirothattermann/storage/internal/config"
	"github.com/lpoirothattermann/storage/internal/disk"
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
	forceCommand, err := cmd.Flags().GetBool("force")
	stateName := args[0]

	state, exists := config.GetConfig().States[stateName]
	if exists == false {
		log.Fatalf("State %q doesn't exists.\n'", stateName)
	}

	file, err := os.Open(state.GetArchivePath())
	if err != nil {
		log.Fatalf("Error while opening encrypted archive: %v\n", err)
	}

	bundleReader, err := bundler.NewReader(file, state.AgeIdentity)
	if err != nil {
		log.Fatalf("Error while opening bundle reader: %v\n", err)
	}

	tmpDirectoryPath := state.GetTemporaryDirectoryPath()
	if disk.FileOrDirectoryExists(tmpDirectoryPath) == true {
		if forceCommand == true {
			if err := os.RemoveAll(tmpDirectoryPath); err != nil {
				log.Fatalf("Error while removing old temporary directory: %v\n", err)
			}
		} else {
			log.Fatalf("State is already open.")
		}
	}

	disk.WriteBundleToPath(bundleReader, tmpDirectoryPath)
	disk.CreateSymlink(tmpDirectoryPath, state.GetSymlinkTargetPath(), true)
}
