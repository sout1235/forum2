package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sout1235/forum2/backend/forum-service/api/proto"
	"github.com/sout1235/forum2/backend/forum-service/internal/config"
	grpcDelivery "github.com/sout1235/forum2/backend/forum-service/internal/delivery/grpc"
	httpDelivery "github.com/sout1235/forum2/backend/forum-service/internal/delivery/http"
	"github.com/sout1235/forum2/backend/forum-service/internal/delivery/http/middleware"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
	"github.com/sout1235/forum2/backend/forum-service/internal/service"
	"github.com/sout1235/forum2/backend/forum-service/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.NewConfig()

	// Подключение к базе данных
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Применяем миграции
	if err := repository.RunMigrations(db); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// Инициализация репозиториев
	topicRepo := repository.NewTopicRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo := repository.NewUserRepository(db, cfg.AuthServiceURL)
	chatRepo := repository.NewChatRepository(db)

	// Инициализация use cases
	commentUseCase := usecase.NewCommentUseCase(commentRepo, userRepo)
	topicService := service.NewTopicService(topicRepo, userRepo)

	// Инициализация HTTP сервера
	authConfig := &middleware.AuthConfig{
		AuthServiceURL: cfg.AuthServiceURL,
	}

	router := httpDelivery.NewRouter(
		topicService,
		commentUseCase,
		userRepo,
		chatRepo,
		cfg.AuthServiceURL,
		authConfig,
	)

	// Запуск HTTP сервера
	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.HTTPPort)
		if err := router.Run(":" + cfg.HTTPPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Инициализация gRPC сервера
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	commentServer := grpcDelivery.NewCommentServer(commentUseCase)
	proto.RegisterCommentServiceServer(grpcServer, commentServer)

	// Запуск gRPC сервера
	log.Printf("Starting gRPC server on :%s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
