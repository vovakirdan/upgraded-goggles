package post

import "errors"

// ValidatePostData проверяет корректность данных поста
func ValidatePostData(title, content string) error {
	if title == "" {
		return errors.New("title is required")
	}
	if content == "" {
		return errors.New("content is required")
	}
	return nil
}
