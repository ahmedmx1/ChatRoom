# Simple Chatroom

A basic chatroom application built with Go using RPC (Remote Procedure Call) for communication between clients and server.

## Features

- Multiple users can connect and chat simultaneously
- Real-time message updates (polls every 2 seconds)
- Full chat history stored on the server
- Messages appear naturally as they're sent
- Simple and clean interface

## Files

- `server.go` - RPC server that stores and manages chat messages
- `client.go` - Client application for sending and receiving messages

## How to Run

### 1. Start the Server

Open a terminal and run:

```bash
go run server.go
```

You should see:
```
Chat server running on port 1234...
```

### 2. Start Client(s)

Open one or more additional terminals and run:

```bash
go run client.go
```

Each client will:
1. Connect to the server
2. Prompt you to enter a username
3. Show any existing chat history
4. Let you start chatting

## Usage

[demo](https://drive.google.com/file/d/1IQEK-uaqqeOampRAnOv6V4HM8m9nEmMW/view?usp=drive_link)
### Sending Messages

Simply type your message and press Enter:

```
Alice: Hello everyone!
Bob: Hi Alice!
Alice: How are you?
```

### Exiting

Type `exit` and press Enter to quit the chatroom:

```
Alice: exit
Goodbye!
```

Or press `Ctrl+C` to force quit.

## Technical Details

- **Protocol**: TCP/IP with Go's net/rpc
- **Port**: 1234 (localhost)
- **Message Storage**: In-memory list on server
- **Update Frequency**: Clients poll for new messages every 2 seconds
- **Thread Safety**: Server uses mutex for concurrent access

## Requirements

- Server and clients must run on the same machine (localhost)

