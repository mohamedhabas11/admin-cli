package internal

import (
	"fmt"
	"log/syslog"
)

// SyslogLogger holds the syslog writer and format settings
type SyslogLogger struct {
	writer *syslog.Writer
	format string
}

// NewSyslogLogger initializes a new SyslogLogger
func NewSyslogLogger(format, protocol, remoteAddr string) *SyslogLogger {
	var writer *syslog.Writer
	var err error

	// If remoteAddr is provided, use it, otherwise use the local syslog
	if remoteAddr != "" {
		// Dial remote syslog server using the provided protocol (udp/tcp)
		writer, err = syslog.Dial(protocol, remoteAddr, syslog.LOG_INFO|syslog.LOG_DAEMON, "admin-cli")
		if err != nil {
			fmt.Printf("Failed to connect to remote syslog server: %v\n", err)
			return nil
		}
	} else {
		// Default to local syslog server
		writer, err = syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "admin-cli")
		if err != nil {
			fmt.Printf("Failed to connect to local syslog server: %v\n", err)
			return nil
		}
	}

	return &SyslogLogger{writer: writer, format: format}
}

// Log logs a message with the specified priority level
func (s *SyslogLogger) Log(priority syslog.Priority, message string) {
	// Format the message based on the desired format
	formattedMessage := s.formatMessage(priority, message)

	// Attempt to log the message
	_, err := s.writer.Write([]byte(formattedMessage))
	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
	}
}

// formatMessage formats the syslog message based on the configured format
func (s *SyslogLogger) formatMessage(priority syslog.Priority, message string) string {
	level := priorityToString(priority)

	switch s.format {
	case "RFC5424":
		return fmt.Sprintf("<%d>1 admin-cli - - - [%s] %s", priority, level, message)
	case "RFC3164":
		return fmt.Sprintf("<%d>admin-cli: [%s] %s", priority, level, message)
	default:
		return fmt.Sprintf("[%s] %s", level, message)
	}
}

// priorityToString converts syslog priority to string
func priorityToString(priority syslog.Priority) string {
	switch priority {
	case syslog.LOG_DEBUG:
		return "DEBUG"
	case syslog.LOG_INFO:
		return "INFO"
	case syslog.LOG_NOTICE:
		return "NOTICE"
	case syslog.LOG_WARNING:
		return "WARNING"
	case syslog.LOG_ERR:
		return "ERROR"
	case syslog.LOG_CRIT:
		return "CRITICAL"
	case syslog.LOG_ALERT:
		return "ALERT"
	case syslog.LOG_EMERG:
		return "EMERGENCY"
	default:
		return "INFO"
	}
}
