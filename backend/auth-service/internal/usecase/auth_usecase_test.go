package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
)

func TestAuthUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	useCase := NewAuthUseCase(mockRepo, "test-key")

	tests := []struct {
		name          string
		user          *entity.User
		setupMock     func()
		expectedError string
	}{
		{
			name: "successful registration",
			user: &entity.User{
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
			},
			setupMock: func() {
				mockRepo.EXPECT().ExistsByUsername(gomock.Any(), "testuser").Return(false, nil)
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), "test@example.com").Return(false, nil)
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "username already exists",
			user: &entity.User{
				Username:     "existinguser",
				Email:        "test@example.com",
				PasswordHash: "hashed_password",
			},
			setupMock: func() {
				mockRepo.EXPECT().ExistsByUsername(gomock.Any(), "existinguser").Return(true, nil)
			},
			expectedError: "username already exists",
		},
		{
			name: "email already exists",
			user: &entity.User{
				Username:     "testuser",
				Email:        "existing@example.com",
				PasswordHash: "hashed_password",
			},
			setupMock: func() {
				mockRepo.EXPECT().ExistsByUsername(gomock.Any(), "testuser").Return(false, nil)
				mockRepo.EXPECT().ExistsByEmail(gomock.Any(), "existing@example.com").Return(true, nil)
			},
			expectedError: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := useCase.Register(context.Background(), tt.user)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	useCase := NewAuthUseCase(mockRepo, "test-key")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		username      string
		password      string
		setupMock     func()
		expectedError string
	}{
		{
			name:     "successful login",
			username: "testuser",
			password: "correct_password",
			setupMock: func() {
				mockRepo.EXPECT().GetByUsername(gomock.Any(), "testuser").Return(&entity.User{
					ID:           "1",
					Username:     "testuser",
					PasswordHash: string(hashedPassword),
				}, nil)
			},
			expectedError: "",
		},
		{
			name:     "user not found",
			username: "nonexistent",
			password: "any_password",
			setupMock: func() {
				mockRepo.EXPECT().GetByUsername(gomock.Any(), "nonexistent").Return(nil, assert.AnError)
			},
			expectedError: "invalid credentials",
		},
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrong_password",
			setupMock: func() {
				mockRepo.EXPECT().GetByUsername(gomock.Any(), "testuser").Return(&entity.User{
					ID:           "1",
					Username:     "testuser",
					PasswordHash: string(hashedPassword),
				}, nil)
			},
			expectedError: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			tokens, user, err := useCase.Login(context.Background(), tt.username, tt.password)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, tokens)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotNil(t, user)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}

func TestAuthUseCase_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	useCase := NewAuthUseCase(mockRepo, "test-key")

	// Создаем тестовый токен
	claims := jwt.MapClaims{
		"user_id":  "123",
		"username": "testuser",
		"type":     "access",
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, _ := token.SignedString([]byte("test-key"))

	tests := []struct {
		name          string
		token         string
		expectedID    string
		expectedError string
	}{
		{
			name:       "valid token",
			token:      validToken,
			expectedID: "123",
		},
		{
			name:          "invalid token",
			token:         "invalid.token.here",
			expectedError: "token is malformed",
		},
		{
			name:          "expired token",
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTYxNjE2MTYsInVzZXJfaWQiOjF9.expired_signature",
			expectedError: "token is malformed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := useCase.ValidateToken(tt.token)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}

func TestAuthUseCase_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	useCase := NewAuthUseCase(mockRepo, "test-key")

	// Создаем тестовый refresh токен
	claims := jwt.MapClaims{
		"user_id":  "123",
		"username": "testuser",
		"type":     "refresh",
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validRefreshToken, _ := token.SignedString([]byte("test-key"))

	tests := []struct {
		name          string
		refreshToken  string
		setupMock     func()
		expectedError string
	}{
		{
			name:         "successful refresh",
			refreshToken: validRefreshToken,
			setupMock: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "123").Return(&entity.User{
					ID:       "123",
					Username: "testuser",
				}, nil)
			},
		},
		{
			name:          "invalid token",
			refreshToken:  "invalid.token.here",
			expectedError: "invalid refresh token",
		},
		{
			name:         "user not found",
			refreshToken: validRefreshToken,
			setupMock: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "123").Return(nil, assert.AnError)
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}
			tokens, err := useCase.RefreshToken(tt.refreshToken)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, tokens)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}

func TestAuthUseCase_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	useCase := NewAuthUseCase(mockRepo, "test-key")

	tests := []struct {
		name          string
		userID        string
		setupMock     func()
		expectedUser  *entity.User
		expectedError string
	}{
		{
			name:   "successful get user",
			userID: "123",
			setupMock: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "123").Return(&entity.User{
					ID:       "123",
					Username: "testuser",
					Email:    "test@example.com",
				}, nil)
			},
			expectedUser: &entity.User{
				ID:       "123",
				Username: "testuser",
				Email:    "test@example.com",
			},
		},
		{
			name:   "user not found",
			userID: "nonexistent",
			setupMock: func() {
				mockRepo.EXPECT().GetByID(gomock.Any(), "nonexistent").Return(nil, ErrUserNotFound)
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, err := useCase.GetUserByID(context.Background(), tt.userID)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}
		})
	}
}

func TestAuthUseCase_GenerateTokenPair(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	useCase := NewAuthUseCase(mockRepo, "test-key")

	user := &entity.User{
		ID:       "123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Настраиваем мок для вызова GetByID в RefreshToken
	mockRepo.EXPECT().GetByID(gomock.Any(), "123").Return(user, nil)

	tokens, err := useCase.GenerateTokenPair(user)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)

	// Проверяем, что токены можно валидировать
	userID, err := useCase.ValidateToken(tokens.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userID)

	// Проверяем refresh token
	newTokens, err := useCase.RefreshToken(tokens.RefreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, newTokens)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
}
