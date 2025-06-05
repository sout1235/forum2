package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/sout1235/forum2/backend/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtKey string) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, user *entity.User) error {
	// Проверяем уникальность username
	exists, err := uc.userRepo.ExistsByUsername(ctx, user.Username)
	if err != nil {
		return fmt.Errorf("error checking username existence: %v", err)
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	// Проверяем уникальность email
	exists, err = uc.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("error checking email existence: %v", err)
	}
	if exists {
		return fmt.Errorf("email already exists")
	}

	// Пароль уже хеширован в handler, поэтому просто сохраняем его
	return uc.userRepo.Create(ctx, user)
}

func (uc *AuthUseCase) Login(ctx context.Context, username, password string) (*entity.TokenPair, *entity.User, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Генерируем токены
	tokens, err := uc.GenerateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating tokens: %v", err)
	}

	return tokens, user, nil
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return uc.jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, userID)
}

func (uc *AuthUseCase) GenerateTokenPair(user *entity.User) (*entity.TokenPair, error) {
	// Генерируем access token (15 минут)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"type":     "access",
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(uc.jwtKey)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh token (7 дней)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"type":     "refresh",
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(uc.jwtKey)
	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(refreshToken string) (*entity.TokenPair, error) {
	// Проверяем refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return uc.jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token claims")
	}

	// Проверяем тип токена
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Получаем пользователя
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	user, err := uc.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Генерируем новую пару токенов
	return uc.GenerateTokenPair(user)
}
