import { Component, Input } from '@angular/core';

import { Tweet, TweetService } from '../shared';


@Component({
  selector: 'create-tweet',
  templateUrl: './create-tweet.component.html',
  styleUrls: ['./create-tweet.component.scss']
})
export class CreateTweetComponent {
  tweet: Tweet

  constructor(
    private _tweetService: TweetService
  ) {
    this.tweet = {
      content: ""
    }
  }

  onSubmit(): void {
    this._tweetService.createTweet(this.tweet)
      .subscribe(
        result => {
          // TODO: message
          this.tweet = {
            content: ""
          }
        },
        error => {
          // TODO: message
        }
      )
  }
}
