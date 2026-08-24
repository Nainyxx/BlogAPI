package service

import (
	"BlogAPI/models"
	"errors"

	"github.com/google/uuid"
)

/*
Blog service - main data controlling component

Implented main methods:

	RegisterUser - user registration method
	CreatePost - creating post A by user B. post A.author = B.ID
	UpdatePost - post editing method (change name, body or image url)
	DeletePost - deleting post method
	WriteComment
	LikeComment
	DislikeComment
	EditComment

Implented utils methods:

	SearchUserByID - searching user object by his ID in BlogService.users
	SearchPostByID - searching post object by his ID in BlogService.posts
	IsLoginTaken - is user login taken or not validate

TODO methods:
*/
type BlogService struct {
	users    []models.User
	posts    []models.Post
	comments []models.Comment
}

func NewBlogService() *BlogService {
	return &BlogService{
		users:    make([]models.User, 0),
		posts:    make([]models.Post, 0),
		comments: make([]models.Comment, 0),
	}
}

func (bs *BlogService) SearchUserByID(userID uuid.UUID) (*models.User, error) {
	for i := range bs.users {
		if bs.users[i].ID == userID {
			return &bs.users[i], nil
		}
	}
	return nil, errors.New("User not found")
}

func (bs *BlogService) SearchPostByID(postID uuid.UUID) (*models.Post, error) {
	for i := range bs.posts {
		if bs.posts[i].ID == postID {
			return &bs.posts[i], nil
		}
	}
	return nil, errors.New("post not found")
}

func (bs *BlogService) SearchCommentByID(commentID uuid.UUID) (*models.Comment, error) {
	for i := range bs.comments {
		if bs.comments[i].ID == commentID {
			return &bs.comments[i], nil
		}
	}
	return nil, errors.New("comment not found")
}

func (bs *BlogService) UserIsExist(userID uuid.UUID) error {
	for i := range bs.users {
		if bs.users[i].ID == userID {
			return nil
		}
	}
	return errors.New("User not found")
}

func (bs *BlogService) PostIsExist(postID uuid.UUID) error {
	for i := range bs.posts {
		if bs.posts[i].ID == postID {
			return nil
		}
	}
	return errors.New("post not found")
}

func (bs *BlogService) CommentIsExist(commentID uuid.UUID) error {
	for i := range bs.comments {
		if bs.comments[i].ID == commentID {
			return nil
		}
	}
	return errors.New("comment not found")
}

func (bs *BlogService) IsLoginTaken(login string) bool {
	for _, user := range bs.users {
		if user.Login == login {
			return true
		}
	}
	return false
}

func (bs *BlogService) RegisterUser(name, surname, login, email, password string) error {
	if bs.IsLoginTaken(login) {
		return errors.New("Login taken")
	}
	newUser, err := models.CreateUser(name, surname, login, email, password)
	if err != nil {
		return err
	}

	bs.users = append(bs.users, newUser)
	return nil
}

func (bs *BlogService) CreatePost(userID uuid.UUID, title, body, imageURL string) error {
	user, err := bs.SearchUserByID(userID)
	if err != nil {
		return err
	}

	newPost, err := models.CreatePost(userID, title, body, imageURL)
	if err != nil {
		return err
	}
	bs.posts = append(bs.posts, newPost)
	user.AddPost(newPost.ID)
	return nil
}

func (bs *BlogService) UpdatePost(userID, postID uuid.UUID, newTitle, newBody, newImageURL string) error {
	editedPost, err := bs.SearchPostByID(postID)
	if err != nil {
		return err
	}

	hasChanges := false

	if editedPost.Title != newTitle {
		err = editedPost.UpdateTitle(newTitle, userID)
		if err != nil {
			return err
		}
		hasChanges = true
	}
	if editedPost.Body != newBody {
		err = editedPost.UpdateBody(newBody, userID)
		if err != nil {
			return err
		}
		hasChanges = true
	}
	if editedPost.ImageURL != newImageURL {
		err = editedPost.UpdateImageURL(newImageURL, userID)
		if err != nil {
			return err
		}
		hasChanges = true
	}
	if hasChanges == false {
		return errors.New("No changes")
	}
	return nil
}

func (bs *BlogService) DeletePost(userID, postID uuid.UUID) error {
	editedPost, err := bs.SearchPostByID(postID)
	if err != nil {
		return err
	}

	if editedPost.AuthorID != userID {
		return errors.New("not authorized to delete this post")
	}

	user, err := bs.SearchUserByID(userID)
	if err == nil {
		user.RemovePost(postID)
	}

	postIndex := -1
	for i := range bs.posts {
		if bs.posts[i].ID == postID {
			postIndex = i
			break
		}
	}

	if postIndex == -1 {
		return errors.New("post not found in global list")
	}

	bs.posts[postIndex] = bs.posts[len(bs.posts)-1]
	bs.posts = bs.posts[:len(bs.posts)-1]

	return nil
}

func (bs *BlogService) WriteComment(userID, postID uuid.UUID, body string) error {
	post, err := bs.SearchPostByID(postID)
	if err != nil {
		return err
	}
	user, err := bs.SearchUserByID(userID)
	if err != nil {
		return err
	}
	newComment, err := models.CreateComment(userID, postID, body)
	if err != nil {
		return err
	}
	bs.comments = append(bs.comments, newComment)
	post.CommentsID = append(post.CommentsID, newComment.ID)
	post.CommentCount++
	user.CommentsID = append(user.CommentsID, newComment.ID)
	return nil
}

func (bs *BlogService) LikeComment(userID, commentID uuid.UUID) error {
	user, err := bs.SearchUserByID(userID)
	if err != nil {
		return err
	}
	comment, err := bs.SearchCommentByID(commentID)
	if err != nil {
		return err
	}
	user.LikeComment(commentID)
	return comment.LikeComment(userID)
}

func (bs *BlogService) DislikeComment(userID, commentID uuid.UUID) error {
	user, err := bs.SearchUserByID(userID)
	if err != nil {
		return err
	}
	comment, err := bs.SearchCommentByID(commentID)
	if err != nil {
		return err
	}
	user.UnlikeComment(commentID)
	return comment.DislikeComment(userID)
}

func (bs *BlogService) EditComment(userID, commentID uuid.UUID, newBody string) error {
	err := bs.UserIsExist(userID)
	if err != nil {
		return err
	}
	editComment, err := bs.SearchCommentByID(commentID)
	if err != nil {
		return err
	}
	err = editComment.ChangeBody(userID, newBody)
	return err
}
