package models

import (
	"errors"

	"github.com/google/uuid"
)

/*
GET

POST
CreateUser
CreatePost
CreateComment

PATCH
UpdateUser
UpdatePost

DELETE
DeleteUser
DeletePost
DeleteComment
*/

type BlogService struct {
	users    []User
	posts    []Post
	comments []Comment
}

func NewBlogService() *BlogService {
	return &BlogService{
		users:    make([]User, 0),
		posts:    make([]Post, 0),
		comments: make([]Comment, 0),
	}
}

func (this *BlogService) SearchUserByLogin(login string) (*User, error) {
	for _, user := range this.users {
		if user.Login == login {
			return &user, nil
		}
	}
	return nil, errors.New("user not found by login")
}

func (this *BlogService) SearchUserByEmail(email string) (*User, error) {
	for _, user := range this.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("user not found by email")
}

func (this *BlogService) SearchUserByID(id uuid.UUID) (*User, error) {
	for _, user := range this.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found by id")
}

func (this *BlogService) SearchPostByID(id uuid.UUID) (*Post, error) {
	for _, post := range this.posts {
		if post.ID == id {
			return &post, nil
		}
	}
	return nil, errors.New("post not found by id")
}

func (this *BlogService) SearchCommentByID(id uuid.UUID) (*Comment, error) {
	for _, comment := range this.comments {
		if comment.ID == id {
			return &comment, nil
		}
	}
	return nil, errors.New("comment not found by id")
}

// POST

func (this *BlogService) CreateUser(name, surname, login, email, password string) error {
	_, err := this.SearchUserByLogin(login)
	if err == nil {
		return errors.New("login already taken")
	}
	newUser, err := CreateUser(name, surname, login, email, password)
	if err != nil {
		return err
	}
	this.users = append(this.users, newUser)
	return nil
}

func (this *BlogService) CreatePost(userID uuid.UUID, title, body, imageURL string) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	newPost, err := CreatePost(userID, title, body, imageURL)
	if err != nil {
		return err
	}
	this.posts = append(this.posts, newPost)
	actualUser.CreatePost(newPost.ID)
	return nil
}

func (this *BlogService) CreateComment(userID, postID uuid.UUID, body string) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	actualPost, err := this.SearchPostByID(postID)
	if err != nil {
		return err
	}
	newComment, err := CreateComment(userID, postID, body)
	if err != nil {
		return err
	}
	this.comments = append(this.comments, newComment)
	actualUser.AddComment(newComment.ID)
	actualPost.AddComment(newComment.ID)
	return nil
}

// LIKE / UNLIKE

func (this *BlogService) LikePost(userID, postID uuid.UUID) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	actualPost, err := this.SearchPostByID(postID)
	if err != nil {
		return err
	}
	actualUser.LikePost(postID)
	actualPost.LikePost(postID)
	return nil
}

func (this *BlogService) UnlikePost(userID, postID uuid.UUID) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	actualPost, err := this.SearchPostByID(postID)
	if err != nil {
		return err
	}
	actualUser.UnlikePost(postID)
	actualPost.UnlikePost(postID)
	return nil
}

func (this *BlogService) LikeComment(userID, commentID uuid.UUID) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	actualComment, err := this.SearchCommentByID(commentID)
	if err != nil {
		return err
	}
	actualUser.LikeComment(commentID)
	actualComment.LikeComment(commentID)
	return nil
}

func (this *BlogService) UnlikeComment(userID, commentID uuid.UUID) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	actualComment, err := this.SearchCommentByID(commentID)
	if err != nil {
		return err
	}
	actualUser.UnlikeComment(commentID)
	actualComment.UnlikeComment(commentID)
	return nil
}

// PATCH

func (this *BlogService) EditUser(userID uuid.UUID, new_name, surname string) error {
	actualUser, err := this.SearchUserByID(userID)
	if err != nil {
		return err
	}
	actualUser.UpdateName(new_name)
	actualUser.UpdateSurname(surname)
	return nil
}

func (this *BlogService) EditPost(userID, postID uuid.UUID, title, new_body, imageURL string) error {
	actualPost, err := this.SearchPostByID(postID)
	if err != nil {
		return err
	}
	actualPost.UpdateTitle(title, userID)
	actualPost.UpdateBody(new_body, userID)
	actualPost.UpdateImageURL(imageURL, userID)
	return nil
}

// DELETE

func (p *BlogService) DeleteUser(userID uuid.UUID) error {
	for i, user := range p.users {
		if user.ID == userID {
			p.users[i] = p.users[len(p.users)-1]
			p.users = p.users[:len(p.users)-1]
			return nil
		}
	}
	return errors.New("user not found by id")
}

func (p *BlogService) DeletePost(postID uuid.UUID) error {
	for i, post := range p.posts {
		if post.ID == postID {
			p.posts[i] = p.posts[len(p.posts)-1]
			p.posts = p.posts[:len(p.posts)-1]
			return nil
		}
	}
	return errors.New("post not found by id")
}
