package cli

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "huston",
	Short: "Huston is a CLI tool to invoke actions on spacelifts before_ and after_ hooks.",
}

var (
	dryRun bool
)

func init() {
	RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "To run the command in dry-run mode, default is false")
}
