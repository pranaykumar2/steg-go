package main

import (
	"log"
	"time"

	"github.com/pranaykumar2/steg-go/api"
	"github.com/pranaykumar2/steg-go/api/utils"
)

func main() {
	// Create required directories
	if err := utils.EnsureDirectoryExists(utils.UploadDir); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	if err := utils.EnsureDirectoryExists(utils.TempDir); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	// Set up scheduled cleanup of temp files (files older than 1 hour)
	go func() {
		for {
			if err := utils.CleanupTempFiles(1 * time.Hour); err != nil {
				log.Printf("Error cleaning up temp files: %v", err)
			}
			time.Sleep(15 * time.Minute)
		}
	}()

	// Create and start the server
	server := api.NewServer()
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
