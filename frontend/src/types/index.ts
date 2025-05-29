export interface User {
  id: number;
  username: string;
  email: string;
  avatar?: string;
  role: string;
}

export interface Topic {
  id: number;
  title: string;
  content: string;
  authorId: number;
  author: User;
  categoryId: number;
  category: {
    id: number;
    name: string;
  };
  views: number;
  comments: number;
  tags: Array<{
    id: number;
    name: string;
  }>;
  createdAt: string;
  updatedAt: string;
}

export interface Comment {
  id: number;
  content: string;
  authorId: number;
  author: User;
  topicId: number;
  parentId?: number;
  replies: Comment[];
  likes: number;
  createdAt: string;
  updatedAt: string;
}

export interface Category {
  id: number;
  name: string;
  description?: string;
} 