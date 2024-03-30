package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Level defines the type for log levels.
type Level string

const (
	// InfoLevel represents informational messages.
	InfoLevel Level = "INFO"
	// WarningLevel represents warning messages.
	WarningLevel Level = "WARNING"
	// ErrorLevel represents error messages.
	ErrorLevel Level = "ERROR"
)

// Log outputs a colorized log message based on the log level.
// If the log level is ErrorLevel, the program will exit with status code 1.
func Log(level Level, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	switch level {
	case InfoLevel:
		fmt.Print("[")
		color.New(color.FgCyan).Print(level)
		fmt.Printf("] %s\n", message)
	case WarningLevel:
		fmt.Print("[")
		color.New(color.FgYellow).Print(level)
		fmt.Printf("] %s\n", message)
	case ErrorLevel:
		fmt.Print("[")
		color.New(color.FgRed).Print(level)
		fmt.Printf("] %s\n", message)
		os.Exit(1)
	default:
		fmt.Printf("[UNKNOWN] %s\n", message)
	}
}
