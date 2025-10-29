package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

// Message represents a chat message
type Message struct {
	Username  string
	Content   string
	Timestamp time.Time
}

// SendMessageArgs represents the arguments for sending a message
type SendMessageArgs struct {
	Username string
	Content  string
}

// printMessage prints a single message
func printMessage(msg Message) {
	fmt.Printf("%s: %s\n", msg.Username, msg.Content)
}

// pollForUpdates periodically checks for new messages
func pollForUpdates(client *rpc.Client, username string, lastCount *int) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var history []Message
		err := client.Call("ChatService.GetHistory", username, &history)
		if err != nil {
			continue
		}

		// If there are new messages, print only the new ones
		if len(history) > *lastCount {
			fmt.Print("\r\033[K") // Clear the current line, for proper printing
			for i := *lastCount; i < len(history); i++ {
				printMessage(history[i])
			}
			*lastCount = len(history)
			fmt.Printf("%s: ", username)
		}
	}
}

func main() {
	// Connect to server
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}

	fmt.Println("Connected to server\n")

	// Get username using bufio.Reader to read full line
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading username:", err)
	}
	username = strings.TrimSpace(username)

	if username == "" {
		log.Fatal("Username cannot be empty")
	}

	fmt.Println("Type 'exit' to quit.\n")

	// Get initial history
	var history []Message
	err = client.Call("ChatService.GetHistory", username, &history)
	if err != nil {
		log.Println("Error fetching history:", err)
	} else {
		// Display existing messages
		for _, msg := range history {
			printMessage(msg)
		}
		if len(history) > 0 {
			fmt.Println()
		}
	}

	lastMessageCount := len(history)

	// Start polling for updates in background
	go pollForUpdates(client, username, &lastMessageCount)

	// Main loop - runs forever until "exit" or Ctrl+C
	for {
		fmt.Printf("%s: ", username)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		message = strings.TrimSpace(message)

		// Check for exit command
		if message == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Skip empty messages
		if message == "" {
			continue
		}

		// Send message to server
		args := SendMessageArgs{
			Username: username,
			Content:  message,
		}
		var reply []Message

		err = client.Call("ChatService.SendMessage", args, &reply)
		if err != nil {
			log.Println("Error sending message:", err)
			fmt.Println("Server may be down.")
			continue
		}

		// Update message count
		lastMessageCount = len(reply)
	}
}
