import { Component, Input } from '@angular/core';

import { Tweet, TweetService } from '../../shared';


import '../../../../public/scss/abstract.scss';


@Component({
  selector: 'tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.scss']
})
export class TweetComponent {
  @Input() tweet: Tweet

  constructor(private _tweetService: TweetService) {
  }

  private _follow() {
    this.tweet.author.following = true
  }

  private _toggleLike() {
    this.tweet.liked = !this.tweet.liked

    let toggleFunc = this._tweetService.like(this.tweet.id)
    if (this.tweet.liked) {
      toggleFunc = this._tweetService.unlike(this.tweet.id)
    }

    toggleFunc
      .subscribe(
        result => {
          // this.tweet = result
        },
        error => {}
      )
  }

  private _retweet() {
    this.tweet.retweeted = true
    this._tweetService.retweet(this.tweet.id)
      .subscribe(
        result => {
          // this.tweet = result
        },
        error => {}
      )
  }
}
