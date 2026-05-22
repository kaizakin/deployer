package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SystemStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Engine    string    `json:"engine"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SystemStatus{
		Status:    "HEALTHY",
		Timestamp: time.Now(),
		Engine:    "Go Runtime v1.26",
	})
}

func main() {
	http.HandleFunc("/api/v1/status", statusHandler)
	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Orchestration target initializing securely on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		os.Exit(1)
	}
}
