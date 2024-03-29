package main

import (
	"flag"
	"fmt"
	"os"
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
	targetPtr := backupCmd.String("target", "", "Object target location (Required)")
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
		os.Exit(1)
	}

	switch os.Args[1] {
	case "backup":
		backupCmd.Parse(os.Args[2:])
		if *bucketPtr == "" || *targetPtr == "" {
			backupCmd.Usage()
			os.Exit(1)
		}
		fmt.Println("Starting backup of terraform state.")
	case "check-x":
		checkXCmd.Parse(os.Args[2:])
		fmt.Println("Performing check-x...")
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
		os.Exit(1)
	}
}
