package Gin

import (
	"BlogAPI/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitGin инициализирует Gin router и регистрирует все routes
func InitGin(bs *service.BlogService) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// ============ PUBLIC ROUTES (GET) ============

	// Users routes
	router.GET("/users", GetUsers(bs))
	router.GET("/users/:login", GetUser(bs))

	// Posts routes
	router.GET("/posts", GetPosts(bs))
	router.GET("/posts/:postID", GetPost(bs))

	// Comments routes
	router.GET("/comments", GetComments(bs))
	router.GET("/comments/:commentID", GetComment(bs))

	// ============ AUTHENTICATION ROUTES ============

	router.POST("/auth/register", RegisterUser(bs))
	// router.POST("/auth/login", LoginUser(bs))  // TODO: добавить JWT

	// ============ POST ROUTES (CREATION) ============

	router.POST("/posts", CreatePost(bs))
	router.POST("/comments", CreateComment(bs))

	// ============ PUT ROUTES (UPDATES) ============

	router.PUT("/users/:id", UpdateUser(bs))
	router.PUT("/posts/:id", UpdatePost(bs))
	router.PUT("/comments/:id", UpdateComment(bs))

	// ============ DELETE ROUTES ============

	router.DELETE("/posts/:id", DeletePost(bs))
	router.DELETE("/comments/:id", DeleteComment(bs))

	// ============ LIKE/UNLIKE ROUTES ============

	router.POST("/posts/:id/like", LikePost(bs))
	router.POST("/posts/:id/unlike", DislikePost(bs))
	router.POST("/comments/:id/like", LikeComment(bs))
	router.POST("/comments/:id/unlike", DislikeComment(bs))

	return router
}
