package types

import (
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	ID        int
	UserID    int
	Username  string
	Content   string
	Timestamp time.Time
}

// NewMessage creates a new message
func NewMessage(id int, userID int, username string, content string) *Message {
	return &Message{
		ID:        id,
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// ContainsKeyword checks if message contains the given keyword
func (m *Message) ContainsKeyword(keyword string) bool {
	return strings.Contains(strings.ToLower(m.Content), strings.ToLower(keyword))
}

// FormatMessage formats a message for display
func FormatMessage(msg *Message) string {
	timeStr := msg.Timestamp.Format(time.RFC3339)
	return strings.Join([]string{"[", timeStr, "] ", msg.Username, ": ", msg.Content}, "")
}
