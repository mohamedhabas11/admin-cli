package cmd

import (
	"admin-cli/internal"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var (
	address net.IP
	port    int
	paths   []string
	logFile string
)

var serveCmd = &cobra.Command{
	Use:   "httpserver",
	Short: "Serve a list of paths over HTTP, defaulting to the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Call the function from the internal package to start the server
		err := internal.ServeHTTP(address.String(), port, paths, logFile)
		if err != nil {
			return fmt.Errorf("error starting HTTP server: %v", err)
		}
		return nil
	},
}

func init() {
	// Define flags for the httpserver command
	serveCmd.Flags().IPVar(&address, "address", net.ParseIP("127.0.0.1"), "Address to listen on")
	serveCmd.Flags().IntVar(&port, "port", 8080, "Port to listen on")
	serveCmd.Flags().StringVar(&logFile, "log", "", "Log file to write to")
	serveCmd.Flags().StringSliceVarP(&paths, "path", "P", []string{"."}, "Paths to serve")

	// Add the serve command to the root command
	rootCmd.AddCommand(serveCmd)
}
