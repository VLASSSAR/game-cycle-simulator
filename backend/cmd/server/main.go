package main

import (
	"log"
	"net/http"
	"os"

	"game-cycle-simulator/internal/config"
	apphttp "game-cycle-simulator/internal/transport/http"
)

func main() {
	cfg := config.New()
	router := apphttp.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("server started on :%s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
