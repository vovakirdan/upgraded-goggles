package grpcserver

import (
	"context"
	"fmt"

	"upgraded-goggles/internal/user"

	userpb "upgraded-goggles/api/proto/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserServer реализует gRPC-сервис для работы с пользователями
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	service user.Service
}

// NewUserServer создаёт новый экземпляр UserServer
func NewUserServer(srv user.Service) *UserServer {
	return &UserServer{
		service: srv,
	}
}

// convertToProto преобразует внутреннюю модель пользователя в proto-сообщение
func convertToProto(u *user.User) *userpb.User {
	return &userpb.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
}

// Register обрабатывает запрос на регистрацию пользователя
func (s *UserServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	newUser, err := s.service.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	return &userpb.RegisterResponse{
		User: convertToProto(newUser),
	}, nil
}

// Login обрабатывает запрос на авторизацию пользователя
func (s *UserServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	token, usr, err := s.service.LoginUser(req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login user: %w", err)
	}
	return &userpb.LoginResponse{
		Token: token,
		User:  convertToProto(usr),
	}, nil
}

// GetUser обрабатывает запрос на получение информации о пользователе
func (s *UserServer) GetUser(ctx context.Context, req *userpb.UserRequest) (*userpb.User, error) {
	usr, err := s.service.GetUserByID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if usr == nil {
		return nil, fmt.Errorf("user not found")
	}
	return convertToProto(usr), nil
}
