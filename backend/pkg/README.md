# Forum Common

Общие модули для форум-проекта.

## Модули

### Logger

Модуль логирования на основе zap, предоставляющий единый интерфейс для логирования во всех сервисах.

#### Использование

```go
import "github.com/yourusername/forum-common/logger"

func main() {
    logger.Init()
    logger.Info("Application started")
    
    // Использование с полями
    logger.Info("User logged in",
        zap.String("username", "john"),
        zap.Int("user_id", 123),
    )
}
```

## Установка

```bash
go get github.com/yourusername/forum-common
```

## Лицензия

MIT 