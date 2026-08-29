package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID       uuid.UUID `json:"id"`
	Body     string    `json:"body"`
	AuthorID uuid.UUID `json:"author_id"`

	UsersWhoLiked []uuid.UUID `json:"who_liked"`
	LikedCount    int         `json:"liked_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateComment(userid, postid uuid.UUID, body string) (Comment, error) {
	if len(body) == 0 {
		return Comment{}, errors.New("comment body cannot be empty")
	}
	return Comment{
		ID:            uuid.New(),
		Body:          body,
		AuthorID:      userid,
		UsersWhoLiked: make([]uuid.UUID, 0),
		LikedCount:    0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (c *Comment) ChangeBody(userid uuid.UUID, body string) error {
	if len(body) == 0 {
		return errors.New("updated comment body cannot be empty")
	}
	if c.AuthorID != userid {
		return errors.New("user is not the author of this comment")
	}
	c.Body = body
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Comment) LikeComment(userid uuid.UUID) error {
	for _, u := range c.UsersWhoLiked {
		if u == userid {
			return errors.New("user has already liked this comment")
		}
	}
	c.UsersWhoLiked = append(c.UsersWhoLiked, userid)
	c.LikedCount++
	return nil
}

func (c *Comment) UnlikeComment(userid uuid.UUID) error {
	likedIndex := -1
	for i, u := range c.UsersWhoLiked {
		if u == userid {
			likedIndex = i
			break
		}
	}
	if likedIndex == -1 {
		return errors.New("user has not liked this comment")
	}
	c.UsersWhoLiked[likedIndex] = c.UsersWhoLiked[len(c.UsersWhoLiked)-1]
	c.UsersWhoLiked = c.UsersWhoLiked[:len(c.UsersWhoLiked)-1]
	c.LikedCount--
	return nil
}
