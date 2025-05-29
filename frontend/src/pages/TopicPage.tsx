import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../store';
import { setCurrentTopic } from '../store/slices/topicsSlice';
import TopicDetail from '../components/TopicDetail';
import { Topic } from '../types';

const TopicPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const dispatch = useDispatch();
  const { currentTopic, loading, error } = useSelector((state: RootState) => state.topics);

  useEffect(() => {
    const fetchTopic = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/topics/${id}`);
        const data = await response.json();
        dispatch(setCurrentTopic(data as Topic));
      } catch (error) {
        console.error('Error fetching topic:', error);
      }
    };

    fetchTopic();
  }, [dispatch, id]);

  return (
    <div className="container mx-auto px-4 py-8">
      <TopicDetail />
    </div>
  );
};

export default TopicPage; 