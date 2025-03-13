package cmd

import (
	"admin-cli/http"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	httpAddress   string
	httpPort      int
	httpPaths     []string
	httpUploadDir string
	httplogFile   string
)

var httpServeCmd = &cobra.Command{
	Use:   "httpserver",
	Short: "Serve a list of paths over HTTP, defaulting to the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a new HTTP server instance
		server, err := http.NewServer(httpAddress, httpPort, httpPaths, httpUploadDir, httplogFile)
		if err != nil {
			return fmt.Errorf("error starting HTTP server: %v", err)
		}
		// Start the server
		if err := server.Start(); err != nil {
			return fmt.Errorf("error running HTTP server: %v", err)
		}
		return nil
	},
}

func init() {
	// Define flags for the httpserver command
	httpServeCmd.Flags().StringVar(&httpAddress, "address", "127.0.0.1", "Address to listen on")
	httpServeCmd.Flags().IntVar(&httpPort, "port", 8080, "Port to listen on")
	httpServeCmd.Flags().StringSliceVarP(&httpPaths, "path", "P", []string{"."}, "Paths to serve")
	httpServeCmd.Flags().StringVarP(&httpUploadDir, "upload-dir", "u", "", "Directory for file uploads")
	httpServeCmd.Flags().StringVar(&httplogFile, "log", "", "Log file to write to")

	// Add the serve command to the root command
	rootCmd.AddCommand(httpServeCmd)
}
