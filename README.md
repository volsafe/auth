# Authentication Service

A robust authentication service built with Go and Gin framework, providing secure user management and JWT-based authentication.

## Features

- User registration and authentication
- JWT-based authentication with access and refresh tokens
- Password hashing using bcrypt
- PostgreSQL database integration
- Swagger API documentation
- Secure password reset functionality
- User profile management

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker (optional)

## Configuration

The service uses a YAML configuration file located at `config/config.yaml`. Here's an example configuration:

```yaml
server:
  port: 8081

database:
  host: "localhost"
  port: 5432
  user: "pgsql"
  password: "pgsql"
  name: "users"

jwt:
  secret: "your-jwt-secret"
```

## Environment Variables

The following environment variables are required:

- `ACCESS_TOKEN_SECRET`: Secret key for signing access tokens
- `REFRESH_TOKEN_SECRET`: Secret key for signing refresh tokens

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd auth
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
# Using Docker
docker run --name postgres -e POSTGRES_USER=pgsql -e POSTGRES_PASSWORD=pgsql -e POSTGRES_DB=users -p 5432:5432 -d postgres:14

# Or use your existing PostgreSQL installation
```

4. Run the service:
```bash
ACCESS_TOKEN_SECRET=your-access-token-secret REFRESH_TOKEN_SECRET=your-refresh-token-secret go run cmd/main.go
```

## API Documentation

The API documentation is available through Swagger UI at `/swagger/index.html` when the service is running.

### Available Endpoints

- `POST /signup`: Create a new user account
- `POST /signin`: Authenticate and receive JWT tokens
- `PUT /forget-password`: Update user password (requires authentication)
- `POST /delete-profile`: Delete user profile (requires authentication)
- `POST /refresh`: Get new access token using refresh token

## Security

- Passwords are hashed using bcrypt
- JWT tokens for authentication
- Input validation and sanitization
- Rate limiting (TODO)
- CORS protection (TODO)

## Testing

Run the tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 