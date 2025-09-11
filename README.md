# SecureChat

A real-time chat application built with Go that supports both text and audio messaging with end-to-end privacy controls.

## Features

- **Real-time messaging** using WebSockets
- **Audio message support** with server-side file storage
- **Private messaging** - only sender and recipient can access messages
- **Authentication and authorization** system
- **Secure file access** for audio messages

## Tech Stack

- **Backend**: Go with WebSocket support
- **Frontend**: Vanilla JavaScript + Go templates
- **Storage**: File system for audio messages
- **Real-time**: WebSocket connections

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd securechat
```

2. Install dependencies:
```bash
go mod tidy
```

3. Require env:
```bash
cp .env.example .env
```

4. Run the application:
```bash
go run main.go
```

5. Access the application at `http://localhost:8080`


## API Endpoints

### Authentication
- `GET /login` - Login page
- `GET /register` - Registration page
- `POST /api/users/login` - Login user
- `POST /api/users/create` - Register user
- `POST /api/users/logout` - End user session

### Chat
- `GET /chat` - Main chat interface
- `WS /users/{userid}` - WebSocket connection for real-time messaging

### Audio Messages
- `GET /files/{filename}` - Retrieve audio file (authorized users only)

## WebSocket Message Format

```json
{
  "message_type": "TEXT|AUDIO",
  "content": "message content", // this is optional
  "content_data": "audio bytes", // this is optional, send this or content
  "recipient_id": "user_id", // must be integer
  "time": "2024-01-01T00:00:00Z" // current timestamp
}
```

## Security Features

- **Private messaging**: Messages only accessible to sender and recipient
- **Secure audio files**: Audio uploads restricted to authorized users
- **Session-based auth**: User authentication with session management
- **File access control**: Audio files served only to message participants

## Development

### Running in development mode:
```bash
air
```

### Building for production:
```bash
go build -o securechat
./securechat
```

## Usage

1. **Login/Register**: Access the authentication page
2. **Start chatting**: Send text messages in real-time
3. **Audio messages**: Record and send voice messages
4. **Private conversations**: All messages are private between participants

## Browser Support

- Modern browsers with WebSocket and MediaRecorder API support
- Chrome, Firefox, Safari, Edge (latest versions)

## License

TBD

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
