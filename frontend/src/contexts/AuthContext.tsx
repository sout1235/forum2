import React, { createContext, useState, useEffect } from 'react';
import axios from 'axios';
import { User } from '../types/user';
import { AuthContextType } from '../hooks/useAuth';
import { API_ENDPOINTS, API_CONFIG } from '../config/api';

export const AuthContext = createContext<AuthContextType | null>(null);

// Создаем axios instance с перехватчиками
const api = axios.create();

// Добавляем перехватчик ответов
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Если ошибка 401 и это не запрос на обновление токена
    if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/auth/refresh')) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refreshToken');
        if (!refreshToken) {
          window.location.href = '/login';
          return Promise.reject(error);
        }

        // Пытаемся обновить токен
        const response = await api.post(API_ENDPOINTS.AUTH.REFRESH, {
          refresh_token: refreshToken
        });

        const { access_token, refresh_token } = response.data;

        // Сохраняем новые токены
        localStorage.setItem('token', access_token);
        localStorage.setItem('refreshToken', refresh_token);

        // Обновляем заголовок в оригинальном запросе
        originalRequest.headers['Authorization'] = `Bearer ${access_token}`;

        // Повторяем оригинальный запрос
        return api(originalRequest);
      } catch (refreshError) {
        localStorage.removeItem('token');
        localStorage.removeItem('refreshToken');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const checkIsAdmin = (user: User) => {
    return user.role === 'admin';
  };

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('token');
      
      if (!token) {
        setIsLoading(false);
        return;
      }

      try {
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        
        const response = await api.get<User>(API_ENDPOINTS.AUTH.PROFILE);

        if (response.data) {
          setUser(response.data);
          setIsAdmin(checkIsAdmin(response.data));
        }
      } catch (error) {
        if (axios.isAxiosError(error) && error.response?.status === 401) {
          localStorage.removeItem('token');
          localStorage.removeItem('refreshToken');
          delete api.defaults.headers.common['Authorization'];
          setUser(null);
          setIsAdmin(false);
        }
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (username: string, password: string) => {
    try {
      console.log('Начало процесса входа для пользователя:', username);
      console.log('Отправка запроса на:', API_ENDPOINTS.AUTH.LOGIN);
      
      // Создаем копию axios с дополнительными настройками для отладки
      const debugApi = axios.create({
        baseURL: API_CONFIG.AUTH_API_URL,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      });

      // Добавляем перехватчик для логирования запроса
      debugApi.interceptors.request.use(request => {
        console.log('Отправляемый запрос:', {
          url: request.url,
          method: request.method,
          headers: request.headers,
          data: request.data
        });
        return request;
      });

      // Добавляем перехватчик для логирования ответа
      debugApi.interceptors.response.use(
        response => {
          console.log('Получен ответ:', {
            status: response.status,
            headers: response.headers,
            data: response.data
          });
          return response;
        },
        error => {
          console.error('Ошибка запроса:', {
            status: error.response?.status,
            data: error.response?.data,
            headers: error.response?.headers,
            message: error.message
          });
          return Promise.reject(error);
        }
      );

      const response = await debugApi.post('/api/auth/login', {
        username,
        password,
      });

      console.log('Получен ответ от сервера:', response.status);
      console.log('Данные ответа:', {
        hasAccessToken: !!response.data.access_token,
        hasRefreshToken: !!response.data.refresh_token,
        hasUser: !!response.data.user
      });

      if (!response.data.access_token || !response.data.refresh_token || !response.data.user) {
        console.error('Неполный ответ от сервера:', response.data);
        throw new Error('Неверный ответ от сервера');
      }

      const { access_token, refresh_token, user } = response.data;
      console.log('Сохранение токенов и данных пользователя');
      
      localStorage.setItem('token', access_token);
      localStorage.setItem('refreshToken', refresh_token);
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
      setUser(user);
      setIsAdmin(checkIsAdmin(user));
      console.log('Вход успешно завершен');
    } catch (error) {
      console.error('Ошибка входа:', error);
      if (axios.isAxiosError(error)) {
        console.error('Детали ошибки:', {
          status: error.response?.status,
          data: error.response?.data,
          headers: error.response?.headers,
          config: {
            url: error.config?.url,
            method: error.config?.method,
            headers: error.config?.headers,
            data: error.config?.data
          }
        });
        if (error.response?.status === 401) {
          throw new Error('Неверное имя пользователя или пароль');
        }
        throw new Error(error.response?.data?.message || 'Ошибка входа');
      }
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('refreshToken');
    delete api.defaults.headers.common['Authorization'];
    setUser(null);
    setIsAdmin(false);
  };

  const register = async (username: string, email: string, password: string) => {
    try {
      const response = await api.post<{ access_token: string; refresh_token: string; user: User }>(API_ENDPOINTS.AUTH.REGISTER, {
        username,
        email,
        password,
      });

      if (!response.data.access_token || !response.data.refresh_token || !response.data.user) {
        throw new Error('Неверный ответ от сервера');
      }

      const { access_token, refresh_token, user } = response.data;
      
      localStorage.setItem('token', access_token);
      localStorage.setItem('refreshToken', refresh_token);
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
      setUser(user);
      setIsAdmin(checkIsAdmin(user));
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 409) {
          throw new Error('Пользователь с таким именем или email уже существует');
        }
        throw new Error(error.response?.data?.message || 'Ошибка регистрации');
      }
      throw error;
    }
  };

  if (isLoading) {
    return <div>Загрузка...</div>;
  }

  return (
    <AuthContext.Provider value={{ user, isAdmin, login, logout, register }}>
      {children}
    </AuthContext.Provider>
  );
}; 