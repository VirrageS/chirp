import { Component, Input } from '@angular/core';

import { AlertType } from '../../core/alerts';
import { StoreHelper } from '../../shared';
import { Tweet, TweetService } from '../shared';


@Component({
  selector: 'create-tweet',
  templateUrl: './create-tweet.component.html',
  styleUrls: ['./create-tweet.component.scss']
})
export class CreateTweetComponent {
  private tweet: Tweet = {
    content: ""
  }

  constructor(
    private storeHelper: StoreHelper,
    private tweetService: TweetService,
  ) {}

  onSubmit(): void {
    this.tweetService.createTweet(this.tweet)
      .subscribe(
        result => {
          this.storeHelper.add("alerts", {message: "Tweet has been added successfully", type: AlertType.success});
          this.tweet = {content: ""}
        },
        error => {
          this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
        }
      )
  }
}
