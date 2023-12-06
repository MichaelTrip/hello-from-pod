package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func getHostname() (string, error) {
	cmd := exec.Command("hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, err := getHostname()
	if err != nil {
		http.Error(w, "Error getting hostname", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("Hello from Pod %s", hostname)
	fmt.Fprintln(w, message)

	// Log the client's visit
	log.Printf("Client %s visited the server\n", r.RemoteAddr)
}

func main() {
	// Create a log file for more persistent logging if needed
	logFile, err := os.Create("server.log")
	if err != nil {
		fmt.Println("Error creating log file:", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Set up logging to both stdout and the log file
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Create a custom server with keepalive initially enabled
	server := &http.Server{
		Addr:           ":8080",
		Handler:        http.HandlerFunc(handler),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    0,       // Disable keepalive
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Disable keepalives explicitly
	server.SetKeepAlivesEnabled(false)

	fmt.Println("Server is running on :8080...")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
}
