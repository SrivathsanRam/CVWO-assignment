export interface User {
  id: number;
  username: string;
  created_at: string;
}

export interface Topic {
  id: number;
  title: string;
  description: string;
  user_id: number;
  username?: string; // Optional: include username
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  topic_id: number;
  user_id: number;
  username?: string; // Optional: include username
  topic_name?: string; // Optional: include topic name
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  content: string;
  post_id: number;
  user_id: number;
  username?: string; // Optional: include username
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  user: User;
  message: string;
}

export interface ApiError {
  error: string;
}

// Request type
export interface CreateTopicRequest {
  title: string;
  description: string;
}

export interface UpdateTopicRequest {
  title: string;
  description: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  topic_id: number;
}

export interface UpdatePostRequest {
  title: string;
  content: string;
}

export interface CreateCommentRequest {
  content: string;
  post_id: number;
}

export interface UpdateCommentRequest {
  content: string;
}