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
  private tweets: Array<Tweet> = []

  constructor(
    private userService: UserService,
    private store: Store
  ) {}

  ngOnInit(): void {
    this.userService.getTweets()
      .subscribe((tweets: Array<Tweet>) => this.tweets = tweets)
    this.store.changes("my_tweets")
      .subscribe((tweets: Array<Tweet>) => this.tweets = tweets)
  }
}
