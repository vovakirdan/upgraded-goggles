package grpcserver

import (
	"context"
	"fmt"

	"upgraded-goggles/internal/post"

	postpb "upgraded-goggles/api/proto/post"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PostServer реализует gRPC-сервис для работы с постами
type PostServer struct {
	postpb.UnimplementedPostServiceServer
	service post.Service
}

// NewPostServer создаёт новый экземпляр PostServer
func NewPostServer(srv post.Service) *PostServer {
	return &PostServer{
		service: srv,
	}
}

// convertToProtoPost преобразует внутреннюю модель поста в proto-сообщение
func convertToProtoPost(p *post.Post) *postpb.Post {
	return &postpb.Post{
		Id:        p.ID,
		UserId:    p.UserID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: timestamppb.New(p.CreatedAt),
	}
}

// CreatePost обрабатывает запрос на создание поста
func (s *PostServer) CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.CreatePostResponse, error) {
	newPost, err := s.service.CreatePost(req.UserId, req.Title, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	return &postpb.CreatePostResponse{
		Post: convertToProtoPost(newPost),
	}, nil
}

// GetPost обрабатывает запрос на получение поста
func (s *PostServer) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.GetPostResponse, error) {
	p, err := s.service.GetPostByID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if p == nil {
		return nil, fmt.Errorf("post not found")
	}
	return &postpb.GetPostResponse{
		Post: convertToProtoPost(p),
	}, nil
}

// UpdatePost обрабатывает запрос на обновление поста
func (s *PostServer) UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.UpdatePostResponse, error) {
	updatedPost, err := s.service.UpdatePost(req.Id, req.Title, req.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	return &postpb.UpdatePostResponse{
		Post: convertToProtoPost(updatedPost),
	}, nil
}

// DeletePost обрабатывает запрос на удаление поста
func (s *PostServer) DeletePost(ctx context.Context, req *postpb.DeletePostRequest) (*postpb.DeletePostResponse, error) {
	if err := s.service.DeletePost(req.Id); err != nil {
		return nil, fmt.Errorf("failed to delete post: %w", err)
	}
	return &postpb.DeletePostResponse{
		Message: "post deleted successfully",
	}, nil
}
