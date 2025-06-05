import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Button, Input, List, message, Spin, Popconfirm, Space } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useAuth } from '../hooks/useAuth';
import { Topic, Comment } from '../types/topic';
import { topicApi } from '../api/topic';
import { useDispatch } from 'react-redux';
import { setTopics } from '../store/slices/topicsSlice';

const { TextArea } = Input;

const TopicPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user, isAdmin } = useAuth();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [newComment, setNewComment] = useState('');
  const dispatch = useDispatch();

  useEffect(() => {
    fetchTopic();
  }, [id]);

  const fetchTopic = async () => {
    try {
      setLoading(true);
      const response = await topicApi.getTopic(Number(id));
      console.log('Topic response:', response.data);
      setTopic(response.data);
      setComments(response.data.comments || []);
      
      // Добавляем логирование для отладки
      console.log('Topic loaded:', response.data);
      console.log('Current user:', user);
      console.log('Is admin:', isAdmin);
      if (user && response.data) {
        console.log('User ID:', user.id, 'Type:', typeof user.id);
        console.log('Author ID:', response.data.author_id, 'Type:', typeof response.data.author_id);
        console.log('ID comparison:', String(user.id) === String(response.data.author_id));
      }
    } catch (error) {
      console.error('Ошибка при загрузке темы:', error);
      message.error('Не удалось загрузить тему');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = () => {
    navigate(`/topics/${id}/edit`);
  };

  const handleDelete = async () => {
    try {
      await topicApi.deleteTopic(Number(id));
      message.success('Тема успешно удалена');
      navigate('/');
    } catch (error) {
      console.error('Ошибка при удалении темы:', error);
      message.error('Не удалось удалить тему');
    }
  };

  const handleAddComment = async () => {
    if (!newComment.trim()) {
      message.warning('Комментарий не может быть пустым');
      return;
    }

    try {
      await topicApi.createComment(Number(id), {
        content: newComment
      });
      setNewComment('');
      message.success('Комментарий добавлен');
      
      // Обновляем данные темы и комментарии
      const topicResponse = await topicApi.getTopic(Number(id));
      setTopic(topicResponse.data);
      setComments(topicResponse.data.comments || []);
      
      // Обновляем список тем на главной странице
      const topicsResponse = await topicApi.getAllTopics();
      dispatch(setTopics(topicsResponse.data));
    } catch (error) {
      console.error('Ошибка при добавлении комментария:', error);
      message.error('Не удалось добавить комментарий');
    }
  };

  const handleDeleteComment = async (commentId: number) => {
    const comment = comments.find(c => c.id === commentId);
    if (!comment) return;

    // Проверяем права на удаление только при нажатии кнопки удаления
    const canDeleteComment = user && (isAdmin || String(user.id) === String(comment.author_id));
    if (!canDeleteComment) {
      message.error('У вас нет прав на удаление этого комментария');
      return;
    }

    try {
      await topicApi.deleteComment(Number(id), commentId);
      message.success('Комментарий удален');
      
      // Обновляем данные темы и комментарии
      const topicResponse = await topicApi.getTopic(Number(id));
      setTopic(topicResponse.data);
      setComments(topicResponse.data.comments || []);
      
      // Обновляем список тем на главной странице
      const topicsResponse = await topicApi.getAllTopics();
      dispatch(setTopics(topicsResponse.data));
    } catch (error) {
      console.error('Ошибка при удалении комментария:', error);
      message.error('Не удалось удалить комментарий');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  if (!topic) {
    return <div className="text-center p-4">Тема не найдена</div>;
  }

  // Добавляем логирование для отладки
  console.log('Topic in render:', topic);
  console.log('Author:', topic.author);
  console.log('Author username:', topic.author?.username);

  // Правильное сравнение ID пользователя и автора темы
  const canEditTopic = user && (isAdmin || String(user.id) === String(topic.author_id));
  
  // Добавляем логирование перед рендерингом
  console.log('Render check:');
  console.log('User:', user);
  console.log('Is admin:', isAdmin);
  console.log('Topic author ID:', topic.author_id);
  console.log('Can edit topic:', canEditTopic);

  return (
    <div className="container mx-auto px-4 py-8">
      <Card
        title={topic.title}
        extra={
          canEditTopic && (
            <Space>
              <Button
                type="text"
                icon={<EditOutlined />}
                onClick={handleEdit}
              >
                Редактировать
              </Button>
              <Popconfirm
                title="Удалить тему?"
                description="Вы уверены, что хотите удалить эту тему?"
                onConfirm={handleDelete}
                okText="Да"
                cancelText="Нет"
              >
                <Button
                  type="text"
                  danger
                  icon={<DeleteOutlined />}
                >
                  Удалить
                </Button>
              </Popconfirm>
            </Space>
          )
        }
      >
        <div className="topic-content">
          <p className="whitespace-pre-wrap">{topic.content}</p>
          <div className="topic-meta mt-4 text-gray-500">
            <span className="mr-4">Автор: {topic.author?.username || 'Unknown User'}</span>
            <span>Создано: {new Date(topic.created_at).toLocaleDateString()}</span>
          </div>
        </div>
      </Card>

      <div className="mt-8">
        <h3 className="text-xl font-bold mb-4">Комментарии</h3>
        {user ? (
          <div className="mb-4">
            <TextArea
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              placeholder="Написать комментарий..."
              rows={4}
              className="mb-2"
            />
            <Button type="primary" onClick={handleAddComment}>
              Отправить
            </Button>
          </div>
        ) : (
          <p className="text-gray-500 mb-4">
            Войдите в систему, чтобы оставить комментарий
          </p>
        )}

        <List
          dataSource={comments}
          renderItem={comment => {
            const canDeleteComment = user && (isAdmin || String(user.id) === String(comment.author_id));
            return (
              <List.Item
                actions={
                  canDeleteComment ? [
                    <Popconfirm
                      key="delete"
                      title="Удалить комментарий?"
                      description="Вы уверены, что хотите удалить этот комментарий?"
                      onConfirm={() => handleDeleteComment(comment.id)}
                      okText="Да"
                      cancelText="Нет"
                    >
                      <Button
                        type="text"
                        danger
                        icon={<DeleteOutlined />}
                      >
                        Удалить
                      </Button>
                    </Popconfirm>
                  ] : []
                }
              >
                <List.Item.Meta
                  title={
                    <div className="flex items-center">
                      <span className="font-medium text-blue-600">{comment.author.username}</span>
                      <span className="text-gray-500 text-sm ml-2">
                        {new Date(comment.created_at).toLocaleString()}
                      </span>
                    </div>
                  }
                  description={
                    <div className="mt-2 whitespace-pre-wrap text-gray-800">
                      {comment.content}
                    </div>
                  }
                />
              </List.Item>
            );
          }}
        />
      </div>
    </div>
  );
};

export default TopicPage; 