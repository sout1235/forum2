package chat

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Message представляет сообщение чата
type Message struct {
	UserID    int64
	Username  string
	Content   string
	Timestamp time.Time
}

// Chat управляет чатом и WebSocket-подключениями
type Chat struct {
	clients    map[*websocket.Conn]bool
	messages   []Message
	mutex      sync.Mutex
	upgrader   websocket.Upgrader
	cleanupAge time.Duration
}

// NewChat создает новый экземпляр чата
func NewChat(cleanupAge time.Duration) *Chat {
	return &Chat{
		clients:    make(map[*websocket.Conn]bool),
		messages:   make([]Message, 0),
		upgrader:   websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		cleanupAge: cleanupAge,
	}
}

// HandleWebSocket обрабатывает WebSocket-подключения
func (c *Chat) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении соединения: %v", err)
		return
	}
	defer conn.Close()

	c.mutex.Lock()
	c.clients[conn] = true
	c.mutex.Unlock()

	// Отправляем историю сообщений новому клиенту
	c.mutex.Lock()
	for _, msg := range c.messages {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Ошибка при отправке истории: %v", err)
			break
		}
	}
	c.mutex.Unlock()

	// Читаем сообщения от клиента
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			break
		}
		msg.Timestamp = time.Now()
		c.mutex.Lock()
		c.messages = append(c.messages, msg)
		// Рассылаем сообщение всем клиентам
		for client := range c.clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
				client.Close()
				delete(c.clients, client)
			}
		}
		c.mutex.Unlock()
	}

	c.mutex.Lock()
	delete(c.clients, conn)
	c.mutex.Unlock()
}

// StartCleanup запускает периодическую очистку старых сообщений
func (c *Chat) StartCleanup() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			c.mutex.Lock()
			now := time.Now()
			var newMessages []Message
			for _, msg := range c.messages {
				if now.Sub(msg.Timestamp) < c.cleanupAge {
					newMessages = append(newMessages, msg)
				}
			}
			c.messages = newMessages
			c.mutex.Unlock()
		}
	}()
}
