package user

import (
	"errors"
	"regexp"
)

const regEmail = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

// ValidateUserData проверяет корректность данных пользователя
func ValidateUserData(username, email, password string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	// Простейшая проверка формата email
	re := regexp.MustCompile(regEmail)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
