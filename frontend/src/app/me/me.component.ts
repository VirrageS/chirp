import { Component } from '@angular/core';

import { Tweet } from '../shared';
import { Store } from '../store';


@Component({
  selector: 'me',
  templateUrl: './me.component.html',
  styleUrls: ['./me.component.scss']
})
export class MeComponent {
  following: number = 384
  followers: number = 2934890
  tweets: number = 0

  constructor(private _store: Store) {
    this._store.changes.pluck("my_tweets")
      .subscribe((tweets: any) => this.tweets = tweets.length)
  }
}
