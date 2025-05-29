import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import CommentList from './CommentList';
import { Topic } from '../types';

const TopicDetail: React.FC = () => {
  const { currentTopic, loading, error } = useSelector((state: RootState) => state.topics);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center text-red-600 p-4">
        {error}
      </div>
    );
  }

  if (!currentTopic) {
    return (
      <div className="text-center text-gray-600 p-4">
        Тема не найдена
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-900">{currentTopic.title}</h1>
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
            {currentTopic.category.name}
          </span>
        </div>
        
        <div className="mt-4 flex items-center text-sm text-gray-500">
          <span className="flex items-center">
            <img
              className="h-8 w-8 rounded-full mr-2"
              src={currentTopic.author.avatar || 'https://via.placeholder.com/40'}
              alt=""
            />
            {currentTopic.author.username}
          </span>
          <span className="mx-2">•</span>
          <span>{new Date(currentTopic.createdAt).toLocaleDateString()}</span>
          <span className="mx-2">•</span>
          <span>{currentTopic.views} просмотров</span>
        </div>

        <div className="mt-6 prose max-w-none">
          {currentTopic.content}
        </div>

        <div className="mt-6 flex flex-wrap gap-2">
          {currentTopic.tags.map((tag) => (
            <span
              key={tag.id}
              className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
            >
              {tag.name}
            </span>
          ))}
        </div>
      </div>

      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Комментарии</h2>
        <CommentList topicId={currentTopic.id} />
      </div>
    </div>
  );
};

export default TopicDetail; 