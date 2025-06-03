package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Port        string `toml:"port"`
    LogLevel    string `toml:"log_level"`
    ReadTimeout int    `toml:"read_timeout"`
    WriteTimeout int   `toml:"write_timeout"`
}

// Logger wrapper for structured logging
type Logger struct {
    *log.Logger
}

func NewLogger(level string) *Logger {
    logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
    return &Logger{Logger: logger}
}

func (l *Logger) Info(format string, v ...interface{}) {
    l.Printf("[INFO] "+format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
    l.Printf("[ERROR] "+format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
    l.Printf("[DEBUG] "+format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
    l.Printf("[WARN] "+format, v...)
}

var logger *Logger

func main() {
    // Initialize logger
    logger = NewLogger("INFO")
    logger.Info("Starting time server application...")

    // Load configuration
    config, err := loadConfig("config.toml")
    if err != nil {
        logger.Error("Failed to load configuration: %v", err)
        os.Exit(1)
    }

    // Set up HTTP server with timeouts
    server := &http.Server{
        Addr:         ":" + config.Port,
        Handler:      setupRoutes(),
        ReadTimeout:  time.Duration(config.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(config.WriteTimeout) * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in a goroutine
    go func() {
        logger.Info("Server starting on port %s", config.Port)
        logger.Info("Server configuration: ReadTimeout=%ds, WriteTimeout=%ds", 
            config.ReadTimeout, config.WriteTimeout)
        
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("Server failed to start: %v", err)
            os.Exit(1)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    logger.Info("Shutting down server...")
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        logger.Error("Server forced to shutdown: %v", err)
    } else {
        logger.Info("Server gracefully stopped")
    }
}

func loadConfig(filename string) (*Config, error) {
    // Default configuration
    config := &Config{
        Port:         "8080",
        LogLevel:     "INFO",
        ReadTimeout:  10,
        WriteTimeout: 10,
    }

    // Try to load from file
    if _, err := os.Stat(filename); err == nil {
        if _, err := toml.DecodeFile(filename, config); err != nil {
            return nil, fmt.Errorf("error parsing config file: %w", err)
        }
        logger.Info("Configuration loaded from %s", filename)
    } else {
        logger.Warn("Config file %s not found, using defaults", filename)
    }

    return config, nil
}

func setupRoutes() *http.ServeMux {
    mux := http.NewServeMux()
    
    // Add logging middleware
    mux.HandleFunc("/", loggingMiddleware(timeHandler))
    mux.HandleFunc("/health", loggingMiddleware(healthHandler))
    
    return mux
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a custom ResponseWriter to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next(wrapped, r)
        
        duration := time.Since(start)
        logger.Info("Request: %s %s - Status: %d - Duration: %v - IP: %s - User-Agent: %s",
            r.Method,
            r.URL.Path,
            wrapped.statusCode,
            duration,
            getClientIP(r),
            r.UserAgent())
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    currentTime := time.Now()
    
    // Support different time formats based on query parameter
    format := r.URL.Query().Get("format")
    var timeStr string
    
    switch format {
    case "unix":
        timeStr = fmt.Sprintf("%d", currentTime.Unix())
    case "iso":
        timeStr = currentTime.Format(time.RFC3339)
    case "rfc":
        timeStr = currentTime.Format(time.RFC1123)
    default:
        timeStr = currentTime.Format("2006-01-02 15:04:05")
    }

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("Cache-Control", "no-cache")
    
    response := fmt.Sprintf("Current server time: %s\nTimezone: %s\nUTC Offset: %s\n", 
        timeStr, 
        currentTime.Location().String(),
        currentTime.Format("-07:00"))
    
    fmt.Fprint(w, response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, `{"status":"healthy","timestamp":"`+time.Now().Format(time.RFC3339)+`"}`)
}

func getClientIP(r *http.Request) string {
    // Check X-Forwarded-For header first (for load balancers/proxies)
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return forwarded
    }
    
    // Check X-Real-IP header
    realIP := r.Header.Get("X-Real-IP")
    if realIP != "" {
        return realIP
    }
    
    // Fall back to RemoteAddr
    return r.RemoteAddr
}
