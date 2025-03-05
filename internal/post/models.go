package post

import "time"

// Post представляет пост пользователя в системе
type Post struct {
	ID        int64     `json:"id" db:"id"`                   // Идентификатор поста
	UserID    int64     `json:"user_id" db:"user_id"`         // Идентификатор автора (пользователя)
	Title     string    `json:"title" db:"title"`             // Заголовок поста
	Content   string    `json:"content" db:"content"`         // Текст поста
	CreatedAt time.Time `json:"created_at" db:"created_at"`   // Дата создания
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`   // Дата последнего обновления
}
