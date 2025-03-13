package http

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Server wraps the HTTP server and logger.
type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

// NewServer initializes a new server with the given configuration.
func NewServer(address string, port int, servePaths []string, uploadDir string, logFile string) (*Server, error) {
	// Set up the logger.
	logger, err := setupLogger(logFile)
	if err != nil {
		return nil, err
	}

	// Create the HTTP multiplexer.
	mux := http.NewServeMux()

	// Configure directories to serve, if provided.
	var dirs []Dir
	if len(servePaths) > 0 {
		dirs, err = setupDirs(servePaths, logger)
		if err != nil {
			return nil, err
		}

		// Parse the index template.
		indexTmpl, err := parseIndexTemplate()
		if err != nil {
			return nil, err
		}

		// Register file servers for each directory.
		for _, d := range dirs {
			registerFileServer(mux, d, logger)
		}

		// Set up the index page handler.
		mux.HandleFunc("/", createIndexHandler(indexTmpl, dirs, logger))
	}

	// If uploadDir is specified, set up the upload handler.
	if uploadDir != "" {
		// Ensure the upload directory exists.
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, fmt.Errorf("error creating upload directory: %v", err)
		}
		mux.HandleFunc("/upload", createUploadHandler(uploadDir, logger))
	}

	// Apply logging middleware.
	handler := loggingMiddleware(logger)(mux)

	// Configure the HTTP server.
	addr := fmt.Sprintf("%s:%d", address, port)
	httpSrv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return &Server{httpServer: httpSrv, logger: logger}, nil
}

// Dir represents a directory to be served.
type Dir struct {
	URL  string
	Path string
}

// setupLogger configures a logger writing to a file or stdout.
func setupLogger(logFile string) (*log.Logger, error) {
	var out io.Writer
	if logFile != "" {
		absPath, err := filepath.Abs(logFile)
		if err != nil {
			return nil, fmt.Errorf("error resolving log file path: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
			return nil, fmt.Errorf("error creating log directory: %v", err)
		}
		f, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("error opening log file: %v", err)
		}
		out = f
	} else {
		out = os.Stdout
	}
	return log.New(out, "", log.LstdFlags), nil
}

// setupDirs validates and prepares directories to serve.
func setupDirs(paths []string, logger *log.Logger) ([]Dir, error) {
	var dirs []Dir
	usedURLs := make(map[string]struct{})
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			logger.Printf("Error resolving path %s: %v", path, err)
			continue
		}
		info, err := os.Stat(absPath)
		if err != nil || !info.IsDir() {
			logger.Printf("Skipping invalid directory %s", absPath)
			continue
		}

		base := filepath.Base(absPath)
		urlPath := "/" + base + "/"
		suffix := 0
		for {
			if _, exists := usedURLs[urlPath]; !exists {
				break
			}
			suffix++
			urlPath = fmt.Sprintf("/%s-%d/", base, suffix)
		}
		usedURLs[urlPath] = struct{}{}
		dirs = append(dirs, Dir{URL: urlPath, Path: absPath})
	}
	return dirs, nil
}

// parseIndexTemplate creates and returns the index page template.
func parseIndexTemplate() (*template.Template, error) {
	tmpl, err := template.New("index").Parse(`
        <html><body>
        <h1>Directories Served</h1>
        <ul>
            {{range .}}
            <li><a href="{{.URL}}">{{.Path}}</a></li>
            {{end}}
        </ul>
        </body></html>`)
	if err != nil {
		return nil, fmt.Errorf("error creating index template: %v", err)
	}
	return tmpl, nil
}

// registerFileServer sets up a file server for a directory.
func registerFileServer(mux *http.ServeMux, d Dir, logger *log.Logger) {
	fileServer := http.FileServer(http.Dir(d.Path))
	mux.Handle(d.URL, http.StripPrefix(d.URL, fileServer))
	logger.Printf("Serving directory %s at %s", d.Path, d.URL)
}

// createIndexHandler returns a handler for the index page.
func createIndexHandler(tmpl *template.Template, dirs []Dir, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, dirs); err != nil {
			logger.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// createUploadHandler returns a handler for file uploads.
func createUploadHandler(uploadDir string, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, "Failed to read multipart form", http.StatusBadRequest)
			return
		}

		hasFiles := false
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "Failed to read part", http.StatusBadRequest)
				return
			}

			if part.FileName() == "" {
				continue
			}

			hasFiles = true
			fileName := sanitizeFileName(part.FileName())
			filePath := filepath.Join(uploadDir, fileName)

			file, err := os.Create(filePath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to create file %s", fileName), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			_, err = io.Copy(file, part)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to write file %s", fileName), http.StatusInternalServerError)
				return
			}

			logger.Printf("Saved file: %s", filePath)
		}

		if !hasFiles {
			http.Error(w, "No files provided", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Files uploaded successfully")
	}
}

// sanitizeFileName removes path components to prevent directory traversal.
func sanitizeFileName(name string) string {
	return filepath.Base(name)
}

// loggingMiddleware logs details of incoming HTTP requests.
func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
				clientIP = ip
			}
			logger.Printf("%s %s from %s [User-Agent: %s]", r.Method, r.URL.Path, clientIP, r.UserAgent())
			next.ServeHTTP(w, r)
		})
	}
}

// Start launches the server and handles graceful shutdown.
func (s *Server) Start() error {
	s.logger.Printf("Starting server on %s", s.httpServer.Addr)

	// Start the server in a goroutine.
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for shutdown signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Println("Shutdown signal received, shutting down server...")

	// Perform graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
