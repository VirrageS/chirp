import { Component } from '@angular/core';

import { Tweet } from '../shared';
import { Store } from '../store';


@Component({
  selector: 'me',
  templateUrl: './me.component.html',
  styleUrls: ['./me.component.scss']
})
export class MeComponent {
  following: number = 0
  followers: number = 0
  tweets: number = 0

  constructor(private _store: Store) {
    this._store.changes.pluck("my_tweets")
      .subscribe((tweets: any) => this.tweets = tweets.length)

    this._store.changes.pluck("user")
      .subscribe((user: any) => {
        // this.following = user.followee_count
        // this.followers = user.follower_count
      })
  }
}
