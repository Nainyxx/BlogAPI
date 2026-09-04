// BlogHandler: HTTP handlers implementing interfaces.HandlerInterface, delegating to the BlogService interface.
package handlers

import (
	"net/http"

	"BlogAPI/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogHandler struct {
	service interfaces.BlogService
}

func NewBlogHandler(service interfaces.BlogService) interfaces.HandlerInterface {
	return &BlogHandler{
		service: service,
	}
}

func (h *BlogHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.service.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func (h *BlogHandler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		login := c.Param("login")

		err := h.service.SearchUserByLogin(login)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User found", "login": login})
	}
}

func (h *BlogHandler) GetPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := h.service.GetAllPosts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, posts)
	}
}

func (h *BlogHandler) GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		err = h.service.SearchPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post found", "id": postID})
	}
}

func (h *BlogHandler) GetComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get all comments"})
	}
}

func (h *BlogHandler) GetComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("commentID")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		err = h.service.SearchCommentByID(commentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Comment found", "id": commentID})
	}
}

func (h *BlogHandler) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Surname  string `json:"surname" binding:"required"`
			Login    string `json:"login" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := h.service.RegisterUser(input.Name, input.Surname, input.Login, input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

func (h *BlogHandler) CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			UserID   string `json:"user_id" binding:"required"`
			Title    string `json:"title" binding:"required"`
			Body     string `json:"body" binding:"required"`
			ImageURL string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.CreatePost(userID, input.Title, input.Body, input.ImageURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
	}
}

func (h *BlogHandler) CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			UserID string `json:"user_id" binding:"required"`
			PostID string `json:"post_id" binding:"required"`
			Body   string `json:"body" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		postID, err := uuid.Parse(input.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		err = h.service.WriteComment(userID, postID, input.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully"})
	}
}

func (h *BlogHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Update user"})
	}
}

func (h *BlogHandler) UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		var input struct {
			UserID   string `json:"user_id" binding:"required"`
			Title    string `json:"title"`
			Body     string `json:"body"`
			ImageURL string `json:"image_url"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.UpdatePost(userID, postID, input.Title, input.Body, input.ImageURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
	}
}

func (h *BlogHandler) UpdateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		var input struct {
			UserID string `json:"user_id" binding:"required"`
			Body   string `json:"body" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.EditComment(userID, commentID, input.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
	}
}

func (h *BlogHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Delete user"})
	}
}

func (h *BlogHandler) DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		var input struct {
			UserID string `json:"user_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.DeletePost(userID, postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}

func (h *BlogHandler) DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Delete comment"})
	}
}

func (h *BlogHandler) LikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Like post"})
	}
}

func (h *BlogHandler) DislikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Unlike post"})
	}
}

func (h *BlogHandler) LikeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		var input struct {
			UserID string `json:"user_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.LikeComment(userID, commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comment liked successfully"})
	}
}

func (h *BlogHandler) DislikeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		var input struct {
			UserID string `json:"user_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.DislikeComment(userID, commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comment unliked successfully"})
	}
}
