package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thoughtgears/spacelift-runner/internal/backup"
)

func main() {
	// Top-level command
	hustonCmd := flag.NewFlagSet("huston", flag.ExitOnError)
	hustonCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", hustonCmd.Name())
		hustonCmd.PrintDefaults()
		fmt.Println("\nSubcommands:")
		fmt.Println("  backup: Perform a backup operation")
		fmt.Println("  check-x: Perform check-x operation")
	}

	// Subcommand: Backup
	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	bucketPtr := backupCmd.String("bucket", "", "Name of the bucket (Required)")
	stackPtr := backupCmd.String("stack", "", "Name of the stack running (Required)")
	filePathPtr := backupCmd.String("file-path", "", "File path to back up (Required)")
	backupCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", backupCmd.Name())
		backupCmd.PrintDefaults()
	}

	// Subcommand: check-x
	checkXCmd := flag.NewFlagSet("check-x", flag.ExitOnError)
	checkXCmd.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", checkXCmd.Name())
		checkXCmd.PrintDefaults()
	}

	if len(os.Args) < 2 {
		hustonCmd.Usage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "backup":
		backupCmd.Parse(os.Args[2:])
		if *bucketPtr == "" || *filePathPtr == "" || *stackPtr == "" {
			backupCmd.Usage()
			os.Exit(0)
		}
		if err := backup.Run(*bucketPtr, *stackPtr, *filePathPtr); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "check-x":
		checkXCmd.Parse(os.Args[2:])
	case "help":
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "backup":
				backupCmd.Usage()
			case "check-x":
				checkXCmd.Usage()
			default:
				hustonCmd.Usage()
			}
		} else {
			hustonCmd.Usage()
		}
	default:
		hustonCmd.Usage()
	}
}
