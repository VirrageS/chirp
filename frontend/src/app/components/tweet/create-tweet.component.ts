import { Component, Input } from '@angular/core';

import { Tweet, TweetService, UserService } from '../../shared';


@Component({
  selector: 'create-tweet',
  templateUrl: './create-tweet.component.html',
  styleUrls: ['./create-tweet.component.scss']
})
export class CreateTweetComponent {
  tweet: Tweet

  constructor(
    private _tweetService: TweetService,
    private _userService: UserService
  ) {
    // TODO: can we avoid initalization?
    this.tweet = {
      content: ""
    }
  }

  onSubmit(): void {
    this._tweetService.createTweet(this.tweet)
  }
}
