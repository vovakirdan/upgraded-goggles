package user

import (
	"database/sql"
	"fmt"
	"time"
)

// Repository описывает интерфейс для работы с пользователями в БД.
type Repository interface {
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository создаёт новый репозиторий пользователей.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreateUser создает нового пользователя и записывает его в БД.
func (r *repository) CreateUser(user *User) error {
	query := `INSERT INTO users (username, email, password, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, time.Now(), time.Now()).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID возвращает пользователя по его ID.
func (r *repository) GetUserByID(id int64) (*User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail возвращает пользователя по email.
func (r *repository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}
