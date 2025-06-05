export interface Author {
  id: number;
  username: string;
  avatar?: string;
}

export interface Category {
  id: number;
  name: string;
  description?: string;
}

export interface Tag {
  id: number;
  name: string;
}

export interface Topic {
  id: number;
  title: string;
  content: string;
  author_id: number;
  author: Author;
  category_id: number;
  category: Category;
  views: number;
  tags: Tag[];
  created_at: string;
  updated_at: string;
  comment_count: number;
  comments?: Comment[];
}

export interface Comment {
  id: number;
  content: string;
  author_id: number;
  author: Author;
  topic_id: number;
  created_at: string;
  updated_at: string;
}

export interface CreateTopicDto {
  title: string;
  content: string;
  category_id: number;
  tags?: number[];
}

export interface CreateCommentDto {
  content: string;
  topic_id: number;
} 