# 📝 Blog API - Clean Architecture Project

A Blog REST API written in Go, built with Clean Architecture, interfaces, and Dependency Injection.

---

## 🎯 Architecture

```
HTTP Request
    ↓
Gin Router
    ↓
HandlerInterface (BlogHandler)
    ↓ (depends on)
BlogService Interface (BlogService)
    ↓ (uses)
Repository Interfaces
    ├─ UserRepository
    ├─ PostRepository
    └─ CommentRepository
    ↓ (implementation)
In-Memory Database or PostgreSQL
```

---

## 📦 Project structure

```
BlogAPI/
├── handlers/
│   └── blog_handler.go          ← HandlerInterface implementation
├── service/
│   └── BlogService.go           ← Business logic
├── repository/
│   ├── postgres_repository.go   ← PostgreSQL Repository implementation (default)
│   └── in_memory_repository.go  ← In-Memory Repository implementation (dev fallback)
├── db/
│   └── postgres.go              ← PostgreSQL connection setup
├── models/
│   ├── User.go
│   ├── Post.go
│   └── Comment.go
├── interfaces/
│   ├── handler_interface.go      ← HandlerInterface
│   ├── blog_service_interface.go ← BlogService Interface
│   └── repository_interfaces.go  ← Repository Interfaces
├── router/
│   └── router.go                ← Gin router setup
├── utils/
│   ├── Generate.go
│   └── Validate.go
├── test_frontend/
│   └── index.html                ← Simple manual test UI
├── main.go                      ← Entry point + DI container
├── Makefile
├── db_schema.sql                 ← PostgreSQL schema (tables, indexes, FKs)
├── .env                          ← Configuration (do not commit!)
├── go.mod
└── README.md
```

---

## 🚀 Quick start

### 1️⃣ Install Go (if not already installed)

```bash
# macOS
brew install go

# Linux
sudo apt-get install golang-go

# Windows
# Download from https://golang.org/dl/
```

### 2️⃣ Enter the project

```bash
cd BlogAPI
```

### 3️⃣ Install dependencies

```bash
make install
```

### 4️⃣ Configure .env and the database

Edit `.env` with your PostgreSQL connection details (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`) and server port. Then create the database and apply the schema:

```bash
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
make db-migrate
```

### 5️⃣ Run the server

```bash
make run
```

**Output:**
```
Server starting on http://localhost:8080
Available routes:
  GET    /users                 - Get all users
  GET    /users/:login          - Get user by login
  POST   /auth/register         - Register new user
  ...
```

---

## 🔌 API Endpoints

### Users

```bash
# Register
POST /auth/register
Content-Type: application/json

{
  "name": "John",
  "surname": "Doe",
  "login": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}

# Get all users
GET /users

# Get user by login
GET /users/johndoe
```

### Posts

```bash
# Create post
POST /posts
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "My First Post",
  "body": "This is my first blog post!",
  "image_url": "https://example.com/image.jpg"
}

# Get all posts
GET /posts

# Get post by ID
GET /posts/550e8400-e29b-41d4-a716-446655440001

# Update post
PUT /posts/550e8400-e29b-41d4-a716-446655440001
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Updated Title",
  "body": "Updated content",
  "image_url": "https://example.com/new-image.jpg"
}

# Delete post
DELETE /posts/550e8400-e29b-41d4-a716-446655440001
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Comments

```bash
# Create comment
POST /comments
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "post_id": "550e8400-e29b-41d4-a716-446655440001",
  "body": "Great post!"
}

# Get all comments
GET /comments

# Get comment by ID
GET /comments/550e8400-e29b-41d4-a716-446655440002

# Like a comment
POST /comments/550e8400-e29b-41d4-a716-446655440002/like
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}

# Unlike a comment
POST /comments/550e8400-e29b-41d4-a716-446655440002/unlike
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 🔄 Dependency Injection in main.go

```go
// 1. Connect to PostgreSQL
conn, err := db.Connect()

// 2. Create repositories (PostgreSQL-backed)
var userRepo interfaces.UserRepository = repository.NewPostgresUserRepository(conn)
var postRepo interfaces.PostRepository = repository.NewPostgresPostRepository(conn)
var commentRepo interfaces.CommentRepository = repository.NewPostgresCommentRepository(conn)

// 3. Create the service
var svc interfaces.BlogService = service.NewBlogService(userRepo, postRepo, commentRepo)

// 4. Create the handler
var handler interfaces.HandlerInterface = handlers.NewBlogHandler(svc)

// 5. Initialize the router
r := router.InitRouter(handler)

// 6. Run the server
r.Run(":8080")
```

### Switching back to in-memory storage

Useful for quick local testing without a database. In `main.go`, replace the repository constructors:

```go
// PostgreSQL:
var userRepo interfaces.UserRepository = repository.NewPostgresUserRepository(conn)

// In-memory:
var userRepo interfaces.UserRepository = repository.NewInMemoryUserRepository()
```

Same for `postRepo` and `commentRepo`. Nothing else changes — both implementations satisfy the same repository interfaces. Note that in-memory data is lost on restart.

### How saving works

Every write (`Create`, `Update`, `Delete`, `Like`, `Unlike`) executes synchronously as part of the HTTP request — the handler calls the service, which calls the repository, which runs the SQL statement and waits for it to complete before the handler responds. There is no batching or delayed flush: by the time the client gets a response, the data is already committed in PostgreSQL.

---

## 🎯 Interfaces

### HandlerInterface
- GetUsers(), GetUser(), GetPosts(), GetPost(), GetComments(), GetComment()
- RegisterUser(), CreatePost(), CreateComment()
- UpdateUser(), UpdatePost(), UpdateComment()
- DeleteUser(), DeletePost(), DeleteComment()
- LikePost(), DislikePost(), LikeComment(), DislikeComment()

### BlogService Interface
- SearchUserByID(), SearchUserByLogin(), SearchPostByID(), SearchCommentByID()
- RegisterUser(), CreatePost(), UpdatePost(), DeletePost()
- WriteComment(), LikeComment(), DislikeComment(), EditComment()

### Repository Interfaces
- **UserRepository**: Create(), GetByID(), GetByEmail(), GetByLogin(), GetAll(), Update(), Delete()
- **PostRepository**: Create(), GetByID(), GetByAuthorID(), GetAll(), Update(), Delete()
- **CommentRepository**: Create(), GetByID(), GetByPostID(), GetByAuthorID(), GetAll(), Update(), Delete()

---

## 📊 Architecture layers

| Layer | Description | Interface |
|-------|-------------|-----------|
| **HTTP** | Request handling | HandlerInterface |
| **Service** | Business logic | BlogService Interface |
| **Repository** | Data access | UserRepository, PostRepository, CommentRepository |
| **Database** | Data storage | In-Memory or PostgreSQL |

---

## ✅ Architecture benefits

✅ **Separation of Concerns** — each layer has a single responsibility
✅ **Dependency Injection** — implementations are easy to swap
✅ **Testability** — straightforward to write unit tests
✅ **Flexibility** — in-memory → PostgreSQL in one line
✅ **SOLID Principles** — all five principles are followed

---

## 🧪 Testing with Postman / cURL

To test the API, use:
- **Postman** (GUI) — import a collection
- **cURL** (terminal) — see the examples above
- The bundled `test_frontend/index.html` for a quick manual UI

---

## 🐛 Troubleshooting

**Error: "cannot find package"**
```bash
go mod tidy
go mod download
```

**Error: "bind: address already in use"**
```bash
# Change the port in .env or use a different one
SERVER_PORT=9000
```

**Error connecting to the database**
```bash
# Check that PostgreSQL is running
pg_isready -h $DB_HOST -p $DB_PORT

# Check the credentials in .env, and that the DB_USER has privileges on the tables:
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO your_db_user;
```

---

## 📝 License

MIT

---

## 👨‍💻 Author

Built as an example of Clean Architecture in Go.

---

**Happy coding! 🚀**
