# Forum Project

Микросервисный форум с чатом на WebSocket.

## Структура проекта

```
.
├── backend/
│   ├── pkg/           # Общие пакеты
│   └── services/
│       ├── auth/      # Сервис аутентификации
│       └── forum/     # Сервис форума
└── frontend/          # React приложение
```

## Технологии

### Backend
- Go
- PostgreSQL
- gRPC
- WebSocket
- Zap Logger

### Frontend
- React
- Tailwind CSS
- WebSocket

## Установка и запуск

### Backend

1. Установите зависимости:
```bash
cd backend/services/auth
go mod tidy
cd ../forum
go mod tidy
```

2. Запустите сервисы:
```bash
# В первом терминале
cd backend/services/auth
go run main.go

# Во втором терминале
cd backend/services/forum
go run main.go
```

### Frontend

1. Установите зависимости:
```bash
cd frontend
npm install
```

2. Запустите приложение:
```bash
npm start
```

## API

### Auth Service
- gRPC порт: 50051
- HTTP порт: 8080

### Forum Service
- gRPC порт: 50052
- HTTP порт: 8081
- WebSocket: ws://localhost:8081/ws/chat

## Лицензия

MIT 