# A simple chat application using Go language

## Features
- Users send messages
- Users receive messages
- Persist chat history
- Message filtering
- Search by user or keyword

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
