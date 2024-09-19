package tests

import (
	"chat-app/services"
	"chat-app/storage"
	"log"
	"testing"
)

func TestRecoverability(t *testing.T) {
	connStr := "user=netpo4ki dbname=postgres sslmode=disable password=19770811Ee"
	store, err := storage.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	messageService := services.NewMessageService(*store)

	// Simulate adding a message
	messageService.AddMessage("User", "Test message before crash")

	// Simulate crash by reinitializing the service and store
	store, err = storage.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	messageService = services.NewMessageService(*store)

	// Check if messages are lost
	messages, err := messageService.GetMessages()
	if len(messages) != 0 {
		t.Errorf("Messages were not recovered after crash. Expected 0, got %d", len(messages))
	}
}
