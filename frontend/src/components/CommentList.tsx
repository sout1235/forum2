import React, { useState } from 'react';
import { List, Button, Input, message, Popconfirm } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';
import { useAuth } from '../hooks/useAuth';
import { Comment } from '../types/topic';
import { topicApi } from '../api/topic';

const { TextArea } = Input;

interface CommentListProps {
  topicId: number;
  comments: Comment[];
  onCommentsChange: () => void;
}

export const CommentList: React.FC<CommentListProps> = ({
  topicId,
  comments,
  onCommentsChange,
}) => {
  const { user, isAdmin } = useAuth();
  const [newComment, setNewComment] = useState('');

  const handleAddComment = async () => {
    if (!newComment.trim()) {
      message.warning('Комментарий не может быть пустым');
      return;
    }

    try {
      await topicApi.createComment(topicId, {
        content: newComment
      });
      setNewComment('');
      onCommentsChange();
      message.success('Комментарий добавлен');
    } catch (error) {
      console.error('Ошибка при добавлении комментария:', error);
      message.error('Не удалось добавить комментарий');
    }
  };

  const handleDeleteComment = async (commentId: number) => {
    try {
      await topicApi.deleteComment(topicId, commentId);
      onCommentsChange();
      message.success('Комментарий удален');
    } catch (error) {
      console.error('Ошибка при удалении комментария:', error);
      message.error('Не удалось удалить комментарий');
    }
  };

  return (
    <div className="comment-list">
      {user ? (
        <div className="mb-4">
          <TextArea
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Написать комментарий..."
            rows={4}
            className="mb-2"
          />
          <Button type="primary" onClick={handleAddComment}>
            Отправить
          </Button>
        </div>
      ) : (
        <p className="text-gray-500 mb-4">
          Войдите в систему, чтобы оставить комментарий
        </p>
      )}

      <List
        dataSource={comments}
        renderItem={(comment) => {
          const canDeleteComment = user && (isAdmin || user.id === comment.author_id);

          return (
            <List.Item
              actions={
                canDeleteComment ? [
                  <Popconfirm
                    key="delete"
                    title="Удалить комментарий?"
                    description="Вы уверены, что хотите удалить этот комментарий?"
                    onConfirm={() => handleDeleteComment(comment.id)}
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
                ] : []
              }
            >
              <List.Item.Meta
                title={
                  <div className="flex items-center">
                    <span className="font-medium">{comment.author.username}</span>
                    <span className="text-gray-500 text-sm ml-2">
                      {new Date(comment.created_at).toLocaleString()}
                    </span>
                  </div>
                }
                description={
                  <div className="whitespace-pre-wrap">{comment.content}</div>
                }
              />
            </List.Item>
          );
        }}
      />
    </div>
  );
}; 