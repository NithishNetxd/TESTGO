package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Port string `toml:"port"`
}

func main() {
    http.HandleFunc("/", timeHandler)

    var config Config
    if _, err := toml.DecodeFile("config.toml", &config); err != nil {
        log.Fatalf("Error loading config.toml: %v", err)
    }


    fmt.Printf("Server is running on port %s...\n", config.Port)
    if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format("2006-01-02 15:04:05") // Format time as YYYY-MM-DD HH:MM:SS
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "Current server time: %s\n", currentTime)
}

