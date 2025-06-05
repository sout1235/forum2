import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../store';
import { setTopics, setLoading, setError } from '../store/slices/topicsSlice';
import { TopicList } from '../components/TopicList';
import { topicApi } from '../api/topic';
import { message } from 'antd';

const Home: React.FC = () => {
  const dispatch = useDispatch();
  const { topics, loading, error } = useSelector((state: RootState) => state.topics);

    const fetchTopics = async () => {
      try {
      dispatch(setLoading(true));
      const response = await topicApi.getAllTopics();
      dispatch(setTopics(response.data));
      } catch (error) {
      console.error('Ошибка при загрузке тем:', error);
      const errorMessage = error instanceof Error ? error.message : 'Не удалось загрузить темы';
      dispatch(setError(errorMessage));
      message.error(errorMessage);
    } finally {
      dispatch(setLoading(false));
      }
    };

  useEffect(() => {
    fetchTopics();
  }, [dispatch]);

  if (error) {
    return <div className="text-center text-red-500 p-4">{error}</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <TopicList topics={topics} onTopicsChange={fetchTopics} />
    </div>
  );
};

export default Home; 