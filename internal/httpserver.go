package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func ServeHTTP(address string, port int, paths []string, logFile string) error {
	logHandler, err := setupLogHandler(logFile)
	if err != nil {
		return fmt.Errorf("error setting up log handler: %v", err)
	}

	httpIndex, err := setupHTTPIndex(paths)
	if err != nil {
		return fmt.Errorf("error setting up HTTP index: %v", err)
	}

	// create a new ServeMux for routing requests
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logHandler.ServeHTTP(w, r)
		httpIndex.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%d", address, port)
	log.Printf("starting server on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func setupLogHandler(logFile string) (http.Handler, error) {
	// Use IO Writer to handle logoutput
	var logOutput io.Writer
	if logFile != "" {
		// resolve log file path
		logFilePath, err := filepath.Abs(logFile)
		if err != nil {
			return nil, fmt.Errorf("error resolving log file path: %v", err)
		}

		// Ensure log directory exists
		if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
			return nil, fmt.Errorf("error creating log directory: %v", err)
		}

		// Open log file
		logFileHandle, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("error opening log file: %v", err)
		}

		// Set log output to open log file
		logOutput = logFileHandle
	} else {
		// Set log output to stdout if no log file is provided
		logOutput = os.Stdout
	}

	log.SetOutput(logOutput)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture the client IP address, Override if Client is behind a proxy
		clientIP := r.RemoteAddr
		if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			clientIP = ip
		}
		log.Printf("Request: %s %s from %s [User-Agent: %s]", r.Method, r.URL.Path, clientIP, r.UserAgent())
	}), nil
}

func setupHTTPIndex(paths []string) (http.Handler, error) {
	// Create a new ServeMux for routing requests
	mux := http.NewServeMux()

	// Prepare an index page for the root path
	indexPage := `<html><body><h1>Directories Served</h1><ul>`
	for _, path := range paths {
		// Resolve the absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("error resolving absolute path for %s: %v", path, err)
		}

		// Check if the path exists and is a directory
		fileInfo, err := os.Stat(absPath)
		if err != nil {
			return nil, fmt.Errorf("error accessing path %s: %v", absPath, err)
		}
		if !fileInfo.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", absPath)
		}

		// Normalize the URL path
		urlPath := "/" + filepath.Base(absPath) + "/"
		indexPage += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, urlPath, absPath)

		// Create a handler that serves the directory
		fileServer := http.FileServer(http.Dir(absPath))
		mux.Handle(urlPath, http.StripPrefix(urlPath, fileServer))

		// Log the directories being served
		log.Printf("Serving directory: %s at URL path: %s", absPath, urlPath)
	}

	// Complete the index page and serve it at '/'
	indexPage += `</ul></body></html>`
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(indexPage))
	})
	return mux, nil
}
