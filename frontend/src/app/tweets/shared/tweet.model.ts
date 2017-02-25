import { User } from '../../users/shared/user.model'

export interface Tweet {
  id?: number
  author?: User
  like_count?: number
  retweet_count?: number
  liked?: boolean
  retweeted?: boolean
  created_at?: string
  content: string
}
