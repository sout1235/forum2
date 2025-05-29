import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../store';
import { setTopics } from '../store/slices/topicsSlice';
import TopicList from '../components/TopicList';
import { Topic } from '../types';

const Home: React.FC = () => {
  const dispatch = useDispatch();
  const { topics, loading, error } = useSelector((state: RootState) => state.topics);

  useEffect(() => {
    const fetchTopics = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/topics');
        const data = await response.json();
        dispatch(setTopics(data as Topic[]));
      } catch (error) {
        console.error('Error fetching topics:', error);
      }
    };

    fetchTopics();
  }, [dispatch]);

  return (
    <div className="container mx-auto px-4 py-8">
      <TopicList />
    </div>
  );
};

export default Home; 