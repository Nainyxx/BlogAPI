package models

import (
	"BlogAPI/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Login        string    `json:"login"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`

	Posts      []uuid.UUID `json:"posts"`
	LikedPosts []uuid.UUID `json:"liked_posts"`

	CommentsID      []uuid.UUID `json:"commentsID"`
	LikedCommentsID []uuid.UUID `json:"liked_commentsID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUser(Name, Surname, Login, Email, Password string) (User, error) {
	if utils.IsNameValid(Name) != nil {
		return User{}, errors.New("Name is not valid")
	} else if utils.IsNameValid(Surname) != nil {
		return User{}, errors.New("Surname is not valid")
	} else if utils.IsLoginValid(Login) != nil {
		return User{}, errors.New("Login is not valid")
	} else if utils.IsEmailValid(Email) != nil {
		return User{}, errors.New("Email is not valid")
	}

	passwordHash, err := utils.HashPassword(Password)
	if err != nil {
		return User{}, err
	}

	newUser := User{
		ID:              uuid.New(),
		Name:            Name,
		Surname:         Surname,
		Login:           Login,
		Email:           Email,
		PasswordHash:    passwordHash,
		Posts:           make([]uuid.UUID, 0),
		LikedPosts:      make([]uuid.UUID, 0),
		CommentsID:      make([]uuid.UUID, 0),
		LikedCommentsID: make([]uuid.UUID, 0),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return newUser, nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) AddPost(postID uuid.UUID) {
	for _, post := range u.Posts {
		if post == postID {
			return
		}
	}
	u.Posts = append(u.Posts, postID)
	u.UpdatedAt = time.Now()
}

func (u *User) RemovePost(postID uuid.UUID) {
	for i, post := range u.Posts {
		if post == postID {
			u.Posts[i] = u.Posts[len(u.Posts)-1]
			u.Posts = u.Posts[:len(u.Posts)-1]

			u.UpdatedAt = time.Now()
			return
		}
	}
}

func (u *User) LikePost(postID uuid.UUID) {
	for _, post := range u.LikedPosts {
		if post == postID {
			return
		}
	}
	u.LikedPosts = append(u.LikedPosts, postID)
	u.UpdatedAt = time.Now()
}

func (u *User) UnlikePost(postID uuid.UUID) {
	for i, id := range u.LikedPosts {
		if id == postID {
			u.LikedPosts[i] = u.LikedPosts[len(u.LikedPosts)-1]
			u.LikedPosts = u.LikedPosts[:len(u.LikedPosts)-1]

			u.UpdatedAt = time.Now()
			return
		}
	}
}

func (u *User) LikeComment(commentID uuid.UUID) {
	for _, id := range u.LikedCommentsID {
		if id == commentID {
			return
		}
	}
	u.LikedCommentsID = append(u.LikedCommentsID, commentID)
	u.UpdatedAt = time.Now()
}

func (u *User) UnlikeComment(commentID uuid.UUID) {
	for i, id := range u.LikedCommentsID {
		if id == commentID {
			u.LikedCommentsID[i] = u.LikedCommentsID[len(u.LikedCommentsID)-1]
			u.LikedCommentsID = u.LikedCommentsID[:len(u.LikedCommentsID)-1]

			u.UpdatedAt = time.Now()
			return
		}
	}
}
