package main

import (
	"log"
	stdhttp "net/http"

	"game-cycle-simulator/internal/config"
	apphttp "game-cycle-simulator/internal/transport/http"
)

func main() {
	cfg := config.New()

	router := apphttp.NewRouter()

	log.Printf("server started on :%s", cfg.Port)

	if err := stdhttp.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
