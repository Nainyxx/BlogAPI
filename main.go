// Application entry point: loads config, wires dependencies, and starts the HTTP server.
package main

import (
	"BlogAPI/db"
	"BlogAPI/handlers"
	"BlogAPI/interfaces"
	"BlogAPI/repository"
	"BlogAPI/router"
	"BlogAPI/service"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var userRepo interfaces.UserRepository = repository.NewPostgresUserRepository(conn)
	var postRepo interfaces.PostRepository = repository.NewPostgresPostRepository(conn)
	var commentRepo interfaces.CommentRepository = repository.NewPostgresCommentRepository(conn)

	var svc interfaces.BlogService = service.NewBlogService(userRepo, postRepo, commentRepo)

	var handler interfaces.HandlerInterface = handlers.NewBlogHandler(svc)

	r := router.InitRouter(handler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("Available routes:")
	log.Printf("  GET    /users                 - Get all users")
	log.Printf("  GET    /users/:login          - Get user by login")
	log.Printf("  POST   /auth/register         - Register new user")
	log.Printf("  GET    /posts                 - Get all posts")
	log.Printf("  GET    /posts/:postID         - Get post by ID")
	log.Printf("  POST   /posts                 - Create post")
	log.Printf("  PUT    /posts/:id             - Update post")
	log.Printf("  DELETE /posts/:id             - Delete post")
	log.Printf("  GET    /comments              - Get all comments")
	log.Printf("  GET    /comments/:commentID   - Get comment by ID")
	log.Printf("  POST   /comments              - Create comment")
	log.Printf("  PUT    /comments/:id          - Update comment")
	log.Printf("  DELETE /comments/:id          - Delete comment")
	log.Printf("  POST   /comments/:id/like     - Like comment")
	log.Printf("  POST   /comments/:id/unlike   - Unlike comment")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
