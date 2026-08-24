package interfaces

import "github.com/gin-gonic/gin"

type HandlerInterface interface {
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	GetPosts(c *gin.Context)
	GetPost(c *gin.Context)
	GetComments(c *gin.Context)
	GetComment(c *gin.Context)

	RegisterUser(c *gin.Context)
	CreatePost(c *gin.Context)
	CreateComment(c *gin.Context)

	UpdateUser(c *gin.Context)
	UpdatePost(c *gin.Context)
	UpdateComment(c *gin.Context)

	DeleteUser(c *gin.Context)
	DeletePost(c *gin.Context)
	DeleteComment(c *gin.Context)

	LikePost(c *gin.Context)
	DislikePost(c *gin.Context)
	LikeComment(c *gin.Context)
	DislikeComment(c *gin.Context)
}
