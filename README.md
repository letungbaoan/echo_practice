# Echo Practice - RealWorld API Implementation

A learning project implementing the [RealWorld](https://realworld.io/) API specification using **Go + Echo + GORM + PostgreSQL**.

## 🎯 Project Overview

This is an educational project designed to learn backend development with Go and the Echo web framework. The implementation follows the RealWorld API spec, covering:

- User authentication (JWT)
- Article management (CRUD operations)
- Comments on articles
- User profiles
- Article tagging

Progress is tracked in `TODO.md` across 13 sections, with each section building incrementally on previous work.

## 🛠️ Tech Stack

- **Language**: Go 1.26
- **Web Framework**: Echo v4
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: golang-jwt/jwt v5
- **Password Hashing**: bcrypt (golang.org/x/crypto)
- **Environment**: godotenv
- **Testing**: testify
- **Development**: Air (live-reload)

## 📋 Prerequisites

- Go 1.26+
- Docker & Docker Compose (for PostgreSQL)
- Git

## ⚙️ Setup & Installation

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd echo_practice
go mod download
```

### 2. Start PostgreSQL

```bash
docker compose up -d
```

PostgreSQL will run on **port 5433** (mapped from container's 5432).

### 3. Configure Environment

Create a `.env` file in the project root:

```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=echo_practice

# JWT
JWT_SECRET=your-secret-key-here
```

### 4. Run the Server

**Option A: With live-reload (development)**
```bash
air
```

**Option B: Direct run**
```bash
go run ./cmd/server
```

**Option C: Build and run**
```bash
go build -o server ./cmd/server
./server
```

## 📁 Project Structure

```
echo_practice/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── config/               # Environment configuration
│   ├── database/             # GORM setup & migrations
│   ├── models/               # GORM models
│   ├── dto/                  # Request/Response structures
│   ├── repositories/         # Data access layer
│   ├── services/             # Business logic
│   ├── controllers/          # HTTP handlers
│   ├── middlewares/          # JWT, error handling
│   ├── routes/               # Route registration
│   └── utils/                # Utilities (JWT, password, slug)
├── docker-compose.yml        # PostgreSQL container
├── go.mod & go.sum          # Dependencies
└── README.md                 # This file
```

## 📚 Key Conventions

### Response Format
Follows RealWorld API spec - wrapping responses in resource objects:
```json
{
  "user": {
    "username": "john",
    "email": "john@example.com",
    "image": "...",
    "bio": "..."
  }
}
```

### Error Format
```json
{
  "errors": {
    "body": ["error message 1", "error message 2"]
  }
}
```

### Database Migrations
Uses GORM `AutoMigrate` for schema management (no separate migration tool).

### Logging
Uses Go's standard `slog` with Echo's `RequestLoggerWithConfig`.

## 🧪 Building & Testing

### Verify Build
```bash
go build ./...
```

### Run Tests
```bash
go test ./...
```

## 🐛 Troubleshooting

### Port Already in Use
If port 8080 is busy, change `PORT` in `.env`.

### Database Connection Failed
- Ensure Docker is running: `docker ps`
- Check credentials in `.env` match `docker-compose.yml`
- Reset database: `docker compose down -v && docker compose up -d`

### Import/Build Errors
```bash
go clean -modcache
go mod tidy
go mod download
```

## 📮 Postman Collection

1. **Import**: Open Postman → Click **Import** → Select the Postman collection file (`.json`)
2. **Environment**: Set up a variable `base_url = http://localhost:8080/api`
3. **Authentication**: JWT token from login response is auto-saved to `{{token}}`
4. **Requests**: All endpoints are organized in folders (auth, articles, comments, profiles, tags)

## 📖 Learning Resources

- [RealWorld API Spec](https://realworld.io/docs/specs/backend-specs/introduction)
- [Echo Documentation](https://echo.labstack.com/)
- [GORM Guides](https://gorm.io/docs/)
- [Go Best Practices](https://golang.org/doc/effective_go)

## 📄 License

Learning project for educational purposes.

---

Happy learning! 🚀
