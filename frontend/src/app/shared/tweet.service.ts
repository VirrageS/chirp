import { Injectable } from '@angular/core';

import { ApiService } from './api.service';
import { StoreHelper } from './store-helper';
import { Tweet } from './tweet.model';


@Injectable()
export class TweetService {
  tweet_path: string = "/tweets"

  constructor(
    private _apiService: ApiService,
    private _storeHelper: StoreHelper
  ) {

  }

  createTweet(tweet: Tweet) {
    return this._apiService.post(this.tweet_path, tweet)
      .do((tweet: Tweet) => this._storeHelper.add("my_tweets", tweet))
  }

  like(tweet_id: number) {
    return this._apiService.post(this.tweet_path + "/" + tweet_id + "/like", {})
  }

  unlike(tweet_id: number) {
    return this._apiService.post(this.tweet_path + "/" + tweet_id + "/unlike", {})
  }
}
