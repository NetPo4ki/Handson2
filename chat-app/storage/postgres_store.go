package storage

import (
	"chat-app/models"
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check if the database is reachable
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) AddMessage(message models.Message) error {
	_, err := s.db.Exec(
		"INSERT INTO messages (username, content, timestamp) VALUES ($1, $2, $3)",
		message.Username, message.Content, message.Timestamp,
	)
	return err
}

func (s *PostgresStore) GetMessages() ([]models.Message, error) {
	rows, err := s.db.Query("SELECT username, content, timestamp FROM messages ORDER BY timestamp")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(&message.Username, &message.Content, &message.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (s *PostgresStore) GetMessageCount() (int, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	return count, err
}
