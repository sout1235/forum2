import React, { useEffect } from 'react';
import { List, Card, Button, Space, Tag, Spin, Popconfirm, message } from 'antd';
import { EditOutlined, DeleteOutlined, EyeOutlined, PlusOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Topic } from '../types/topic';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../store';
import { topicApi } from '../api/topic';
import { setTopics } from '../store/slices/topicsSlice';

interface TopicListProps {
  topics?: Topic[];
  onEdit?: (id: number) => void;
  onDelete?: (id: number) => void;
  onTopicsChange?: () => void;
}

export const TopicList: React.FC<TopicListProps> = ({ 
  topics: propTopics, 
  onEdit, 
  onDelete,
  onTopicsChange 
}) => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { user, isAdmin } = useAuth();
  const { topics: storeTopics, loading } = useSelector((state: RootState) => state.topics);

  const topics = propTopics || storeTopics || [];

  // Загружаем темы при монтировании компонента
  useEffect(() => {
    if (!topics.length) {
      fetchTopics();
    }
  }, []);

  const fetchTopics = async () => {
    try {
      const response = await topicApi.getAllTopics();
      dispatch(setTopics(response.data));
    } catch (error) {
      console.error('Ошибка при загрузке тем:', error);
      message.error('Не удалось загрузить темы');
    }
  };

  const handleEdit = (id: number) => {
    if (onEdit) {
      onEdit(id);
    } else {
      navigate(`/topics/${id}/edit`);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await topicApi.deleteTopic(id);
      message.success('Тема успешно удалена');
      
      // Обновляем список тем
      if (onTopicsChange) {
        onTopicsChange();
      } else {
        const response = await topicApi.getAllTopics();
        dispatch(setTopics(response.data));
      }
    } catch (error) {
      console.error('Ошибка при удалении темы:', error);
      message.error('Не удалось удалить тему');
    }
  };

  const handleView = (id: number) => {
    navigate(`/topics/${id}`);
  };

  const handleCreate = () => {
    navigate('/topics/create');
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div className="topic-list">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Темы</h2>
        {user && (
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={handleCreate}
          >
            Создать тему
          </Button>
        )}
      </div>

      <List
        grid={{ gutter: 16, column: 1 }}
        dataSource={topics}
        renderItem={(topic: Topic) => (
          <List.Item
            actions={[
              <Button
                key="view"
                type="text"
                icon={<EyeOutlined />}
                onClick={() => handleView(topic.id)}
              >
                Просмотр
              </Button>
            ]}
          >
            <Card
              hoverable
              onClick={() => handleView(topic.id)}
              title={topic.title}
              style={{ width: '100%' }}
            >
              <div className="topic-content">
                <p className="whitespace-pre-wrap">{topic.content}</p>
                <div className="topic-meta mt-4 text-gray-500">
                  <span className="mr-4">Автор: {topic.author?.username || 'Unknown User'}</span>
                  <span>Создано: {new Date(topic.created_at).toLocaleDateString()}</span>
                </div>
                <div className="topic-tags mt-2">
                  {topic.tags?.map(tag => (
                    <Tag key={tag.id} color="blue">
                      {tag.name}
                    </Tag>
                  ))}
                </div>
              </div>
            </Card>
          </List.Item>
        )}
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          showTotal: total => `Всего ${total} тем`,
        }}
      />
    </div>
  );
}; 