export interface User {
  id: number;
  username: string;
  email: string;
  avatar?: string;
  isAdmin: boolean;
}

export type { Topic, Author, Category, Tag, CreateTopicDto } from './topic'; 