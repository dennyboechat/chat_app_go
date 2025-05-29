package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/dennyboechat/chat_app_go/pkg/chat"
)

func main() {
	// Create a new chat room
	chatRoom := chat.NewChatRoom()
	
	// Start the chat room
	chatRoom.Start()
	// Create some users
	users := []*chat.User{
		chat.NewUser(1, "Alice"),
		chat.NewUser(2, "Bob"),
		chat.NewUser(3, "Charlie"),
	}
	
	// Register all users
	for _, user := range users {
		chatRoom.RegisterUser(user)
	}

	// Set up simulated user message senders
	var wg sync.WaitGroup

	// Start the listener for each user
	for _, user := range users {
		wg.Add(1)
		go func(u *chat.User) {
			defer wg.Done()

			stream, err := chatRoom.ListenForMessages(u.ID)
			if err != nil {
				fmt.Printf("Error getting message stream for %s: %v\n", u.Username, err)
				return
			}

			// Listen for messages
			for msg := range stream {
				formattedMsg := chat.FormatMessage(msg)
				fmt.Printf("[%s received] %s\n", u.Username, formattedMsg)
			}
		}(user)
	}

	// Start command line interface for interaction
	go commandLineInterface(chatRoom, users)

	// Wait for all goroutines to finish (this won't happen in reality as the CLI keeps running)
	wg.Wait()
}

func commandLineInterface(chatRoom *chat.ChatRoom, users []*chat.User) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nChat Application Commands:")
		fmt.Println("1. Send a message")
		fmt.Println("2. View all messages")
		fmt.Println("3. Filter messages by user")
		fmt.Println("4. Search messages by keyword")
		fmt.Println("5. Exit")
		fmt.Print("\nEnter your choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			sendMessage(chatRoom, users, scanner)
		case "2":
			viewAllMessages(chatRoom)
		case "3":
			filterMessagesByUser(chatRoom, users, scanner)
		case "4":
			searchMessagesByKeyword(chatRoom, scanner)
		case "5":
			fmt.Println("Exiting chat application.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func sendMessage(chatRoom *chat.ChatRoom, users []*chat.User, scanner *bufio.Scanner) {
	// Display available users
	fmt.Println("Available users:")
	for _, user := range users {
		fmt.Printf("%d. %s\n", user.ID, user.Username)
	}

	// Get user ID
	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := scanner.Text()
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Println("Invalid user ID. Please enter a number.")
		return
	}

	// Check if user exists
	userExists := false
	for _, user := range users {
		if user.ID == userID {
			userExists = true
			break
		}
	}

	if !userExists {
		fmt.Println("User with that ID doesn't exist.")
		return
	}

	// Get message content
	fmt.Print("Enter your message: ")
	scanner.Scan()
	content := scanner.Text()

	// Send the message
	err = chatRoom.SendMessage(userID, content)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	fmt.Println("Message sent successfully.")
}

func viewAllMessages(chatRoom *chat.ChatRoom) {
	messages := chatRoom.GetAllMessages()

	if len(messages) == 0 {
		fmt.Println("No messages found.")
		return
	}

	fmt.Println("All Messages:")
	for _, msg := range messages {
		fmt.Println(chat.FormatMessage(msg))
	}
}

func filterMessagesByUser(chatRoom *chat.ChatRoom, users []*chat.User, scanner *bufio.Scanner) {
	// Display available users
	fmt.Println("Available users:")
	for _, user := range users {
		fmt.Printf("%d. %s\n", user.ID, user.Username)
	}

	// Get user ID
	fmt.Print("Enter user ID: ")
	scanner.Scan()
	userIDStr := scanner.Text()
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		fmt.Println("Invalid user ID. Please enter a number.")
		return
	}

	// Check if user exists
	userExists := false
	for _, user := range users {
		if user.ID == userID {
			userExists = true
			break
		}
	}

	if !userExists {
		fmt.Println("User with that ID doesn't exist.")
		return
	}

	// Get messages by user
	messages := chatRoom.GetUserMessages(userID)

	if len(messages) == 0 {
		fmt.Println("No messages found for this user.")
		return
	}

	fmt.Printf("Messages from user ID %d:\n", userID)
	for _, msg := range messages {
		fmt.Println(chat.FormatMessage(msg))
	}
}

func searchMessagesByKeyword(chatRoom *chat.ChatRoom, scanner *bufio.Scanner) {
	// Get keyword
	fmt.Print("Enter keyword to search for: ")
	scanner.Scan()
	keyword := scanner.Text()

	if keyword == "" {
		fmt.Println("Keyword cannot be empty.")
		return
	}

	// Search messages by keyword
	messages := chatRoom.GetMessagesByKeyword(keyword)

	if len(messages) == 0 {
		fmt.Printf("No messages found containing '%s'.\n", keyword)
		return
	}

	fmt.Printf("Messages containing '%s':\n", keyword)
	for _, msg := range messages {
		fmt.Println(chat.FormatMessage(msg))
	}
}
