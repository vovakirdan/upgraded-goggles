package gateway

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	userpb "upgraded-goggles/api/proto/user"
	postpb "upgraded-goggles/api/proto/post"
	gatewaypb "upgraded-goggles/api/proto/gateway"
)

// RegisterRoutes регистрирует маршруты для gRPC-сервисов в gRPC Gateway mux.
func RegisterRoutes(ctx context.Context, mux *runtime.ServeMux, cfg *Config, opts []grpc.DialOption) error {
	// Регистрация маршрута для UserService
	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.UserServiceAddress, opts)
	if err != nil {
		return fmt.Errorf("failed to register UserService: %w", err)
	}

	// Регистрация маршрута для PostService
	err = postpb.RegisterPostServiceHandlerFromEndpoint(ctx, mux, cfg.PostServiceAddress, opts)
	if err != nil {
		return fmt.Errorf("failed to register PostService: %w", err)
	}

	// Регистрация маршрута для API Gateway (например, HealthCheck)
	err = gatewaypb.RegisterAPIGatewayHandlerFromEndpoint(ctx, mux, cfg.GatewayAddress, opts)
	if err != nil {
		return fmt.Errorf("failed to register APIGateway: %w", err)
	}

	return nil
}
