export const API_CONFIG = {
  AUTH_API_URL: 'http://localhost:8080',
  FORUM_API_URL: 'http://localhost:8081',
  GRPC_AUTH_URL: 'localhost:50051',
  WS: 'ws://localhost:8081/ws',
};

// URL для форумного сервиса
export const FORUM_API_URL = process.env.REACT_APP_FORUM_API_URL || 'http://localhost:8081/api/v1';

// URL для сервиса авторизации
export const AUTH_API_URL = process.env.REACT_APP_AUTH_API_URL || 'http://localhost:8080/api';

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: `${AUTH_API_URL}/auth/login`,
    REGISTER: `${AUTH_API_URL}/auth/register`,
    PROFILE: `${AUTH_API_URL}/auth/profile`,
    VERIFY: `${AUTH_API_URL}/auth/verify`,
    REFRESH: `${AUTH_API_URL}/auth/refresh`,
  },
  TOPICS: {
    GET_ALL: `${FORUM_API_URL}/topics`,
    GET_ONE: (id: string) => `${FORUM_API_URL}/topics/${id}`,
    CREATE: `${FORUM_API_URL}/topics`,
    UPDATE: (id: string) => `${FORUM_API_URL}/topics/${id}`,
    DELETE: (id: string) => `${FORUM_API_URL}/topics/${id}`,
    BASE: '/api/v1/topics',
  },
  COMMENTS: {
    GET_ALL: (topicId: string) => `${FORUM_API_URL}/topics/${topicId}/comments`,
    CREATE: (topicId: string) => `${FORUM_API_URL}/topics/${topicId}/comments`,
    UPDATE: (topicId: string, commentId: string) => `${FORUM_API_URL}/topics/${topicId}/comments/${commentId}`,
    DELETE: (topicId: string, commentId: string) => `${FORUM_API_URL}/topics/${topicId}/comments/${commentId}`,
    BASE: '/api/v1/comments',
  },
  CHAT: {
    WS: 'ws://localhost:8081/ws',
    MESSAGES: `${FORUM_API_URL}/chat/messages`,
  },
}; 