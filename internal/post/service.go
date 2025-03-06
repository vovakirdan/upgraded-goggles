package post

import (
	"errors"
	"fmt"
	"time"

	"upgraded-goggles/internal/logger"
)

// Service описывает бизнес-логику для работы с постами
type Service interface {
	CreatePost(userID int64, title, content string) (*Post, error)
	GetPostByID(id int64) (*Post, error)
	UpdatePost(id int64, title, content string) (*Post, error)
	DeletePost(id int64) error
}

type service struct {
	repo Repository
}

// NewService создаёт новый сервис постов
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreatePost создает новый пост после валидации данных
func (s *service) CreatePost(userID int64, title, content string) (*Post, error) {
	if err := ValidatePostData(title, content); err != nil {
		return nil, err
	}
	post := &Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreatePost(post); err != nil {
		return nil, err
	}
	if logger.Logger != nil {
		logger.Logger.Printf("Post created: %d", post.ID)
	} else {
		fmt.Printf("Post created: %d, but logger is nil", post.ID)
	}
	return post, nil
}

// GetPostByID возвращает пост по его ID
func (s *service) GetPostByID(id int64) (*Post, error) {
	return s.repo.GetPostByID(id)
}

// UpdatePost обновляет данные поста после валидации
func (s *service) UpdatePost(id int64, title, content string) (*Post, error) {
	if err := ValidatePostData(title, content); err != nil {
		return nil, err
	}
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("post not found")
	}
	post.Title = title
	post.Content = content
	post.UpdatedAt = time.Now()
	if err := s.repo.UpdatePost(post); err != nil {
		return nil, err
	}
	logger.Logger.Printf("Post updated: %d", post.ID)
	return post, nil
}

// DeletePost удаляет пост по ID
func (s *service) DeletePost(id int64) error {
	err := s.repo.DeletePost(id)
	if err != nil {
		return err
	}
	logger.Logger.Printf("Post deleted: %d", id)
	return nil
}
