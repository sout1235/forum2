export interface User {
  id: number | string;
  username: string;
  email: string;
  role: string;
  avatar?: string;
  isAdmin?: boolean;
} 