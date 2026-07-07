package main

import (
	"Grippy/internal/transport/http/server"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading it:", err)
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "localhost"
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT value '%s': %v", portStr, err)
	}

	srv, err := server.Init(addr, port)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
