# A simple chat application using Go language

## Features  
- Users send messages
- Users receive messages
- Persist chat history
- Message filtering
- Search by user or keyword

## Technical Details
- In-memory database to store chat history with timestamp and user ID
- Built-in concurrency features with goroutines
- Channels for message passing
- Focus on performance

## Project Structure
```
chat_app_go/
├── cmd/
│   └── chat/
│       └── main.go           # Main application entry point
├── internal/
│   ├── models/
│   │   ├── user.go           # User model
│   │   ├── message.go        # Message model
│   │   └── chatroom.go       # ChatRoom implementation with channels
│   └── storage/
│       └── storage.go        # Message storage interface and implementation
├── go.mod                    # Go module file
└── README.md                 # This file
```

## How to Run
1. Ensure Go is installed on your system
2. Clone the repository:
   ```
   git clone https://github.com/dennyboechat/chat_app_go.git
   cd chat_app_go
   ```
3. Run the application:
   ```
   go run cmd/chat/main.go
   ```

## Usage
The application provides a simple command-line interface with the following options:
1. Send a message - Send a message as one of the simulated users
2. View all messages - Display all messages in the chat history
3. Filter messages by user - Show messages from a specific user
4. Search messages by keyword - Find messages containing specific text
5. Exit - Quit the application

## Design Considerations
- **Concurrency**: The application uses goroutines and channels to handle message broadcasting efficiently
- **Thread Safety**: All shared resources are protected with mutexes
- **Interfaces**: The storage layer is defined by interfaces for easy swapping of implementations
- **Simulated Users**: The application comes with simulated users for testing

## GitHub repository  
- https://github.com/dennyboechat/chat_app_go

## Timeline  
- 24-May (Sat) – Send deliverable 1
- 07-Jun (Sat) – Send deliverable 2
- 21-Jun (Sat) - Send deliverable 3
