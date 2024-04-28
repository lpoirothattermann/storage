package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/lpoirothattermann/storage/internal/constants"
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
	commit := "Commit not available"
	// ok is true only when application has been build
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
		}
	}

	fmt.Printf("%v (%v)\n", constants.VERSION, commit)

}
