# VicNotes Backend

A high-performance Go backend for the VicNotes note-taking application.

## Features

- User authentication with JWT tokens
- CRUD operations for notes
- PostgreSQL database with proper indexing
- RESTful API design
- Error handling and validation
- Password hashing with bcrypt
- Request logging and recovery middleware

## Project Structure

```
backend/
├── main.go                 # Application entry point
├── config/                 # Configuration management
├── database/              # Database initialization and migrations
├── models/                # Data models
├── handlers/              # HTTP request handlers
├── middleware/            # HTTP middleware
├── utils/                 # Utility functions (JWT, password hashing)
├── go.mod                 # Go module definition
├── Dockerfile             # Docker containerization
└── README.md              # This file
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker and Docker Compose (optional)

## Setup

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables in `.env`:
```
POSTGRES_USER=vicnotes-test-user
POSTGRES_PASSWORD=vicnotes-test-password
POSTGRES_DB=vicnotes-test-db
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
```

3. Start PostgreSQL (using Docker Compose from root):
```bash
cd ..
docker-compose up -d
```

4. Run the backend:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Docker

Build and run with Docker:
```bash
docker build -t vicnotes-backend .
docker run -p 8080:8080 --env-file ../.env vicnotes-backend
```

## API Endpoints

### Health Check
- `GET /health` - Check server status

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user

### Notes (Protected - requires JWT token)
- `POST /api/v1/notes` - Create a new note
- `GET /api/v1/notes` - List all user's notes
- `GET /api/v1/notes/{id}` - Get a specific note
- `PUT /api/v1/notes/{id}` - Update a note
- `DELETE /api/v1/notes/{id}` - Delete a note

## Authentication

Protected endpoints require an `Authorization` header with a Bearer token:
```
Authorization: Bearer <token>
```

Tokens are obtained from the login or register endpoints and are valid for 24 hours.

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Notes Table
```sql
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Example Requests

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Create Note
```bash
curl -X POST http://localhost:8080/api/v1/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"My Note","content":"Note content here"}'
```

### List Notes
```bash
curl -X GET http://localhost:8080/api/v1/notes \
  -H "Authorization: Bearer <token>"
```

## Best Practices Implemented

- **Security**: Password hashing with bcrypt, JWT authentication
- **Error Handling**: Comprehensive error responses with meaningful messages
- **Database**: Connection pooling, proper indexing, prepared statements
- **Code Organization**: Clear separation of concerns with packages
- **Middleware**: Logging, recovery, and authentication middleware
- **Validation**: Input validation on all endpoints
- **Performance**: Database indexes on frequently queried columns
- **Scalability**: Stateless design suitable for horizontal scaling

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o vicnotes-backend
```

## License

See LICENSE file in the root directory.
