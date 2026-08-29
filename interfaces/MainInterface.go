package interfaces

import "github.com/google/uuid"

type BlogService interface {
	CreateUser(name, surname, login, email, password string) error
	CreatePost(userID uuid.UUID, title, body, imageURL string) error
	CreateComment(userID, postID uuid.UUID, body string) error

	LikePost(userID, postID uuid.UUID) error
	UnlikePost(userID, postID uuid.UUID) error
	LikeComment(userID, commentID uuid.UUID) error
	UnlikeComment(userID, commentID uuid.UUID) error

	EditUser(userID uuid.UUID, new_name, surname string)
	EditPost(userID, postID uuid.UUID, title, new_body, imageURL string) error

	DeleteUser(userID uuid.UUID) error
	DeletePost(userID uuid.UUID) error
}
