import React, { useEffect, useState } from 'react';
import { Card, Typography, Spin, message, Button, Space, Popconfirm } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { topicApi } from '../api/topic';
import { Topic, Comment } from '../types/topic';
import { CommentList } from './CommentList';
import { useAuth } from '../hooks/useAuth';

const { Title, Paragraph } = Typography;

const TopicDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user, isAdmin } = useAuth();
  const [currentTopic, setCurrentTopic] = useState<Topic | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTopic();
  }, [id]);

  const fetchTopic = async () => {
    try {
      setLoading(true);
      const response = await topicApi.getTopic(Number(id));
      setCurrentTopic(response.data);
      setComments(response.data.comments || []);
      
      // Добавляем логирование после загрузки темы
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

  const handleCommentsChange = () => {
    fetchTopic();
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

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  if (!currentTopic) {
    return <div className="text-center p-4">Тема не найдена</div>;
  }

  // Проверяем, является ли пользователь автором темы или администратором
  const canEditTopic = user && (isAdmin || String(user.id) === String(currentTopic.author_id));
  
  // Добавляем логирование перед рендерингом
  console.log('Render check:');
  console.log('User:', user);
  console.log('Is admin:', isAdmin);
  console.log('Topic author ID:', currentTopic.author_id);
  console.log('Can edit topic:', canEditTopic);

  return (
    <div className="container mx-auto px-4 py-8">
      <Card
        title={
          <div className="flex justify-between items-center">
            <Title level={2} style={{ margin: 0 }}>{currentTopic.title}</Title>
            {canEditTopic && (
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
            )}
        </div>
        }
      >
        <Paragraph className="whitespace-pre-wrap">{currentTopic.content}</Paragraph>
        <div className="mt-4 text-gray-500">
          <span className="mr-4">Автор: {currentTopic.author?.username || 'Unknown User'}</span>
          <span className="mr-4">Создано: {new Date(currentTopic.created_at).toLocaleDateString()}</span>
          <span className="mr-4">Комментарии: {currentTopic.comment_count}</span>
          <span>Просмотры: {currentTopic.views}</span>
        </div>
      </Card>

      <Card className="mt-4">
        <Title level={3}>Комментарии ({currentTopic.comment_count})</Title>
        <CommentList 
          topicId={currentTopic.id} 
          comments={comments}
          onCommentsChange={handleCommentsChange}
        />
      </Card>
    </div>
  );
};

export default TopicDetail; 