package grpc

import (
	"context"

	"backend/services/auth/proto"

	"github.com/sout1235/forumski/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client proto.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(authServiceAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to auth service", zap.Error(err))
		return nil, err
	}

	client := proto.NewAuthServiceClient(conn)
	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *AuthClient) VerifyToken(ctx context.Context, token string) (int64, string, error) {
	resp, err := c.client.VerifyToken(ctx, &proto.VerifyTokenRequest{
		Token: token,
	})
	if err != nil {
		logger.Error("Failed to verify token", zap.Error(err))
		return 0, "", err
	}
	return resp.UserId, resp.Role, nil
}

func (c *AuthClient) GetUserByID(ctx context.Context, userID int64) (*proto.User, error) {
	resp, err := c.client.GetUserByID(ctx, &proto.GetUserByIDRequest{
		UserId: userID,
	})
	if err != nil {
		logger.Error("Failed to get user by ID", zap.Error(err))
		return nil, err
	}
	return &proto.User{
		Id:       resp.Id,
		Username: resp.Username,
		Email:    resp.Email,
		Role:     resp.Role,
		Avatar:   resp.Avatar,
	}, nil
}
