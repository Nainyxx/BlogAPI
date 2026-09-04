// HandlerInterface contract: HTTP handler methods registered on the router.
package interfaces

import "github.com/gin-gonic/gin"

type HandlerInterface interface {
	GetUsers() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	GetPosts() gin.HandlerFunc
	GetPost() gin.HandlerFunc
	GetComments() gin.HandlerFunc
	GetComment() gin.HandlerFunc

	RegisterUser() gin.HandlerFunc
	CreatePost() gin.HandlerFunc
	CreateComment() gin.HandlerFunc

	UpdateUser() gin.HandlerFunc
	UpdatePost() gin.HandlerFunc
	UpdateComment() gin.HandlerFunc

	DeleteUser() gin.HandlerFunc
	DeletePost() gin.HandlerFunc
	DeleteComment() gin.HandlerFunc

	LikePost() gin.HandlerFunc
	DislikePost() gin.HandlerFunc
	LikeComment() gin.HandlerFunc
	DislikeComment() gin.HandlerFunc
}
