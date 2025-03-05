package user

import "time"

// User представляет пользователя системы
type User struct {
	ID        int64     `json:"id" db:"id"`                   // Идентификатор пользователя
	Username  string    `json:"username" db:"username"`       // Имя пользователя
	Email     string    `json:"email" db:"email"`             // Электронная почта
	Password  string    `json:"password" db:"password"`       // Хэшированный пароль
	CreatedAt time.Time `json:"created_at" db:"created_at"`   // Дата создания
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`   // Дата последнего обновления
}
