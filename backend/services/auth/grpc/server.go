package grpc

import (
	"context"
	"log"
	"net"

	"backend/services/auth/handlers"
	pb "backend/services/auth/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	authHandler *handlers.AuthHandler
}

func NewServer(authHandler *handlers.AuthHandler) *Server {
	return &Server{
		authHandler: authHandler,
	}
}

func (s *Server) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	claims, err := s.authHandler.ParseJWT(req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.VerifyTokenResponse{
		UserId: int64(claims["user_id"].(float64)),
		Role:   claims["role"].(string),
		Valid:  true,
	}, nil
}

func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	user, err := s.authHandler.Repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Avatar:   user.Avatar,
	}, nil
}

func StartServer(authHandler *handlers.AuthHandler, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, NewServer(authHandler))

	log.Printf("Starting gRPC server on port %s", port)
	return s.Serve(lis)
}
