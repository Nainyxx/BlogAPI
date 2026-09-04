// BlogService: business logic layer implementing interfaces.BlogService on top of the repository interfaces.
package service

import (
	"BlogAPI/interfaces"
	"BlogAPI/models"
	"errors"

	"github.com/google/uuid"
)

type BlogService struct {
	userRepo    interfaces.UserRepository
	postRepo    interfaces.PostRepository
	commentRepo interfaces.CommentRepository
}

func NewBlogService(userRepo interfaces.UserRepository, postRepo interfaces.PostRepository, commentRepo interfaces.CommentRepository) *BlogService {
	return &BlogService{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (s *BlogService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *BlogService) GetAllPosts() ([]models.Post, error) {
	return s.postRepo.GetAll()
}

func (s *BlogService) SearchUserByID(userID uuid.UUID) error {
	_, err := s.userRepo.GetByID(userID)
	return err
}

func (s *BlogService) SearchUserByLogin(login string) error {
	_, err := s.userRepo.GetByLogin(login)
	return err
}

func (s *BlogService) SearchPostByID(postID uuid.UUID) error {
	_, err := s.postRepo.GetByID(postID)
	return err
}

func (s *BlogService) SearchCommentByID(commentID uuid.UUID) error {
	_, err := s.commentRepo.GetByID(commentID)
	return err
}

func (s *BlogService) IsLoginTaken(login string) bool {
	_, err := s.userRepo.GetByLogin(login)
	return err == nil
}

func (s *BlogService) IsEmailTaken(email string) bool {
	_, err := s.userRepo.GetByEmail(email)
	return err == nil
}

func (s *BlogService) UserIsExist(userID uuid.UUID) error {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return errors.New("user does not exist")
	}
	return nil
}

func (s *BlogService) PostIsExist(postID uuid.UUID) error {
	if _, err := s.postRepo.GetByID(postID); err != nil {
		return errors.New("post does not exist")
	}
	return nil
}

func (s *BlogService) CommentIsExist(commentID uuid.UUID) error {
	if _, err := s.commentRepo.GetByID(commentID); err != nil {
		return errors.New("comment does not exist")
	}
	return nil
}

func (s *BlogService) RegisterUser(name, surname, login, email, password string) error {
	if s.IsLoginTaken(login) {
		return errors.New("login already taken")
	}
	if s.IsEmailTaken(email) {
		return errors.New("email already taken")
	}

	user, err := models.CreateUser(name, surname, login, email, password)
	if err != nil {
		return err
	}

	return s.userRepo.Create(&user)
}

func (s *BlogService) CreatePost(userID uuid.UUID, title, body, imageURL string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	post, err := models.CreatePost(userID, title, body, imageURL)
	if err != nil {
		return err
	}

	if err := s.postRepo.Create(&post); err != nil {
		return err
	}

	user.AddPost(post.ID)
	return s.userRepo.Update(user)
}

func (s *BlogService) UpdatePost(userID, postID uuid.UUID, newTitle, newBody, newImageURL string) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	if newTitle != "" {
		if err := post.UpdateTitle(newTitle, userID); err != nil {
			return err
		}
	}
	if newBody != "" {
		if err := post.UpdateBody(newBody, userID); err != nil {
			return err
		}
	}
	if newImageURL != "" {
		if err := post.UpdateImageURL(newImageURL, userID); err != nil {
			return err
		}
	}

	return s.postRepo.Update(post)
}

func (s *BlogService) DeletePost(userID, postID uuid.UUID) error {
	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return errors.New("not authorized")
	}
	return s.postRepo.Delete(postID)
}

func (s *BlogService) WriteComment(userID, postID uuid.UUID, body string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	post, err := s.postRepo.GetByID(postID)
	if err != nil {
		return err
	}

	comment, err := models.CreateComment(userID, postID, body)
	if err != nil {
		return err
	}

	if err := s.commentRepo.Create(&comment); err != nil {
		return err
	}

	post.AddComment(comment.ID)
	if err := s.postRepo.Update(post); err != nil {
		return err
	}

	user.AddComment(comment.ID)
	return s.userRepo.Update(user)
}

func (s *BlogService) LikeComment(userID, commentID uuid.UUID) error {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return err
	}
	if _, err := s.commentRepo.GetByID(commentID); err != nil {
		return err
	}
	return s.commentRepo.Like(commentID, userID)
}

func (s *BlogService) DislikeComment(userID, commentID uuid.UUID) error {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return err
	}
	if _, err := s.commentRepo.GetByID(commentID); err != nil {
		return err
	}
	return s.commentRepo.Unlike(commentID, userID)
}

func (s *BlogService) EditComment(userID, commentID uuid.UUID, newBody string) error {
	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return err
	}
	if err := comment.ChangeBody(userID, newBody); err != nil {
		return err
	}
	return s.commentRepo.Update(comment)
}
