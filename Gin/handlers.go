package Gin

import (
	"BlogAPI/service"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

/*
HTTP handlers for BlogService methods

GET handlers (public):
- GET /users - GetUsers
- GET /users/:login - GetUser
- GET /posts - GetPosts
- GET /posts/:postID - GetPost
- GET /comments - GetComments
- GET /comments/:commentID - GetComment

POST handlers:
- POST /auth/register - RegisterUser
- POST /posts - CreatePost
- POST /comments - CreateComment

PUT handlers:
- PUT /users/:id - UpdateUser
- PUT /posts/:id - UpdatePost
- PUT /comments/:id - UpdateComment

DELETE handlers:
- DELETE /posts/:id - DeletePost
- DELETE /comments/:id - DeleteComment

Like handlers:
- POST /posts/:id/like - LikePost
- POST /posts/:id/unlike - DislikePost
- POST /comments/:id/like - LikeComment
- POST /comments/:id/unlike - DislikeComment
*/

// ============ GET HANDLERS (PUBLIC) ============

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

// ============ POST HANDLERS (REGISTRATION & CREATION) ============

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
			UserID   string `json:"user_id"`
			Title    string `json:"title"`
			Body     string `json:"body"`
			ImageURL string `json:"image_url"`
		}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user_id, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}
		err = bs.CreatePost(user_id, input.Title, input.Body, input.ImageURL)
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
			UserID string `json:"user_id"`
			PostID string `json:"post_id"`
			Body   string `json:"body"`
		}
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user_id, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}
		post_id, err := uuid.Parse(input.PostID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID format"})
			return
		}
		err = bs.WriteComment(user_id, post_id, input.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "comment created successfully"})
	}
}

// ============ PUT HANDLERS (UPDATES) ============

func UpdateUser(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		userID, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		user, err := bs.SearchUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		var input struct {
			Name    string `json:"name"`
			Surname string `json:"surname"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// ✅ Update fields
		if input.Name != "" {
			user.Name = input.Name
		}
		if input.Surname != "" {
			user.Surname = input.Surname
		}

		c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "user": user})
	}
}

func UpdatePost(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID format"})
			return
		}

		var input struct {
			UserID   string `json:"user_id"`
			Title    string `json:"title"`
			Body     string `json:"body"`
			ImageURL string `json:"image_url"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		err = bs.UpdatePost(userID, postID, input.Title, input.Body, input.ImageURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "post updated successfully"})
	}
}

func UpdateComment(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
			Body   string `json:"body"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		err = bs.EditComment(userID, commentID, input.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment updated successfully"})
	}
}

// ============ DELETE HANDLERS ============

func DeletePost(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		err = bs.DeletePost(userID, postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
	}
}

func DeleteComment(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		comment, err := bs.SearchCommentByID(commentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}

		// ✅ Check authorization
		if comment.AuthorID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to delete this comment"})
			return
		}

		// ✅ Find and delete from slice
		for i := range bs.Comments {
			if bs.Comments[i].ID == commentID {
				bs.Comments[i] = bs.Comments[len(bs.Comments)-1]
				bs.Comments = bs.Comments[:len(bs.Comments)-1]
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
	}
}

// ============ LIKE/UNLIKE HANDLERS ============

func LikePost(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		// ✅ Check if post exists
		post, err := bs.SearchPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		// ✅ Check if already liked
		for _, id := range post.UsersWhoLiked {
			if id == userID {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user has already liked this post"})
				return
			}
		}

		// ✅ Add like
		user, err := bs.SearchUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		user.LikePost(postID)
		post.UsersWhoLiked = append(post.UsersWhoLiked, userID)
		post.LikedCount++

		c.JSON(http.StatusOK, gin.H{"message": "post liked successfully"})
	}
}

func DislikePost(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("id")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		// ✅ Check if post exists
		post, err := bs.SearchPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		// ✅ Find and remove like
		likeIndex := -1
		for i, id := range post.UsersWhoLiked {
			if id == userID {
				likeIndex = i
				break
			}
		}

		if likeIndex == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user has not liked this post"})
			return
		}

		// ✅ Remove from post likes
		post.UsersWhoLiked[likeIndex] = post.UsersWhoLiked[len(post.UsersWhoLiked)-1]
		post.UsersWhoLiked = post.UsersWhoLiked[:len(post.UsersWhoLiked)-1]
		post.LikedCount--

		// ✅ Remove from user likes
		user, err := bs.SearchUserByID(userID)
		if err == nil {
			user.UnlikePost(postID)
		}

		c.JSON(http.StatusOK, gin.H{"message": "post unliked successfully"})
	}
}

func LikeComment(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		err = bs.LikeComment(userID, commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment liked successfully"})
	}
}

func DislikeComment(bs *service.BlogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		commentIDStr := c.Param("id")
		commentID, err := uuid.Parse(commentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID format"})
			return
		}

		var input struct {
			UserID string `json:"user_id"`
		}
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		err = bs.DislikeComment(userID, commentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "comment unliked successfully"})
	}
}
