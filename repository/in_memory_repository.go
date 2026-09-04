// In-memory implementation of the repository interfaces (used before a real database is wired in).
package repository

import (
	"BlogAPI/interfaces"
	"BlogAPI/models"
	"errors"
	"github.com/google/uuid"
)

type InMemoryUserRepository struct {
	users []models.User
}

func NewInMemoryUserRepository() interfaces.UserRepository {
	return &InMemoryUserRepository{
		users: make([]models.User, 0),
	}
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.users = append(r.users, *user)
	return nil
}

func (r *InMemoryUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	for i := range r.users {
		if r.users[i].ID == id {
			return &r.users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*models.User, error) {
	for i := range r.users {
		if r.users[i].Email == email {
			return &r.users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) GetByLogin(login string) (*models.User, error) {
	for i := range r.users {
		if r.users[i].Login == login {
			return &r.users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *InMemoryUserRepository) GetAll() ([]models.User, error) {
	return r.users, nil
}

func (r *InMemoryUserRepository) Update(user *models.User) error {
	for i := range r.users {
		if r.users[i].ID == user.ID {
			r.users[i] = *user
			return nil
		}
	}
	return errors.New("user not found")
}

func (r *InMemoryUserRepository) Delete(id uuid.UUID) error {
	for i := range r.users {
		if r.users[i].ID == id {
			r.users[i] = r.users[len(r.users)-1]
			r.users = r.users[:len(r.users)-1]
			return nil
		}
	}
	return errors.New("user not found")
}

type InMemoryPostRepository struct {
	posts []models.Post
}

func NewInMemoryPostRepository() interfaces.PostRepository {
	return &InMemoryPostRepository{
		posts: make([]models.Post, 0),
	}
}

func (r *InMemoryPostRepository) Create(post *models.Post) error {
	r.posts = append(r.posts, *post)
	return nil
}

func (r *InMemoryPostRepository) GetByID(id uuid.UUID) (*models.Post, error) {
	for i := range r.posts {
		if r.posts[i].ID == id {
			return &r.posts[i], nil
		}
	}
	return nil, errors.New("post not found")
}

func (r *InMemoryPostRepository) GetByAuthorID(authorID uuid.UUID) ([]models.Post, error) {
	var result []models.Post
	for i := range r.posts {
		if r.posts[i].AuthorID == authorID {
			result = append(result, r.posts[i])
		}
	}
	return result, nil
}

func (r *InMemoryPostRepository) GetAll() ([]models.Post, error) {
	return r.posts, nil
}

func (r *InMemoryPostRepository) Update(post *models.Post) error {
	for i := range r.posts {
		if r.posts[i].ID == post.ID {
			r.posts[i] = *post
			return nil
		}
	}
	return errors.New("post not found")
}

func (r *InMemoryPostRepository) Delete(id uuid.UUID) error {
	for i := range r.posts {
		if r.posts[i].ID == id {
			r.posts[i] = r.posts[len(r.posts)-1]
			r.posts = r.posts[:len(r.posts)-1]
			return nil
		}
	}
	return errors.New("post not found")
}

type InMemoryCommentRepository struct {
	comments []models.Comment
}

func NewInMemoryCommentRepository() interfaces.CommentRepository {
	return &InMemoryCommentRepository{
		comments: make([]models.Comment, 0),
	}
}

func (r *InMemoryCommentRepository) Create(comment *models.Comment) error {
	r.comments = append(r.comments, *comment)
	return nil
}

func (r *InMemoryCommentRepository) GetByID(id uuid.UUID) (*models.Comment, error) {
	for i := range r.comments {
		if r.comments[i].ID == id {
			return &r.comments[i], nil
		}
	}
	return nil, errors.New("comment not found")
}

func (r *InMemoryCommentRepository) GetByPostID(postID uuid.UUID) ([]models.Comment, error) {
	var result []models.Comment
	for i := range r.comments {
		if r.comments[i].PostID == postID {
			result = append(result, r.comments[i])
		}
	}
	return result, nil
}

func (r *InMemoryCommentRepository) GetByAuthorID(authorID uuid.UUID) ([]models.Comment, error) {
	var result []models.Comment
	for i := range r.comments {
		if r.comments[i].AuthorID == authorID {
			result = append(result, r.comments[i])
		}
	}
	return result, nil
}

func (r *InMemoryCommentRepository) GetAll() ([]models.Comment, error) {
	return r.comments, nil
}

func (r *InMemoryCommentRepository) Update(comment *models.Comment) error {
	for i := range r.comments {
		if r.comments[i].ID == comment.ID {
			r.comments[i] = *comment
			return nil
		}
	}
	return errors.New("comment not found")
}

func (r *InMemoryCommentRepository) Delete(id uuid.UUID) error {
	for i := range r.comments {
		if r.comments[i].ID == id {
			r.comments[i] = r.comments[len(r.comments)-1]
			r.comments = r.comments[:len(r.comments)-1]
			return nil
		}
	}
	return errors.New("comment not found")
}

func (r *InMemoryCommentRepository) Like(commentID, userID uuid.UUID) error {
	for i := range r.comments {
		if r.comments[i].ID == commentID {
			return r.comments[i].LikeComment(userID)
		}
	}
	return errors.New("comment not found")
}

func (r *InMemoryCommentRepository) Unlike(commentID, userID uuid.UUID) error {
	for i := range r.comments {
		if r.comments[i].ID == commentID {
			return r.comments[i].UnlikeComment(userID)
		}
	}
	return errors.New("comment not found")
}
