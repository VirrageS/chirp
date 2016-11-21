export interface User {
  id?: number;
  name?: string;
  username?: string;
  email?: string;
  password?: string;
  created_at?: string;

  following?: boolean;
}
