package grpcserver

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase/mocks"
	pb "github.com/sout1235/forum2/backend/auth-service/proto"
	"github.com/stretchr/testify/assert"
)

func TestAuthServiceServer_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	authUseCase := usecase.NewAuthUseCase(mockRepo, "test-secret")
	server := NewServer(authUseCase)

	tests := []struct {
		name          string
		token         string
		mockBehavior  func()
		expectedValid bool
		expectedError bool
	}{
		{
			name:  "valid token",
			token: "valid-token",
			mockBehavior: func() {
				// Здесь можно добавить моки для проверки токена
			},
			expectedValid: true,
			expectedError: false,
		},
		{
			name:          "invalid token",
			token:         "invalid-token",
			mockBehavior:  func() {},
			expectedValid: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			req := &pb.ValidateTokenRequest{
				Token: tt.token,
			}

			resp, err := server.ValidateToken(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) {
					assert.Equal(t, tt.expectedValid, resp.UserId != "")
				}
			}
		})
	}
}
