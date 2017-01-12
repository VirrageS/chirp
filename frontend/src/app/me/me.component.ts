import { Component } from '@angular/core';

import { Tweet, User, UserService } from '../shared';
import { Store } from '../store';


@Component({
  selector: 'me',
  templateUrl: './me.component.html',
  styleUrls: ['./me.component.scss']
})
export class MeComponent {
  following_count: number = 0
  follower_count: number = 0
  tweet_count: number = 0

  constructor(private _userService: UserService, private _store: Store) {
    this._userService.getTweets()
      .subscribe((tweets: Array<Tweet>) => this.tweet_count = tweets.length)
    this._store.changes.pluck("my_tweets")
      .subscribe((tweets: Array<Tweet>) => this.tweet_count = tweets.length)

    this._userService.getFollowers()
      .subscribe((users: Array<User>) => this.follower_count = users.length)
    this._store.changes.pluck("my_followers")
      .subscribe((users: Array<User>) => this.follower_count = users.length)

    this._userService.getFollowing()
      .subscribe((users: Array<User>) => this.following_count = users.length)
    this._store.changes.pluck("my_following")
      .subscribe((users: Array<User>) => this.following_count = users.length)
  }
}
