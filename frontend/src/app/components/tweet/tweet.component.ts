import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Tweet, TweetService, User, UserService } from '../../shared';
import { Store } from '../../store';
import * as _ from 'lodash';


@Component({
  selector: 'tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.scss']
})
export class TweetComponent {
  @Input() tweet: Tweet
  @Output() tweetChange = new EventEmitter()
  loggedUser: User

  constructor(
    private _tweetService: TweetService,
    private _userService: UserService,
    private _store: Store
  ) {
    this._store.changes.pluck("user")
      .subscribe((user: any) => this.loggedUser = user)
  }

  private _toggleFollow() {
    this.tweet.author.following = !this.tweet.author.following

    // send real request
    let toggleFunc = this._userService.follow(this.tweet.author.id)
    if (!this.tweet.author.following) {
      toggleFunc = this._userService.unfollow(this.tweet.author.id)
    }

    this.tweetChange.emit(this.tweet)
    toggleFunc
      .subscribe(author => {
        _.assign(this.tweet.author, author)
        this.tweetChange.emit(this.tweet)
      })
  }

  private _toggleLike() {
    this.tweet.liked = !this.tweet.liked

    let toggleFunc = this._tweetService.like(this.tweet.id)
    if (!this.tweet.liked) {
      this.tweet.like_count -= 1
      toggleFunc = this._tweetService.unlike(this.tweet.id)
    } else {
      this.tweet.like_count += 1
    }

    toggleFunc
      .subscribe(
        result => _.assign(this.tweet, result),
        error => {}
      )
  }

  // TODO:
  // private _retweet() {
  //   this.tweet.retweeted = true
  //   this._tweetService.retweet(this.tweet.id)
  //     .subscribe(
  //       result => _.assign(this.tweet, result),
  //       error => {}
  //     )
  // }
}
