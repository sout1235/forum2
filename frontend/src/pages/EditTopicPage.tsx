import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Form, Input, Button, message, Spin } from 'antd';
import { useAuth } from '../hooks/useAuth';
import { topicApi } from '../api/topic';
import { useDispatch } from 'react-redux';
import { setTopics } from '../store/slices/topicsSlice';

const { TextArea } = Input;

const EditTopicPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user, isAdmin } = useAuth();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const dispatch = useDispatch();

  useEffect(() => {
    fetchTopic();
  }, [id]);

  const fetchTopic = async () => {
    try {
      setLoading(true);
      const response = await topicApi.getTopic(Number(id));
      const topic = response.data;

      // Проверяем права на редактирование
      if (!user || (!isAdmin && String(user.id) !== String(topic.author_id))) {
        message.error('У вас нет прав на редактирование этой темы');
        navigate(`/topics/${id}`);
        return;
      }

      form.setFieldsValue({
        title: topic.title,
        content: topic.content
      });
    } catch (error) {
      console.error('Ошибка при загрузке темы:', error);
      message.error('Не удалось загрузить тему');
      navigate('/');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (values: { title: string; content: string }) => {
    try {
      setSubmitting(true);
      await topicApi.updateTopic(Number(id), values);
      
      // Обновляем список тем
      const topicsResponse = await topicApi.getAllTopics();
      dispatch(setTopics(topicsResponse.data));
      
      message.success('Тема успешно обновлена');
      navigate(`/topics/${id}`);
    } catch (error) {
      console.error('Ошибка при обновлении темы:', error);
      message.error('Не удалось обновить тему');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <Card title="Редактировать тему">
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={{ title: '', content: '' }}
        >
          <Form.Item
            name="title"
            label="Заголовок"
            rules={[{ required: true, message: 'Пожалуйста, введите заголовок' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            name="content"
            label="Содержание"
            rules={[{ required: true, message: 'Пожалуйста, введите содержание' }]}
          >
            <TextArea rows={6} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting}>
              Сохранить
            </Button>
            <Button 
              className="ml-2" 
              onClick={() => navigate(`/topics/${id}`)}
            >
              Отмена
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default EditTopicPage; 