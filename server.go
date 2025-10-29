package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Message represents a chat message
type Message struct {
	Username  string
	Content   string
	Timestamp time.Time
}

// ChatService is our RPC service
type ChatService struct {
	mu       sync.Mutex
	messages []Message
}

// SendMessageArgs represents the arguments for sending a message
type SendMessageArgs struct {
	Username string
	Content  string
}

// SendMessage is the RPC method for sending messages
// Takes username and content, returns the full chat history
func (c *ChatService) SendMessage(args *SendMessageArgs, reply *[]Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create new message with timestamp
	msg := Message{
		Username:  args.Username,
		Content:   args.Content,
		Timestamp: time.Now(),
	}

	// Add to message list
	c.messages = append(c.messages, msg)

	// Log the message
	fmt.Printf("[%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), msg.Username, msg.Content)

	// Return full history
	*reply = make([]Message, len(c.messages))
	copy(*reply, c.messages)

	return nil
}

// GetHistory is the RPC method for fetching all messages
// Takes username (not used but kept for consistency), returns all messages
func (c *ChatService) GetHistory(username string, reply *[]Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Return full history
	*reply = make([]Message, len(c.messages))
	copy(*reply, c.messages)

	return nil
}

func main() {
	// Listen on port 1234
	listener, _ := net.Listen("tcp", "127.0.0.1:1234")

	fmt.Println("Chat server running on port 1234...")

	// Register service -> Publish the methods for clients
	rpc.Register(new(ChatService))

	// Accept connections forever
	for {
		conn, _ := listener.Accept()
		// ServeConn runs the DefaultServer on a single connection
		go rpc.ServeConn(conn)
	}
}
