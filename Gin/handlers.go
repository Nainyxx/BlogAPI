package Gin

import (
	"BlogAPI/service"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

/*
HTTP handlers for BlogService methods
GET /users - GetUsers
*/

func GetUsers(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, bs.Users)
	}
}

func GetUser(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		login := c.Param("login")

		user, err := bs.SearchUserByLogin(login)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetPosts(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, bs.Posts)
	}
}

func GetPost(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("postID")
		id, err := uuid.Parse(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		post, err := bs.SearchPostByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusOK, post)
	}
}

func GetComments(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, bs.Comments)
	}
}

func GetComment(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("commentID")
		id, err := uuid.Parse(commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		comment, err := bs.SearchCommentByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusOK, comment)
	}
}

func RegisterUser(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name"`
			Surname  string `json:"surname"`
			Login    string `json:"login"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		err = bs.RegisterUser(input.Name, input.Surname, input.Login, input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
	}
}

func CreatePost(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			userID   string `json:"user_id"`
			title    string `json:"title"`
			body     string `json:"body"`
			imageURL string `json:"imageURL"`
		}
		err := c.ShouldBind(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user_id, err := uuid.Parse(input.userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = bs.CreatePost(user_id, input.title, input.body, input.imageURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "post created successfully"})
	}
}

func CreateComment(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			userID string `json:"user_id"`
			postID string `json:"post_id"`
			body   string `json:"body"`
		}
		err := c.ShouldBind(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user_id, err := uuid.Parse(input.userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		post_id, err := uuid.Parse(input.postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = bs.WriteComment(user_id, post_id, input.body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "post created successfully"})
	}
}

func UpdateUser(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
	}
}
