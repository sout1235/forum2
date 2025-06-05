import axiosInstance from '../config/axios';
import { Topic, CreateTopicDto, Comment, CreateCommentDto } from '../types/topic';

export const topicApi = {
  getAllTopics: () => 
    axiosInstance.get<Topic[]>('/topics'),

  getTopic: (id: number) =>
    axiosInstance.get<Topic>(`/topics/${id}`),

  createTopic: (data: CreateTopicDto) =>
    axiosInstance.post<Topic>('/topics', data),

  updateTopic: (id: number, data: Partial<CreateTopicDto>) =>
    axiosInstance.put<Topic>(`/topics/${id}`, data),

  deleteTopic: (id: number) =>
    axiosInstance.delete(`/topics/${id}`),

  createComment: (topicId: number, data: { content: string }) =>
    axiosInstance.post<Comment>(`/topics/${topicId}/comments`, data),

  deleteComment: (topicId: number, commentId: number) =>
    axiosInstance.delete(`/topics/${topicId}/comments/${commentId}`)
}; 