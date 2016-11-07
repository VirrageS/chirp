import { Component } from '@angular/core';

import { Tweet, TweetService } from '../../shared';


@Component({
  selector: 'create-tweet',
  templateUrl: './create-tweet.component.html',
  styleUrls: ['./tweet.component.scss']
})
export class CreateTweetComponent {
  tweet?: Tweet

  constructor(private _tweetService: TweetService) {
    
  }

  onSubmit(): void {
    this._tweetService.createTweet(this.tweet)
  }
}
