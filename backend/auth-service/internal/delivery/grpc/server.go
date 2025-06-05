package grpcserver

import (
	context "context"
	"fmt"

	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	pb "github.com/sout1235/forum2/backend/auth-service/proto"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	authUseCase *usecase.AuthUseCase
}

func NewServer(authUseCase *usecase.AuthUseCase) *AuthServiceServer {
	return &AuthServiceServer{
		authUseCase: authUseCase,
	}
}

func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// TODO: Реализовать валидацию токена через s.authUseCase
	if req.Token == "valid-token" {
		return &pb.ValidateTokenResponse{
			UserId: "1",
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}
