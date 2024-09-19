package tests

import (
	"chat-app/handlers"
	"chat-app/services"
	"chat-app/storage"
	"log"
	"testing"
	"time"
)

func TestMessageCountResponseTime(t *testing.T) {
	connStr := "user=netpo4ki dbname=postgres sslmode=disable password=19770811Ee"
	store, err := storage.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	messageService := services.NewMessageService(*store)

	// Add 100 messages for the test
	for i := 0; i < 100; i++ {
		messageService.AddMessage("User", "Test message")
	}

	startTime := time.Now()
	handlers.GetMessageCount(messageService)
	duration := time.Since(startTime)

	if duration > 100*time.Millisecond {
		t.Errorf("Response took too long: %v", duration)
	}
}
