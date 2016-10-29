import { User } from './user.model'

export class Tweet {
  id: number
  author: User
  like_count: number
  retweet_count: number
  created_at: string
  content: string
}
