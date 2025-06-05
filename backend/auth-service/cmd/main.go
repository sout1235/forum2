package main

import (
	"log"
	"net"

	"github.com/sout1235/forum2/backend/auth-service/internal/config"
	grpcDelivery "github.com/sout1235/forum2/backend/auth-service/internal/delivery/grpc"
	httpDelivery "github.com/sout1235/forum2/backend/auth-service/internal/delivery/http"
	"github.com/sout1235/forum2/backend/auth-service/internal/repository"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	pb "github.com/sout1235/forum2/backend/auth-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Инициализируем репозиторий
	userRepo, err := repository.NewUserRepository(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// Инициализируем юзкейс
	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)

	// Инициализируем HTTP сервер
	router := httpDelivery.NewRouter(authUseCase, cfg.InternalToken).Setup()

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
		if err := router.Run(":" + cfg.HTTPPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Инициализация gRPC сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authServer := grpcDelivery.NewServer(authUseCase)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// Запуск gRPC сервера
	log.Printf("Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
