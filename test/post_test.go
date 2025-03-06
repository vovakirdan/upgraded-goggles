package tests

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"upgraded-goggles/internal/post"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := post.NewRepository(db)
	userID := int64(1)
	title := "Test Post"
	content := "This is a test post."

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO posts (user_id, title, content, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
		WithArgs(userID, title, content, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	p := &post.Post{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.CreatePost(p); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if p.ID != 1 {
		t.Errorf("expected post ID to be 1, got %d", p.ID)
	}
}

func TestGetPostByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := post.NewRepository(db)
	postID := int64(1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = $1`)).
		WithArgs(postID).
		WillReturnError(sql.ErrNoRows)

	p, err := repo.GetPostByID(postID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if p != nil {
		t.Errorf("expected nil post, got %+v", p)
	}
}

// fakePostRepo реализует интерфейс post.Repository для тестирования
type fakePostRepo struct {
	post *post.Post
	err  error
}

func (f *fakePostRepo) CreatePost(p *post.Post) error {
	if f.err != nil {
		return f.err
	}
	p.ID = 1
	return nil
}

func (f *fakePostRepo) GetPostByID(id int64) (*post.Post, error) {
	return f.post, f.err
}

func (f *fakePostRepo) UpdatePost(p *post.Post) error {
	if f.err != nil {
		return f.err
	}
	return nil
}

func (f *fakePostRepo) DeletePost(id int64) error {
	return f.err
}

func TestCreatePostService(t *testing.T) {
	repo := &fakePostRepo{}
	service := post.NewService(repo)
	userID := int64(1)
	title := "Test Post"
	content := "This is a test post."

	p, err := service.CreatePost(userID, title, content)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if p == nil || p.ID != 1 {
		t.Errorf("expected valid post with ID 1, got %+v", p)
	}
}

func TestUpdatePostService_NotFound(t *testing.T) {
	repo := &fakePostRepo{post: nil}
	service := post.NewService(repo)
	title := "Updated Title"
	content := "Updated Content"

	p, err := service.UpdatePost(1, title, content)
	if err == nil {
		t.Error("expected error for post not found, got none")
	}
	if p != nil {
		t.Errorf("expected nil post, got %+v", p)
	}
}
