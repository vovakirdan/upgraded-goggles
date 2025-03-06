package post

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository описывает интерфейс для работы с постами в БД
type Repository interface {
	CreatePost(post *Post) error
	GetPostByID(id int64) (*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id int64) error
}

type repository struct {
	db *sql.DB
}

// NewRepository создаёт новый репозиторий для постов
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreatePost создает новый пост в БД
func (r *repository) CreatePost(post *Post) error {
	query := `INSERT INTO posts (user_id, title, content, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, post.UserID, post.Title, post.Content, time.Now(), time.Now()).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

// GetPostByID возвращает пост по его ID
func (r *repository) GetPostByID(id int64) (*Post, error) {
	query := `SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var post Post
	if err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // пост не найден
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return &post, nil
}

// UpdatePost обновляет данные поста
func (r *repository) UpdatePost(post *Post) error {
	query := `UPDATE posts SET title = $1, content = $2, updated_at = $3 WHERE id = $4`
	res, err := r.db.Exec(query, post.Title, post.Content, time.Now(), post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update post, could not fetch affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no post found with id: %d", post.ID)
	}
	return nil
}

// DeletePost удаляет пост из БД по ID
func (r *repository) DeletePost(id int64) error {
	query := `DELETE FROM posts WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete post, could not fetch affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("no post found with id: %d", id)
	}
	return nil
}
