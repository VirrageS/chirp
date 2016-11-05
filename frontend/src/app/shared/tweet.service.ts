import { Injectable } from '@angular/core';
import { ApiService } from './api.service';

@Injectable()
export class TweetService {
  tweet_path: string = "/tweet"

  constructor(private apiService: ApiService) {
    // TODO: get user_id
  }

  like(tweet_id: number) {
    return this.apiService.post(this.tweet_path + "/" + tweet_id + "/like", {});
  }

  unlike(tweet_id: number) {
    return this.apiService.post(this.tweet_path + "/" + tweet_id + "/unlike", {});
  }
}
