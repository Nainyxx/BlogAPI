package main

import (
	"BlogAPI/my_models/Gin"
	"BlogAPI/my_models/service"
	"log"

	"github.com/gin-contrib/cors"
)

func main() {
	blog_service := service.NewBlogService()
	router := Gin.InitGin(blog_service)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:5500"}, // Укажите точный адрес фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	log.Println("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
