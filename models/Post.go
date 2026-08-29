package models

import (
	"BlogAPI/utils"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	ImageURL string    `json:"image_url"`

	AuthorID uuid.UUID `json:"author_id"`

	LikedCount    int         `json:"liked_count"`
	UsersWhoLiked []uuid.UUID `json:"who_liked"`
	CommentsID    []uuid.UUID `json:"comments"`
	CommentCount  int         `json:"comment_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreatePost(userID uuid.UUID, Title, Body, ImageURL string) (Post, error) {
	if len(Title) == 0 || len(Body) == 0 {
		log.Println("Title or Body of Post required")
		return Post{}, errors.New("title or Body of Post required")
	}
	return Post{
		ID:            uuid.New(),
		Title:         Title,
		Body:          Body,
		ImageURL:      ImageURL,
		AuthorID:      userID,
		LikedCount:    0,
		UsersWhoLiked: make([]uuid.UUID, 0),
		CommentsID:    make([]uuid.UUID, 0),
		CommentCount:  0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (p *Post) UpdateTitle(Title string, UserID uuid.UUID) error {
	if UserID != p.AuthorID {
		log.Printf("403 for updating post title")
		return errors.New("not authorized")
	}
	if len(Title) == 0 {
		log.Println("Title of Post required")
		return errors.New("title is required")
	} else if p.Title == Title {
		log.Println("title must be new")
		return errors.New("title must be new")
	}

	p.Title = Title
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Post) UpdateBody(body string, UserID uuid.UUID) error {
	if UserID != p.AuthorID {
		log.Printf("403 for updating post body")
		return errors.New("not authorized")
	}
	if len(body) == 0 {
		log.Println("body of Post required")
		return errors.New("body is required")
	} else if p.Body == body {
		log.Println("body must be new")
		return errors.New("body must be new")
	}

	p.Body = body
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Post) UpdateImageURL(imageURL string, UserID uuid.UUID) error {
	if UserID != p.AuthorID {
		log.Printf("403 for updating post imageurl")
		return errors.New("not authorized")
	}
	if len(imageURL) > 0 {
		if utils.IsImageURLValid(imageURL) != nil {
			log.Println("invalid imageurl")
			return errors.New("invalid imageurl")
		}
	}
	if p.ImageURL == imageURL {
		log.Println("imageURL must be new")
		return errors.New("imageURL must be new")
	}

	p.ImageURL = imageURL
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Post) AddComment(commentID uuid.UUID) {
	for _, id := range p.CommentsID {
		if id == commentID {
			return
		}
	}
	p.CommentsID = append(p.CommentsID, commentID)
	p.UpdatedAt = time.Now()
}

func (p *Post) RemoveComment(commentID uuid.UUID) {
	for i, id := range p.CommentsID {
		if id == commentID {
			p.CommentsID[i] = p.CommentsID[len(p.CommentsID)-1]
			p.CommentsID = p.CommentsID[:len(p.CommentsID)-1]
			p.UpdatedAt = time.Now()
		}
	}
}

func (p *Post) LikePost(userID uuid.UUID) {
	for _, id := range p.UsersWhoLiked {
		if id == userID {
			return
		}
	}
	p.UsersWhoLiked = append(p.UsersWhoLiked, userID)
	p.LikedCount++
	p.UpdatedAt = time.Now()
}

func (p *Post) UnlikePost(userID uuid.UUID) {
	for i, id := range p.UsersWhoLiked {
		if id == userID {
			p.UsersWhoLiked[i] = p.UsersWhoLiked[len(p.UsersWhoLiked)-1]
			p.UsersWhoLiked = p.UsersWhoLiked[:len(p.UsersWhoLiked)-1]
			p.LikedCount--
			p.UpdatedAt = time.Now()
		}
	}
}
