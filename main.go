package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/michaelyang12/shortcuts/middleware"
)

func main() {
	port := os.Getenv("SHORTCUTS_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware.Auth(mux)))
}
