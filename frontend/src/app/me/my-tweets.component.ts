import { Component, OnInit } from '@angular/core';

import { Tweet } from '../tweets';
import { UserService } from '../users';
import { Store } from '../shared';


@Component({
  template: `
    <create-tweet></create-tweet>
    <tweets [tweets]="tweets"></tweets>
  `
})
export class MyTweetsComponent implements OnInit {
  tweets: Tweet[] = []

  constructor(
    private _userService: UserService,
    private _store: Store
  ) {

  }

  ngOnInit(): void {
    this._userService.getTweets()
      .subscribe((tweets: any) => this.tweets = tweets)

    this._store.changes("my_tweets")
      .subscribe((tweets: any) => this.tweets = tweets)
  }
}
