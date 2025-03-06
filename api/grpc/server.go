package grpcserver

import (
	"log"
	"net"

	"google.golang.org/grpc"

	userpb "upgraded-goggles/api/proto/user"
	postpb "upgraded-goggles/api/proto/post"
)

// StartUserGRPCServer запускает gRPC сервер для пользователей на указанном порту
func StartUserGRPCServer(port string, userService *UserServer) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userService)
	log.Printf("User gRPC server started on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve user gRPC server: %v", err)
	}
}

// StartPostGRPCServer запускает gRPC сервер для постов на указанном порту
func StartPostGRPCServer(port string, postService *PostServer) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	postpb.RegisterPostServiceServer(grpcServer, postService)
	log.Printf("Post gRPC server started on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve post gRPC server: %v", err)
	}
}
