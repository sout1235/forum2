import React from 'react';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import { Topic } from '../types';

const TopicList: React.FC = () => {
  const { topics, loading, error } = useSelector((state: RootState) => state.topics);

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

  return (
    <div className="space-y-4">
      {topics.map((topic: Topic) => (
        <div
          key={topic.id}
          className="bg-white shadow rounded-lg p-6 hover:shadow-md transition-shadow"
        >
          <div className="flex items-center justify-between">
            <div className="flex-1">
              <Link
                to={`/topics/${topic.id}`}
                className="text-xl font-semibold text-gray-900 hover:text-indigo-600"
              >
                {topic.title}
              </Link>
              <div className="mt-2 flex items-center text-sm text-gray-500">
                <span className="flex items-center">
                  <img
                    className="h-5 w-5 rounded-full mr-2"
                    src={topic.author.avatar || 'https://via.placeholder.com/40'}
                    alt=""
                  />
                  {topic.author.username}
                </span>
                <span className="mx-2">•</span>
                <span>{new Date(topic.createdAt).toLocaleDateString()}</span>
                <span className="mx-2">•</span>
                <span>{topic.views} просмотров</span>
                <span className="mx-2">•</span>
                <span>{topic.comments} комментариев</span>
              </div>
            </div>
            <div className="ml-4">
              <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
                {topic.category.name}
              </span>
            </div>
          </div>
          <div className="mt-4 flex flex-wrap gap-2">
            {topic.tags.map((tag) => (
              <span
                key={tag.id}
                className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
              >
                {tag.name}
              </span>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

export default TopicList; 