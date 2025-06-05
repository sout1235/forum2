import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { topicApi } from '../api/topic';
import { useDispatch } from 'react-redux';
import { setTopics } from '../store/slices/topicsSlice';

const { TextArea } = Input;

const NewTopicPage: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values: any) => {
    try {
      setLoading(true);
      await topicApi.createTopic(values);
      message.success('Тема успешно создана');
      
      // Обновляем список тем
      const response = await topicApi.getAllTopics();
      dispatch(setTopics(response.data));
      
      navigate('/');
    } catch (error) {
      console.error('Ошибка при создании темы:', error);
      message.error('Не удалось создать тему');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Создать новую тему</h1>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        className="max-w-2xl"
      >
        <Form.Item
          name="title"
          label="Заголовок"
          rules={[
            { required: true, message: 'Пожалуйста, введите заголовок' },
            { min: 3, message: 'Заголовок должен содержать минимум 3 символа' },
          ]}
        >
          <Input placeholder="Введите заголовок темы" />
        </Form.Item>

        <Form.Item
          name="content"
          label="Содержание"
          rules={[
            { required: true, message: 'Пожалуйста, введите содержание' },
            { min: 10, message: 'Содержание должно содержать минимум 10 символов' },
          ]}
        >
          <TextArea
            rows={6}
            placeholder="Введите содержание темы"
          />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading}>
            Создать тему
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default NewTopicPage; 