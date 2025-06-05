import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import { Layout, Menu, Button } from 'antd';
import { UserOutlined, LogoutOutlined } from '@ant-design/icons';
import { Provider } from 'react-redux';
import { store } from './store';
import { TopicList } from './components/TopicList';
import Chat from './components/Chat';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import NewTopicPage from './pages/NewTopicPage';
import TopicPage from './pages/TopicPage';
import EditTopicPage from './pages/EditTopicPage';
import { useAuth } from './hooks/useAuth';
import { AuthProvider } from './contexts/AuthContext';
import './App.css';

const { Header, Content, Footer } = Layout;

const AppContent: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <Layout className="min-h-screen">
      <Header className="flex justify-between items-center">
        <Menu theme="dark" mode="horizontal" className="flex-1">
          <Menu.Item key="home">
            <Link to="/">Главная</Link>
          </Menu.Item>
          {user ? (
            <>
              <Menu.Item key="chat">
                <Link to="/chat">Чат</Link>
              </Menu.Item>
              <Menu.Item key="logout" className="ml-auto">
                <Button type="link" onClick={logout} className="text-white">
                  <LogoutOutlined /> Выйти
                </Button>
              </Menu.Item>
            </>
          ) : (
            <>
              <Menu.Item key="login" className="ml-auto">
                <Link to="/login">Войти</Link>
              </Menu.Item>
              <Menu.Item key="register">
                <Link to="/register">Регистрация</Link>
              </Menu.Item>
            </>
          )}
        </Menu>
      </Header>

      <Content className="p-6">
        <Routes>
          <Route path="/" element={<TopicList />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/topics/create" element={<NewTopicPage />} />
          <Route path="/topics/:id" element={<TopicPage />} />
          <Route path="/topics/:id/edit" element={<EditTopicPage />} />
          <Route path="/chat" element={<Chat />} />
        </Routes>
      </Content>

      <Footer className="text-center">
        Forum ©{new Date().getFullYear()}
      </Footer>
    </Layout>
  );
};

const App: React.FC = () => {
  return (
    <Provider store={store}>
      <Router>
        <AuthProvider>
          <AppContent />
        </AuthProvider>
      </Router>
    </Provider>
  );
};

export default App;
