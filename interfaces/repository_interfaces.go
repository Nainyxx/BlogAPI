// Repository contracts for persisting users, posts, and comments.
package interfaces

import (
	"BlogAPI/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByLogin(login string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id uuid.UUID) (*models.Post, error)
	GetByAuthorID(authorID uuid.UUID) ([]models.Post, error)
	GetAll() ([]models.Post, error)
	Update(post *models.Post) error
	Delete(id uuid.UUID) error
}

type CommentRepository interface {
	Create(comment *models.Comment) error
	GetByID(id uuid.UUID) (*models.Comment, error)
	GetByPostID(postID uuid.UUID) ([]models.Comment, error)
	GetByAuthorID(authorID uuid.UUID) ([]models.Comment, error)
	GetAll() ([]models.Comment, error)
	Update(comment *models.Comment) error
	Delete(id uuid.UUID) error
	Like(commentID, userID uuid.UUID) error
	Unlike(commentID, userID uuid.UUID) error
}
