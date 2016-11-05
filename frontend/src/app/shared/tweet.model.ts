import { User } from './user.model'

export class Tweet {
  id: number
  author: User = new User()
  likes: number = 0
  retweets: number = 0
  liked: boolean = false
  retweeted: boolean = false
  created_at: string
  content: string = ""
}
