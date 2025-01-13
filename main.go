package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", timeHandler)

	// Start the server on port 8080
	port := "8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05") // Format time as YYYY-MM-DD HH:MM:SS
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Current server time: %s\n", currentTime)
}

