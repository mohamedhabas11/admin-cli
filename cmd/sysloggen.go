package cmd

import (
	"admin-cli/internal"
	"fmt"
	"log/syslog"

	"github.com/spf13/cobra"
)

var (
	logLevel   string
	logFormat  string
	logMessage string
	remoteAddr string
	protocol   string
)

// syslogGenCmd is the command for generating syslog messages
var syslogGenCmd = &cobra.Command{
	Use:   "sysloggen",
	Short: "Generate a syslog message",
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new syslog logger with the specified format and protocol
		logger := internal.NewSyslogLogger(logFormat, protocol, remoteAddr)

		// Use the custom message if provided, otherwise use a default message
		if logMessage == "" {
			logMessage = fmt.Sprintf("Generated syslog message with level: %s", logLevel)
		}

		// Convert logLevel string to syslog.Priority
		var priority syslog.Priority
		switch logLevel {
		case "DEBUG":
			priority = syslog.LOG_DEBUG
		case "INFO":
			priority = syslog.LOG_INFO
		case "WARNING":
			priority = syslog.LOG_WARNING
		case "ERROR":
			priority = syslog.LOG_ERR
		default:
			priority = syslog.LOG_INFO // Default to INFO if an invalid level is provided
		}

		// Log the message with the specified level
		logger.Log(priority, logMessage)
		fmt.Println("Syslog message sent successfully.")
	},
}

func init() {
	// Flags for syslog generation
	syslogGenCmd.Flags().StringVarP(&logLevel, "level", "l", "INFO", "Log level (e.g., INFO, WARNING, ERROR, DEBUG)")
	syslogGenCmd.Flags().StringVarP(&logFormat, "format", "f", "RFC5424", "Syslog message format (RFC5424, RFC3164)")
	syslogGenCmd.Flags().StringVarP(&logMessage, "message", "m", "", "Custom syslog message")
	syslogGenCmd.Flags().StringVarP(&remoteAddr, "remote", "r", "", "Comma-separated list of remote syslog servers (optional)")
	syslogGenCmd.Flags().StringVarP(&protocol, "protocol", "p", "udp", "Protocol to use for syslog (udp/tcp)")

	// Add the command to rootCmd
	rootCmd.AddCommand(syslogGenCmd)
}
