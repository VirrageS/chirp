import { Component, Input } from '@angular/core';

import { Tweet, TweetService } from '../../shared';


@Component({
  selector: 'tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.scss']
})
export class TweetComponent {
  @Input() tweet: Tweet

  constructor(private _tweetService: TweetService) {
  }

  private _like() {
    this._tweetService.like(this.tweet.id)
      .subscribe(
        result => {
          this.tweet = result
        },
        error => {}
      )
  }

  private _unlike() {
    this._tweetService.unlike(this.tweet.id)
      .subscribe(
        result => {
          this.tweet = result
        },
        error => {}
      )
  }
}
