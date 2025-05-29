package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dennyboechat/chat_app_go/internal/storage"
	"github.com/dennyboechat/chat_app_go/internal/types"
)

// ChatRoom manages users and messages
type ChatRoom struct {
	users       map[int]*User
	storage     storage.Storage
	broadcast   chan *types.Message
	register    chan *User
	unregister  chan *User
	userStreams map[int]chan *types.Message
	mu          sync.RWMutex
}

// NewChatRoom creates a new chat room
func NewChatRoom(storage storage.Storage) *ChatRoom {
	return &ChatRoom{
		users:       make(map[int]*User),
		storage:     storage,
		broadcast:   make(chan *Message),
		register:    make(chan *User),
		unregister:  make(chan *User),
		userStreams: make(map[int]chan *Message),
		mu:          sync.RWMutex{},
	}
}

// Start begins the chat room's message handling
func (cr *ChatRoom) Start() {
	go func() {
		for {
			select {
			case user := <-cr.register:
				cr.mu.Lock()
				cr.users[user.ID] = user
				cr.userStreams[user.ID] = make(chan *Message, 100)
				cr.mu.Unlock()
				log.Printf("User registered: %s\n", user.Username)

			case user := <-cr.unregister:
				cr.mu.Lock()
				if _, ok := cr.users[user.ID]; ok {
					delete(cr.users, user.ID)
					close(cr.userStreams[user.ID])
					delete(cr.userStreams, user.ID)
				}
				cr.mu.Unlock()
				log.Printf("User unregistered: %s\n", user.Username)

			case message := <-cr.broadcast:
				// Store message
				err := cr.storage.SaveMessage(message)
				if err != nil {
					log.Printf("Error saving message: %v\n", err)
					continue
				}

				// Broadcast to all users
				cr.mu.RLock()
				for _, stream := range cr.userStreams {
					select {
					case stream <- message:
						// Message sent successfully
					default:
						// Skip if user's channel buffer is full
						log.Println("User message buffer full, skipping message")
					}
				}
				cr.mu.RUnlock()
			}
		}
	}()
}

// RegisterUser adds a user to the chat room
func (cr *ChatRoom) RegisterUser(user *User) {
	cr.register <- user
}

// UnregisterUser removes a user from the chat room
func (cr *ChatRoom) UnregisterUser(user *User) {
	cr.unregister <- user
}

// SendMessage broadcasts a message to all users
func (cr *ChatRoom) SendMessage(userID int, content string) error {
	cr.mu.RLock()
	user, exists := cr.users[userID]
	cr.mu.RUnlock()

	if !exists {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	message := NewMessage(0, user.ID, user.Username, content)
	cr.broadcast <- message

	return nil
}

// GetUserMessages returns messages for a specific user
func (cr *ChatRoom) GetUserMessages(userID int) ([]*Message, error) {
	return cr.storage.GetMessagesByUser(userID)
}

// GetMessagesByKeyword returns messages containing a keyword
func (cr *ChatRoom) GetMessagesByKeyword(keyword string) ([]*Message, error) {
	return cr.storage.GetMessagesByKeyword(keyword)
}

// GetAllMessages returns all messages
func (cr *ChatRoom) GetAllMessages() ([]*Message, error) {
	return cr.storage.GetAllMessages()
}

// ListenForMessages returns a channel for a user to receive messages on
func (cr *ChatRoom) ListenForMessages(userID int) (<-chan *Message, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	stream, exists := cr.userStreams[userID]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not registered", userID)
	}

	return stream, nil
}

// UserExists checks if a user is in the chat room
func (cr *ChatRoom) UserExists(userID int) bool {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	_, exists := cr.users[userID]
	return exists
}

// FormatMessage formats a message for display
func FormatMessage(msg *Message) string {
	timeStr := msg.Timestamp.Format(time.RFC3339)
	return fmt.Sprintf("[%s] %s: %s", timeStr, msg.Username, msg.Content)
}
