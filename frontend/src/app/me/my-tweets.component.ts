import { Component, OnInit } from '@angular/core';

import { UserService, Tweet } from '../shared';
import { Store } from '../store';


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
    this._userService.getFeed()
      .subscribe((tweets: any) => this.tweets = tweets)

    this._store.changes.pluck("tweets")
      .subscribe((tweets: any) => this.tweets = tweets)
  }
}
