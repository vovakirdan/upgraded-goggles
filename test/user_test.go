package tests

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"upgraded-goggles/internal/user"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	repo := user.NewRepository(db)
	username := "testuser"
	email := "test@example.com"
	password := "hashedpassword"

	// Ожидаем запрос INSERT с соответствующими параметрами
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (username, email, password, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
		WithArgs(username, email, password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	u := &user.User{
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.CreateUser(u); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if u.ID != 1 {
		t.Errorf("expected user ID to be 1, got %d", u.ID)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error '%s' when opening stub database", err)
	}
	defer db.Close()

	repo := user.NewRepository(db)
	email := "nonexistent@example.com"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	u, err := repo.GetUserByEmail(email)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if u != nil {
		t.Errorf("expected nil user, got %+v", u)
	}
}

// fakeUserRepo реализует интерфейс user.Repository для тестирования
type fakeUserRepo struct {
	user *user.User
	err  error
}

func (f *fakeUserRepo) CreateUser(u *user.User) error {
	if f.err != nil {
		return f.err
	}
	u.ID = 1
	return nil
}

func (f *fakeUserRepo) GetUserByID(id int64) (*user.User, error) {
	return f.user, f.err
}

func (f *fakeUserRepo) GetUserByEmail(email string) (*user.User, error) {
	if f.user != nil && f.user.Email == email {
		return f.user, nil
	}
	return nil, f.err
}

func TestRegisterUser(t *testing.T) {
	repo := &fakeUserRepo{}
	service := user.NewService(repo)
	username := "testuser"
	email := "test@example.com"
	password := "password123"

	u, err := service.RegisterUser(username, email, password)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if u == nil || u.ID != 1 {
		t.Errorf("expected valid user with ID 1, got %+v", u)
	}
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	// Создаем пользователя с хэшированным паролем
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &user.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo := &fakeUserRepo{user: existingUser}
	service := user.NewService(repo)

	// Пробуем залогиниться с неверным паролем
	token, u, err := service.LoginUser("test@example.com", "wrongpassword")
	if err == nil {
		t.Error("expected error for invalid credentials, got none")
	}
	if token != "" || u != nil {
		t.Errorf("expected empty token and nil user on failure, got token: %s, user: %+v", token, u)
	}
}
