// Router: registers all HTTP routes on the Gin engine and wires the CORS middleware.
package router

import (
	"BlogAPI/interfaces"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(handler interfaces.HandlerInterface) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:5500", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/users", handler.GetUsers())
	router.GET("/users/:login", handler.GetUser())

	router.GET("/posts", handler.GetPosts())
	router.GET("/posts/:postID", handler.GetPost())

	router.GET("/comments", handler.GetComments())
	router.GET("/comments/:commentID", handler.GetComment())

	router.POST("/auth/register", handler.RegisterUser())

	router.POST("/posts", handler.CreatePost())
	router.POST("/comments", handler.CreateComment())

	router.PUT("/users/:id", handler.UpdateUser())
	router.PUT("/posts/:id", handler.UpdatePost())
	router.PUT("/comments/:id", handler.UpdateComment())

	router.DELETE("/posts/:id", handler.DeletePost())
	router.DELETE("/comments/:id", handler.DeleteComment())

	router.POST("/posts/:id/like", handler.LikePost())
	router.POST("/posts/:id/unlike", handler.DislikePost())
	router.POST("/comments/:id/like", handler.LikeComment())
	router.POST("/comments/:id/unlike", handler.DislikeComment())

	return router
}
