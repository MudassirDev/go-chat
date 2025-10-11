# Go-Chat - Chat application

## About

`go-chat` is a chat application that includes text and audio messaging.

Read the [blog post](https://mudassirdev.github.io/posts/chat-application/) where I talk about the building process.

The goal was to learn how to use websockets and make a functional real-time application.

## Requirements

- Go version 1.24.6+
- PostgreSQL

## Installation

Clone the repo:
```bash
git clone https://github.com/MudassirDev/go-chat.git
cd go-chat
```
Install dependencies:
```bash
go mod tidy
```
Copy the env file:
```bash
cp .env.example .env # Setup a postgres DB and update the url
```
Run the application:
```bash
make build && make run
```

## Features

- Real-time text messaging
- Voice message recording and playback
- User authentication (token-based)
- Message persistence
- WebSocket-based communication

## Libraries Used

- Gorilla Websockets — WebSocket implementation for Go
- PostgreSQL — Database for storing messages and user data

## Architecture

- Backend: Go with Gorilla WebSockets
- Frontend: HTML, CSS, JavaScript with MediaRecorder API
- Database: PostgreSQL
- Storage: Local file system for audio files

## Future Improvements

- Message read/unread states
- S3 integration for audio file storage
- Group chat functionality
- Typing indicators
