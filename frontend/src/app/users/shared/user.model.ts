export interface User {
  id?: number;
  name?: string;
  username?: string;
  email?: string;
  password?: string;
  created_at?: string;

  following?: boolean;
  follower_count?: number;
  followee_count?: number;
}
