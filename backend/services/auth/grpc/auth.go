package grpc

import (
	"context"

	"backend/services/auth/handlers"
	pb "backend/services/auth/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authHandler *handlers.AuthHandler
}

func NewAuthServer(authHandler *handlers.AuthHandler) *AuthServer {
	return &AuthServer{
		authHandler: authHandler,
	}
}

func (s *AuthServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
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

func (s *AuthServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
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
