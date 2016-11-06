import { User } from './user.model'

export interface Tweet {
  id?: number
  author: User
  likes?: number
  retweets?: number
  liked?: boolean
  retweeted?: boolean
  created_at?: string
  content: string
}
