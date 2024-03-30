package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thoughtgears/spacelift-runner/internal/backup"
	"github.com/thoughtgears/spacelift-runner/internal/log"
)

func init() {
	backupCmd := &cobra.Command{
		Use:   "backup-state",
		Short: "Perform a backup on the state-file to a GCS bucket",
		Run: func(cmd *cobra.Command, args []string) {
			var object string
			stateFile, _ := cmd.Flags().GetString("state-file")
			stackID, _ := cmd.Flags().GetString("stack-id")
			bucket, _ := cmd.Flags().GetString("bucket")

			if !strings.HasSuffix(stateFile, ".tfstate") {
				log.Log(log.ErrorLevel, "State file must be a .tfstate file")
			}

			if strings.Contains(stateFile, "/") {
				parts := strings.Split(stateFile, "/")
				object = fmt.Sprintf("%s/%s", stackID, parts[len(parts)-1])
			} else {
				object = fmt.Sprintf("%s/%s", stackID, stateFile)
			}

			if dryRun {
				log.Log(log.InfoLevel, fmt.Sprintf("DRY-RUN: Backup of: %s to %s/%s complete", stateFile, bucket, object))
			} else {
				client, err := backup.NewClient(bucket, object)
				if err != nil {
					log.Log(log.ErrorLevel, err.Error())
				}

				if err := client.Execute(stateFile); err != nil {
					log.Log(log.ErrorLevel, err.Error())
				}
				log.Log(log.InfoLevel, fmt.Sprintf("Backup of: %s to %s/%s complete", stateFile, bucket, object))
			}
		},
	}
	backupCmd.Flags().StringP("state-file", "s", "terraform.tfstate", "Specifies the location of the state file, default is terraform.tfstate")
	backupCmd.Flags().StringP("bucket", "b", "", "Specifies the backup GCS bucket name")
	backupCmd.Flags().StringP("stack-id", "i", "", "Specifies the Spacelift stack ID")
	backupCmd.MarkFlagRequired("bucket")
	backupCmd.MarkFlagRequired("stack-id")

	RootCmd.AddCommand(backupCmd)
}
