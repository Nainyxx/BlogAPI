// PostgreSQL implementation of the repository interfaces, backed by db_schema.sql.
package repository

import (
	"BlogAPI/interfaces"
	"BlogAPI/models"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) interfaces.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, name, surname, login, email, password_hash, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Name, user.Surname, user.Login, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil && isUniqueViolation(err) {
		return errors.New("login or email already taken")
	}
	return err
}

func scanUser(row *sql.Row) (*models.User, error) {
	u := &models.User{
		PostsID:         make([]uuid.UUID, 0),
		LikedPostsID:    make([]uuid.UUID, 0),
		CommentsID:      make([]uuid.UUID, 0),
		LikedCommentsID: make([]uuid.UUID, 0),
	}
	err := row.Scan(&u.ID, &u.Name, &u.Surname, &u.Login, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}

const userColumns = `id, name, surname, login, email, password_hash, created_at, updated_at`

func (r *PostgresUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	row := r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *PostgresUserRepository) GetByLogin(login string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE login = $1`, login)
	return scanUser(row)
}

func (r *PostgresUserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT ` + userColumns + ` FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		u := models.User{
			PostsID:         make([]uuid.UUID, 0),
			LikedPostsID:    make([]uuid.UUID, 0),
			CommentsID:      make([]uuid.UUID, 0),
			LikedCommentsID: make([]uuid.UUID, 0),
		}
		if err := rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Login, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *PostgresUserRepository) Update(user *models.User) error {
	res, err := r.db.Exec(
		`UPDATE users SET name=$1, surname=$2, login=$3, email=$4, password_hash=$5, updated_at=$6 WHERE id=$7`,
		user.Name, user.Surname, user.Login, user.Email, user.PasswordHash, user.UpdatedAt, user.ID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return errors.New("login or email already taken")
		}
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *PostgresUserRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("user not found")
	}
	return nil
}

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) interfaces.PostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) Create(post *models.Post) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO posts (id, title, body, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		post.ID, post.Title, post.Body, post.AuthorID, post.CreatedAt, post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	if post.ImageURL != "" {
		_, err = tx.Exec(
			`INSERT INTO images (id, image_url, post_id, created_at) VALUES ($1, $2, $3, $4)`,
			uuid.New(), post.ImageURL, post.ID, post.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresPostRepository) enrich(post *models.Post) error {
	err := r.db.QueryRow(`SELECT image_url FROM images WHERE post_id=$1 ORDER BY created_at DESC LIMIT 1`, post.ID).Scan(&post.ImageURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id=$1`, post.ID).Scan(&post.CommentCount); err != nil {
		return err
	}
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM likes WHERE post_id=$1`, post.ID).Scan(&post.LikedCount); err != nil {
		return err
	}
	return nil
}

const postColumns = `id, title, body, author_id, created_at, updated_at`

func (r *PostgresPostRepository) scanPost(row *sql.Row) (*models.Post, error) {
	post := &models.Post{
		UsersWhoLiked: make([]uuid.UUID, 0),
		CommentsID:    make([]uuid.UUID, 0),
	}
	err := row.Scan(&post.ID, &post.Title, &post.Body, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	if err := r.enrich(post); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostgresPostRepository) GetByID(id uuid.UUID) (*models.Post, error) {
	row := r.db.QueryRow(`SELECT `+postColumns+` FROM posts WHERE id = $1`, id)
	return r.scanPost(row)
}

func (r *PostgresPostRepository) queryPosts(query string, args ...any) ([]models.Post, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{
			UsersWhoLiked: make([]uuid.UUID, 0),
			CommentsID:    make([]uuid.UUID, 0),
		}
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, err
		}
		if err := r.enrich(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (r *PostgresPostRepository) GetByAuthorID(authorID uuid.UUID) ([]models.Post, error) {
	return r.queryPosts(`SELECT `+postColumns+` FROM posts WHERE author_id = $1`, authorID)
}

func (r *PostgresPostRepository) GetAll() ([]models.Post, error) {
	return r.queryPosts(`SELECT ` + postColumns + ` FROM posts`)
}

func (r *PostgresPostRepository) Update(post *models.Post) error {
	res, err := r.db.Exec(`UPDATE posts SET title=$1, body=$2, updated_at=$3 WHERE id=$4`, post.Title, post.Body, post.UpdatedAt, post.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("post not found")
	}

	if post.ImageURL != "" {
		_, err = r.db.Exec(
			`INSERT INTO images (id, image_url, post_id, created_at) VALUES ($1, $2, $3, $4)`,
			uuid.New(), post.ImageURL, post.ID, time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresPostRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("post not found")
	}
	return nil
}

type PostgresCommentRepository struct {
	db *sql.DB
}

func NewPostgresCommentRepository(db *sql.DB) interfaces.CommentRepository {
	return &PostgresCommentRepository{db: db}
}

const commentColumns = `id, body, author_id, post_id, created_at, updated_at`

func (r *PostgresCommentRepository) Create(comment *models.Comment) error {
	_, err := r.db.Exec(
		`INSERT INTO comments (id, body, author_id, post_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		comment.ID, comment.Body, comment.AuthorID, comment.PostID, comment.CreatedAt, comment.UpdatedAt,
	)
	return err
}

func (r *PostgresCommentRepository) scanComment(row *sql.Row) (*models.Comment, error) {
	c := &models.Comment{
		UsersWhoLiked: make([]uuid.UUID, 0),
	}
	err := row.Scan(&c.ID, &c.Body, &c.AuthorID, &c.PostID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM likes WHERE comment_id=$1`, c.ID).Scan(&c.LikedCount); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *PostgresCommentRepository) GetByID(id uuid.UUID) (*models.Comment, error) {
	row := r.db.QueryRow(`SELECT `+commentColumns+` FROM comments WHERE id = $1`, id)
	return r.scanComment(row)
}

func (r *PostgresCommentRepository) queryComments(query string, args ...any) ([]models.Comment, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0)
	for rows.Next() {
		c := models.Comment{
			UsersWhoLiked: make([]uuid.UUID, 0),
		}
		if err := rows.Scan(&c.ID, &c.Body, &c.AuthorID, &c.PostID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if err := r.db.QueryRow(`SELECT COUNT(*) FROM likes WHERE comment_id=$1`, c.ID).Scan(&c.LikedCount); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *PostgresCommentRepository) GetByPostID(postID uuid.UUID) ([]models.Comment, error) {
	return r.queryComments(`SELECT `+commentColumns+` FROM comments WHERE post_id = $1`, postID)
}

func (r *PostgresCommentRepository) GetByAuthorID(authorID uuid.UUID) ([]models.Comment, error) {
	return r.queryComments(`SELECT `+commentColumns+` FROM comments WHERE author_id = $1`, authorID)
}

func (r *PostgresCommentRepository) GetAll() ([]models.Comment, error) {
	return r.queryComments(`SELECT ` + commentColumns + ` FROM comments`)
}

func (r *PostgresCommentRepository) Update(comment *models.Comment) error {
	res, err := r.db.Exec(`UPDATE comments SET body=$1, updated_at=$2 WHERE id=$3`, comment.Body, comment.UpdatedAt, comment.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func (r *PostgresCommentRepository) Delete(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM comments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("comment not found")
	}
	return nil
}

func (r *PostgresCommentRepository) Like(commentID, userID uuid.UUID) error {
	_, err := r.db.Exec(
		`INSERT INTO likes (id, user_id, comment_id, created_at) VALUES ($1, $2, $3, $4)`,
		uuid.New(), userID, commentID, time.Now(),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return errors.New("user has already liked this comment")
		}
		return err
	}
	return nil
}

func (r *PostgresCommentRepository) Unlike(commentID, userID uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM likes WHERE user_id=$1 AND comment_id=$2`, userID, commentID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("user has not liked this comment")
	}
	return nil
}
