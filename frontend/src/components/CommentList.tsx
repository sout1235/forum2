import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../store';
import { addComment } from '../store/slices/commentsSlice';
import { Comment } from '../types';

interface CommentListProps {
  topicId: number;
}

const CommentList: React.FC<CommentListProps> = ({ topicId }) => {
  const dispatch = useDispatch();
  const { comments, loading, error } = useSelector((state: RootState) => state.comments);
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);
  const [content, setContent] = useState('');
  const [commentError, setCommentError] = useState('');

  useEffect(() => {
    const fetchComments = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/topics/${topicId}/comments`);
        const data = await response.json();
        dispatch({ type: 'comments/setComments', payload: data });
      } catch (error) {
        console.error('Error fetching comments:', error);
      }
    };

    fetchComments();
  }, [dispatch, topicId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setCommentError('');

    if (!content.trim()) {
      setCommentError('Комментарий не может быть пустым');
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/topics/${topicId}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ content }),
      });

      if (!response.ok) {
        throw new Error('Ошибка при создании комментария');
      }

      const data = await response.json();
      dispatch(addComment(data as Comment));
      setContent('');
    } catch (error) {
      setCommentError('Ошибка при создании комментария');
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-32">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
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
    <div className="space-y-6">
      {isAuthenticated && (
        <form onSubmit={handleSubmit} className="space-y-4">
          {commentError && (
            <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded">
              {commentError}
            </div>
          )}
          <div>
            <label htmlFor="comment" className="sr-only">
              Ваш комментарий
            </label>
            <textarea
              id="comment"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={3}
              className="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              placeholder="Написать комментарий..."
            />
          </div>
          <div className="flex justify-end">
            <button
              type="submit"
              className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              Отправить
            </button>
          </div>
        </form>
      )}

      <div className="space-y-4">
        {comments.map((comment) => (
          <div key={comment.id} className="bg-gray-50 rounded-lg p-4">
            <div className="flex items-center space-x-3">
              <img
                className="h-8 w-8 rounded-full"
                src={comment.author.avatar || 'https://via.placeholder.com/40'}
                alt=""
              />
              <div>
                <p className="text-sm font-medium text-gray-900">
                  {comment.author.username}
                </p>
                <p className="text-sm text-gray-500">
                  {new Date(comment.createdAt).toLocaleDateString()}
                </p>
              </div>
            </div>
            <div className="mt-2 text-sm text-gray-700">
              {comment.content}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CommentList; 