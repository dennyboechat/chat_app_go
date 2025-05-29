package storage

import (
	"sync"

	"github.com/dennyboechat/chat_app_go/internal/types"
)

// Storage is an interface for message storage implementations
type Storage interface {
	SaveMessage(msg *types.Message) error
	GetAllMessages() ([]*types.Message, error)
	GetMessagesByUser(userID int) ([]*types.Message, error)
	GetMessagesByKeyword(keyword string) ([]*types.Message, error)
}

// MemoryStorage implements Storage using in-memory data structures
type MemoryStorage struct {
	messages []*types.Message
	nextID   int
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new memory-based storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		messages: make([]*types.Message, 0),
		nextID:   1,
		mu:       sync.RWMutex{},
	}
}

// SaveMessage saves a message to the storage
func (s *MemoryStorage) SaveMessage(msg *types.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Assign an ID if not already set
	if msg.ID == 0 {
		msg.ID = s.nextID
		s.nextID++
	}

	s.messages = append(s.messages, msg)
	return nil
}

// GetAllMessages returns all messages
func (s *MemoryStorage) GetAllMessages() ([]*types.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]*types.Message, len(s.messages))
	copy(result, s.messages)

	return result, nil
}

// GetMessagesByUser returns messages from a specific user
func (s *MemoryStorage) GetMessagesByUser(userID int) ([]*types.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.Message

	for _, msg := range s.messages {
		if msg.UserID == userID {
			result = append(result, msg)
		}
	}

	return result, nil
}

// GetMessagesByKeyword returns messages containing a specific keyword
func (s *MemoryStorage) GetMessagesByKeyword(keyword string) ([]*types.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*types.Message

	for _, msg := range s.messages {
		if msg.ContainsKeyword(keyword) {
			result = append(result, msg)
		}
	}

	return result, nil
}
