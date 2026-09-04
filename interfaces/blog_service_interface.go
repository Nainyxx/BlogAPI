// BlogService contract: business logic consumed by HTTP handlers.
package interfaces

import (
	"BlogAPI/models"

	"github.com/google/uuid"
)

type BlogService interface {
	GetAllUsers() ([]models.User, error)
	GetAllPosts() ([]models.Post, error)

	SearchUserByID(userID uuid.UUID) error
	SearchUserByLogin(login string) error
	SearchPostByID(postID uuid.UUID) error
	SearchCommentByID(commentID uuid.UUID) error

	IsLoginTaken(login string) bool
	IsEmailTaken(email string) bool
	UserIsExist(userID uuid.UUID) error
	PostIsExist(postID uuid.UUID) error
	CommentIsExist(commentID uuid.UUID) error

	RegisterUser(name, surname, login, email, password string) error

	CreatePost(userID uuid.UUID, title, body, imageURL string) error
	UpdatePost(userID, postID uuid.UUID, newTitle, newBody, newImageURL string) error
	DeletePost(userID, postID uuid.UUID) error

	WriteComment(userID, postID uuid.UUID, body string) error
	LikeComment(userID, commentID uuid.UUID) error
	DislikeComment(userID, commentID uuid.UUID) error
	EditComment(userID, commentID uuid.UUID, newBody string) error
}
