package user

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"upgraded-goggles/internal/logger"
)

// Service описывает бизнес-логику для пользователей
type Service interface {
	RegisterUser(username, email, password string) (*User, error)
	LoginUser(email, password string) (string, *User, error) // возвращает токен и данные пользователя
	GetUserByID(id int64) (*User, error)
}

type service struct {
	repo Repository
}

// NewService создаёт новый сервис пользователей
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// RegisterUser регистрирует нового пользователя после валидации данных
func (s *service) RegisterUser(username, email, password string) (*User, error) {
	// Валидация данных
	if err := ValidateUserData(username, email, password); err != nil {
		return nil, err
	}
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	logger.Logger.Printf("User registered: %s", email)
	return user, nil
}

// LoginUser проверяет учетные данные и возвращает JWT-токен (заглушка) и данные пользователя
func (s *service) LoginUser(email, password string) (string, *User, error) {
	if email == "" || password == "" {
		return "", nil, errors.New("email and password are required")
	}
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("user not found")
	}
	// Сравнение хэшированного пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	// Генерация JWT-токена – здесь заглушка, замените на реальную реализацию
	token := "dummy-token"
	logger.Logger.Printf("User logged in: %s", email)
	return token, user, nil
}

// GetUserByID возвращает пользователя по его ID
func (s *service) GetUserByID(id int64) (*User, error) {
	return s.repo.GetUserByID(id)
}
