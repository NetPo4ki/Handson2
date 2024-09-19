package main

import (
	"chat-app/handlers"
	"chat-app/services"
	"chat-app/storage"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}

func main() {
	r := gin.Default()

	// Enable CORS for all origins
	r.Use(cors.Default())

	// Initialize PostgreSQL storage
	connStr := "user=netpo4ki dbname=postgres sslmode=disable password=19770811Ee"
	store, err := storage.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Initialize the service
	messageService := services.NewMessageService(*store)

	// Set up routes
	r.POST("/messages", handlers.PostMessage(messageService))
	r.GET("/messages", handlers.GetMessages(messageService))
	r.GET("/messages/count", handlers.GetMessageCount(messageService))

	// WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		serveWs(messageService, c.Writer, c.Request)
	})

	// Run the server
	r.Run(":8080")
}

func serveWs(service *services.MessageService, w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// Register this connection in the message service
	service.RegisterConnection(conn)

	// Listen for messages to avoid broken pipe errors
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Remove connection when an error occurs (e.g., client disconnect)
			service.UnregisterConnection(conn)
			break
		}
	}
}
