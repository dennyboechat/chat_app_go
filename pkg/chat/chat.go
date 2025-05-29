package chat

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// User represents a chat participant
type User struct {
	ID       int
	Username string
}

// Message represents a chat message
type Message struct {
	ID        int
	UserID    int
	Username  string
	Content   string
	Timestamp time.Time
}

// ChatRoom manages users, messages, and connections
type ChatRoom struct {
	users       map[int]*User
	messages    []*Message
	broadcast   chan *Message
	register    chan *User
	unregister  chan *User
	userStreams map[int]chan *Message
	nextMsgID   int
	mu          sync.RWMutex
}

// NewChatRoom creates a new chat room
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		users:       make(map[int]*User),
		messages:    make([]*Message, 0),
		broadcast:   make(chan *Message),
		register:    make(chan *User),
		unregister:  make(chan *User),
		userStreams: make(map[int]chan *Message),
		nextMsgID:   1,
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
				cr.saveMessage(message)

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

	message := &Message{
		ID:        0, // Will be set by saveMessage
		UserID:    user.ID,
		Username:  user.Username,
		Content:   content,
		Timestamp: time.Now(),
	}
	
	cr.broadcast <- message
	return nil
}

// Internal method to save a message
func (cr *ChatRoom) saveMessage(msg *Message) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	
	// Assign an ID if not already set
	if msg.ID == 0 {
		msg.ID = cr.nextMsgID
		cr.nextMsgID++
	}
	
	cr.messages = append(cr.messages, msg)
}

// GetUserMessages returns messages from a specific user
func (cr *ChatRoom) GetUserMessages(userID int) []*Message {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	
	var result []*Message
	
	for _, msg := range cr.messages {
		if msg.UserID == userID {
			result = append(result, msg)
		}
	}
	
	return result
}

// GetMessagesByKeyword returns messages containing a keyword
func (cr *ChatRoom) GetMessagesByKeyword(keyword string) []*Message {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	
	var result []*Message
	
	for _, msg := range cr.messages {
		if containsKeyword(msg, keyword) {
			result = append(result, msg)
		}
	}
	
	return result
}

// Helper function to check if a message contains a keyword
func containsKeyword(msg *Message, keyword string) bool {
	return strings.Contains(strings.ToLower(msg.Content), strings.ToLower(keyword))
}

// GetAllMessages returns all messages
func (cr *ChatRoom) GetAllMessages() []*Message {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	result := make([]*Message, len(cr.messages))
	copy(result, cr.messages)
	
	return result
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

// FormatMessage formats a message for display
func FormatMessage(msg *Message) string {
	timeStr := msg.Timestamp.Format(time.RFC3339)
	return fmt.Sprintf("[%s] %s: %s", timeStr, msg.Username, msg.Content)
}

// NewUser creates a new user
func NewUser(id int, username string) *User {
	return &User{
		ID:       id,
		Username: username,
	}
}
