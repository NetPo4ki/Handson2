package services

import (
	"chat-app/models"
	"chat-app/storage"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MessageService struct {
	store       storage.PostgresStore // Interface to use any storage
	connections []*websocket.Conn
	connMu      sync.Mutex
}

func NewMessageService(store storage.PostgresStore) *MessageService {
	return &MessageService{
		store:       store,
		connections: []*websocket.Conn{},
	}
}

func (s *MessageService) AddMessage(username, content string) error {
	message := models.Message{
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}
	err := s.store.AddMessage(message)
	if err != nil {
		return err
	}
	s.broadcastMessage(message)
	return nil
}

func (s *MessageService) GetMessages() ([]models.Message, error) {
	return s.store.GetMessages()
}

func (s *MessageService) GetMessageCount() (int, error) {
	return s.store.GetMessageCount()
}

// RegisterConnection adds a new WebSocket connection.
func (s *MessageService) RegisterConnection(conn *websocket.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	s.connections = append(s.connections, conn)

	// Send existing messages to the new connection
	messages, err := s.store.GetMessages()
	if err != nil {

	}
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			// Handle error
		}
	}
}

// UnregisterConnection removes a WebSocket connection.
func (s *MessageService) UnregisterConnection(conn *websocket.Conn) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for i, c := range s.connections {
		if c == conn {
			s.connections = append(s.connections[:i], s.connections[i+1:]...)
			break
		}
	}
}

// broadcastMessage sends a message to all connected clients.
func (s *MessageService) broadcastMessage(message models.Message) {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for _, conn := range s.connections {
		if err := conn.WriteJSON(message); err != nil {
			// Handle error
		}
	}
}
