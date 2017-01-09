import { Component, Input } from '@angular/core';

import { Tweet, TweetService, User, UserService } from '../../shared';
import { Store } from '../../store';


@Component({
  selector: 'tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.scss']
})
export class TweetComponent {
  @Input() tweet: Tweet
  loggedUser: User

  constructor(
    private _tweetService: TweetService,
    private _userService: UserService,
    private _store: Store
  ) {
    this._store.changes.pluck("user")
      .subscribe((user: any) => this.loggedUser = user)
  }

  private _follow() {
    this.tweet.author.following = true

    // send real request
    this._userService.follow(this.tweet.author.id)
      .subscribe(author => this.tweet.author = author)
  }

  private _toggleLike() {
    this.tweet.liked = !this.tweet.liked

    let toggleFunc = this._tweetService.like(this.tweet.id)
    if (!this.tweet.liked) {
      toggleFunc = this._tweetService.unlike(this.tweet.id)
    }

    toggleFunc
      .subscribe(
        result => {
          this.tweet = result
        },
        error => {}
      )
  }

  private _retweet() {
    this.tweet.retweeted = true
    this._tweetService.retweet(this.tweet.id)
      .subscribe(
        result => {
          this.tweet = result
        },
        error => {}
      )
  }
}
