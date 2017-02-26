import { Component, OnInit } from '@angular/core';

import { Tweet } from '../tweets';
import { User, UserService } from '../users';
import { Store } from '../shared';


@Component({
  selector: 'me',
  templateUrl: './me.component.html',
  styleUrls: ['./me.component.scss']
})
export class MeComponent implements OnInit {
  private following_count: number = 0;
  private follower_count: number = 0;
  private tweet_count: number = 0;

  constructor(
    private userService: UserService,
    private store: Store
  ) {}

  ngOnInit() {
    this.userService.getTweets()
      .subscribe((tweets: Array<Tweet>) => this.tweet_count = tweets.length)
    this.store.changes("my_tweets")
      .subscribe((tweets: Array<Tweet>) => this.tweet_count = tweets.length)

    this.userService.getFollowers()
      .subscribe((users: Array<User>) => this.follower_count = users.length)
    this.store.changes("my_followers")
      .subscribe((users: Array<User>) => this.follower_count = users.length)

    this.userService.getFollowing()
      .subscribe((users: Array<User>) => this.following_count = users.length)
    this.store.changes("my_following")
      .subscribe((users: Array<User>) => this.following_count = users.length)
  }
}
