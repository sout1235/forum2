import React, { useState, useEffect, useRef } from 'react';
import { API_ENDPOINTS } from '../config/api';
import { useAuth } from '../hooks/useAuth';
import axios from 'axios';

interface Message {
  type: string;
  content: string;
  author: string;
  timestamp: Date;
  id?: string;
}

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [isConnected, setIsConnected] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<string>('disconnected');
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();
  const { user } = useAuth();
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const messageIdsRef = useRef<Set<string>>(new Set());
  const isSendingRef = useRef<boolean>(false);
  const lastMessageTimestampRef = useRef<number>(0);

  // Функция для создания уникального ID сообщения
  const createMessageId = (message: Message) => {
    if (!message.content || !message.author) {
      console.error('Invalid message format:', message);
      return null;
    }
    return `${message.author}:${message.content}:${message.timestamp.getTime()}`;
  };

  // Функция для проверки и добавления сообщения
  const addMessage = (message: Message) => {
    if (!message.content || !message.author) {
      console.error('Invalid message format:', message);
      return;
    }

    const messageId = message.id || createMessageId(message);
    if (!messageId) return;

    if (messageIdsRef.current.has(messageId)) {
      console.log('Duplicate message detected, skipping:', message);
      return;
    }

    // Обновляем timestamp последнего сообщения
    const messageTime = message.timestamp.getTime();
    if (messageTime > lastMessageTimestampRef.current) {
      lastMessageTimestampRef.current = messageTime;
    }

    messageIdsRef.current.add(messageId);
    setMessages(prev => [...prev, message]);
  };

  // Загрузка сообщений из базы данных
  const loadMessages = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) return;

      const response = await axios.get(API_ENDPOINTS.CHAT.MESSAGES, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.data && Array.isArray(response.data)) {
        // Очищаем существующие сообщения и их ID
        setMessages([]);
        messageIdsRef.current.clear();
        lastMessageTimestampRef.current = 0;
        
        // Добавляем сообщения из базы данных
        response.data.forEach((msg: any) => {
          const message: Message = {
            type: 'message',
            content: msg.content,
            author: msg.author_username,
            timestamp: new Date(msg.created_at),
            id: msg.id
          };
          addMessage(message);
        });
      }
    } catch (err) {
      console.error('Error loading messages:', err);
      setError('Failed to load messages');
    }
  };

  const connectWebSocket = () => {
    if (!user) {
      setError('Please log in to use chat');
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      setError('No authentication token found');
      return;
    }

    // Clear any existing reconnect timeout
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = undefined;
    }

    try {
      console.log('Connecting to WebSocket...');
      setConnectionStatus('connecting');
      
      // Закрываем существующее соединение, если оно есть
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }

      const ws = new WebSocket(API_ENDPOINTS.CHAT.WS);
      wsRef.current = ws;

      ws.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        setError(null);
        setConnectionStatus('connected');
        
        // Отправляем токен для авторизации и последний timestamp
        const authMessage = {
          type: 'auth',
          token: token,
          lastMessageTimestamp: lastMessageTimestampRef.current
        };
        console.log('Sending auth message:', authMessage);
        ws.send(JSON.stringify(authMessage));
      };

      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          console.log('Received message:', data);

          if (data.type === 'auth_success') {
            console.log('Successfully authenticated in WebSocket');
            setIsAuthenticated(true);
            setError(null);
          } else if (data.type === 'error') {
            console.error('WebSocket error:', data.content);
            setError(data.content);
            if (data.content === 'You must authenticate first') {
              setIsAuthenticated(false);
            }
          } else if (data.type === 'message') {
            console.log('Adding message to chat:', data);
            const message: Message = {
              type: data.type,
              content: data.content,
              author: data.author,
              timestamp: new Date(data.timestamp || Date.now()),
              id: data.id
            };
            addMessage(message);
          } else if (data.type === 'pong') {
            console.log('Received pong from server');
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
          setError('Failed to parse message from server');
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setError('Connection error occurred');
        setIsConnected(false);
        setConnectionStatus('error');
      };

      ws.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        setIsConnected(false);
        setConnectionStatus('disconnected');
        
        // Очищаем текущее соединение
        wsRef.current = null;
        
        // Пробуем переподключиться через 5 секунд, только если это не было намеренное закрытие
        if (event.code !== 1000 && !reconnectTimeoutRef.current) {
          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectTimeoutRef.current = undefined;
            connectWebSocket();
          }, 5000);
        }
      };

      // Настраиваем пинг-понг
      const pingInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: 'ping' }));
        }
      }, 25000); // Отправляем пинг каждые 25 секунд

      // Очищаем интервал при размонтировании
      return () => {
        clearInterval(pingInterval);
        if (ws.readyState === WebSocket.OPEN) {
          ws.close(1000, 'Component unmounting');
        }
      };
    } catch (err) {
      console.error('Error creating WebSocket:', err);
      setError('Failed to create WebSocket connection');
      setConnectionStatus('error');
    }
  };

  useEffect(() => {
    if (user) {
      loadMessages();
      connectWebSocket();
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close(1000, 'Component unmounting');
        wsRef.current = null;
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, [user]);

  // Scroll to bottom when messages change
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  const sendMessage = () => {
    if (!wsRef.current || !isAuthenticated) {
      console.error('Cannot send message:', {
        wsExists: !!wsRef.current,
        messageEmpty: !newMessage.trim(),
        isConnected,
        isAuthenticated,
        isSending: isSendingRef.current
      });
      setError('Не удалось отправить сообщение. Проверьте подключение и авторизацию.');
      return;
    }

    if (!newMessage.trim()) {
      return;
    }

    try {
      isSendingRef.current = true;
      const message = {
        type: 'message',
        content: newMessage.trim()
      };
      console.log('Sending message to server:', message);
      wsRef.current.send(JSON.stringify(message));
      setNewMessage('');
    } catch (err) {
      console.error('Error sending message:', err);
      setError('Не удалось отправить сообщение');
    } finally {
      isSendingRef.current = false;
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  if (!user) {
    return (
      <div className="p-4 bg-white rounded-lg shadow">
        <p className="text-gray-600">Please log in to use chat</p>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-[600px] bg-white rounded-lg shadow">
      <div className="p-4 border-b">
        <h2 className="text-xl font-semibold">Chat</h2>
        <div className="flex items-center space-x-2 mt-1">
          <div className={`w-2 h-2 rounded-full ${
            connectionStatus === 'connected' ? 'bg-green-500' :
            connectionStatus === 'connecting' || connectionStatus === 'reconnecting' ? 'bg-yellow-500' :
            'bg-red-500'
          }`} />
          <p className={`text-sm ${
            connectionStatus === 'connected' ? 'text-green-500' :
            connectionStatus === 'connecting' || connectionStatus === 'reconnecting' ? 'text-yellow-500' :
            'text-red-500'
          }`}>
            {connectionStatus === 'connected' ? 'Connected' :
             connectionStatus === 'connecting' ? 'Connecting...' :
             connectionStatus === 'reconnecting' ? 'Reconnecting...' :
             'Disconnected'}
          </p>
        </div>
        {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((msg, index) => (
          <div
            key={index}
            className={`flex flex-col ${msg.author === user.username ? 'items-end' : 'items-start'}`}
          >
            <div
              className={`max-w-[70%] rounded-lg p-3 shadow-md transition-all duration-200 ${
                msg.author === user.username
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-100 text-gray-800'
              }`}
            >
              <p className="text-xs font-semibold mb-1 opacity-70">{msg.author}</p>
              <p className="text-lg leading-snug">{msg.content}</p>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <div className="p-4 border-t">
        <div className="flex space-x-2">
          <textarea
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type a message..."
            className="flex-1 p-2 border rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-blue-500"
            rows={2}
            disabled={!isConnected}
          />
          <button
            onClick={sendMessage}
            disabled={!isConnected || !newMessage.trim()}
            className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
};

export default Chat; 