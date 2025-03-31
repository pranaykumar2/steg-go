package main

import (
	"log"
	"time"
	"github.com/pranaykumar2/steg-go/api"
	"github.com/pranaykumar2/steg-go/api/utils"
)

func main() {
	if err := utils.EnsureDirectoryExists(utils.UploadDir); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	if err := utils.EnsureDirectoryExists(utils.TempDir); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	go func() {
		for {
			if err := utils.CleanupTempFiles(1 * time.Hour); err != nil {
				log.Printf("Error cleaning up temp files: %v", err)
			}
			time.Sleep(15 * time.Minute)
		}
	}()

	server := api.NewServer()
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
