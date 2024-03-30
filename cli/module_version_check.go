package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thoughtgears/spacelift-runner/internal/log"
	"github.com/thoughtgears/spacelift-runner/internal/spacelift"
)

func init() {
	moduleVersionCheckCmd := &cobra.Command{
		Use:   "module-version",
		Short: "Check if the latest version of a module is being used",
		Run: func(cmd *cobra.Command, args []string) {
			terraformRoot, _ := cmd.Flags().GetString("terraform-root")

			data, err := getModuleVersions()
			if err != nil {
				log.Log(log.ErrorLevel, err.Error())
			}

			versionsOk, err := checkModuleVersions(terraformRoot, data)
			if err != nil {
				log.Log(log.ErrorLevel, err.Error())
			}

			if versionsOk {
				log.Log(log.InfoLevel, "Module version check complete, no versions found that do not match expected versions.")
			}
		},
	}
	moduleVersionCheckCmd.Flags().StringP("terraform-root", "t", "./", "Specifies the location of terraform root directory, default is ./")

	RootCmd.AddCommand(moduleVersionCheckCmd)
}

// getModuleVersions interacts with Spacelift to retrieve the latest versions of Terraform modules.
func getModuleVersions() (map[string]string, error) {
	client, err := spacelift.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create spacelift client: %w", err)
	}

	data, err := client.GetLatestModuleVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get module versions: %w", err)
	}

	return data, nil
}

// checkModuleVersions searches for Terraform _files in the specified root directory and
// its subdirectories, then checks their module versions against the expected versions.
func checkModuleVersions(rootDir string, data map[string]string) (bool, error) {
	terraformFilesFound := false
	versionsOk := true

	// Walk through the root directory to find Terraform _files, skipping .terraform and
	//modules directories.
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip .terraform and modules directories
		if info.IsDir() && (info.Name() == ".terraform" || info.Name() == "modules") {
			return filepath.SkipDir
		}

		// Process Terraform _files to check their module versions.
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
			// Read and process the Terraform file
			versionsOk = processTerraformFile(path, data)
			terraformFilesFound = true
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	// Warn if no Terraform _files were found in the directory.
	if !terraformFilesFound {
		log.Log(log.WarningLevel, "No Terraform _files found in the specified directory.")
	}

	return versionsOk, nil
}

// processTerraformFile reads a Terraform file, extracts module source and version,
// and checks them against expected versions.
func processTerraformFile(path string, data map[string]string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Log(log.WarningLevel, fmt.Sprintf("Error reading file: %s %v", path, err))
		return false
	}

	contentStr := string(content)

	// Define regex patterns for matching module source and version.
	// This is based on standard HCL syntax for module blocks.
	// Define regex to match each module block.
	moduleBlockRegex := regexp.MustCompile(`module\s*"[^"]+"\s*\{[^\}]+\}`)
	modules := moduleBlockRegex.FindAllString(contentStr, -1)

	versionsOk := true

	for _, module := range modules {
		// For each module, find source and version.
		sourceMatch := regexp.MustCompile(`source\s*=\s*"([^"]*spacelift\.io[^"]+)"`).FindStringSubmatch(module)
		versionMatch := regexp.MustCompile(`version\s*=\s*"([^"]+)"`).FindStringSubmatch(module)

		if sourceMatch != nil && versionMatch != nil {
			source := sourceMatch[1]
			version := versionMatch[1]

			// Check if the extracted source URL contains the module name key for any expected version.
			for key, expectedVersion := range data {
				if strings.Contains(source, key) && version != expectedVersion {
					log.Log(log.WarningLevel, fmt.Sprintf("File: %s - Module '%s' with source '%s' uses version '%s', expected '%s'.", path, key, source, version, expectedVersion))
					versionsOk = false
				}
			}
		}
	}

	return versionsOk
}
