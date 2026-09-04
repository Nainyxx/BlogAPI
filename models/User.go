// User model: domain data and behavior for accounts (registration, posts/comments tracking, likes).
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

	PostsID      []uuid.UUID `json:"posts"`
	LikedPostsID []uuid.UUID `json:"liked_posts"`

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
		PostsID:         make([]uuid.UUID, 0),
		LikedPostsID:    make([]uuid.UUID, 0),
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
	for _, post := range u.PostsID {
		if post == postID {
			return
		}
	}
	u.PostsID = append(u.PostsID, postID)
	u.UpdatedAt = time.Now()
}

func (u *User) AddComment(commentID uuid.UUID) {
	for _, comment := range u.CommentsID {
		if comment == commentID {
			return
		}
	}
	u.CommentsID = append(u.CommentsID, commentID)
	u.UpdatedAt = time.Now()
}

func (u *User) RemovePost(postID uuid.UUID) {
	for i, post := range u.PostsID {
		if post == postID {
			u.PostsID[i] = u.PostsID[len(u.PostsID)-1]
			u.PostsID = u.PostsID[:len(u.PostsID)-1]

			u.UpdatedAt = time.Now()
			return
		}
	}
}

func (u *User) LikePost(postID uuid.UUID) {
	for _, post := range u.LikedPostsID {
		if post == postID {
			return
		}
	}
	u.LikedPostsID = append(u.LikedPostsID, postID)
	u.UpdatedAt = time.Now()
}

func (u *User) UnlikePost(postID uuid.UUID) {
	for i, id := range u.LikedPostsID {
		if id == postID {
			u.LikedPostsID[i] = u.LikedPostsID[len(u.LikedPostsID)-1]
			u.LikedPostsID = u.LikedPostsID[:len(u.LikedPostsID)-1]

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

func (u *User) CreatePost(postID uuid.UUID) {
	for _, id := range u.PostsID {
		if id == postID {
			return
		}
	}
	u.PostsID = append(u.PostsID, postID)
	u.UpdatedAt = time.Now()
}

func (u *User) UpdateName(newName string) error {
	if u.Name != newName {
		if utils.IsNameValid(newName) == nil {
			u.Name = newName
			u.UpdatedAt = time.Now()
			return nil
		} else {
			return errors.New("name is not valid")
		}
	}
	return errors.New("new name must be different")
}

func (u *User) UpdateSurname(newName string) error {
	if u.Surname != newName {
		if utils.IsNameValid(newName) == nil {
			u.Surname = newName
			u.UpdatedAt = time.Now()
			return nil
		} else {
			return errors.New("surname is not valid")
		}
	}
	return errors.New("surnew name must be different")
}
