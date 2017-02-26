import { Component, OnInit } from '@angular/core';

import { User, UserService } from '../users';
import { Store } from '../shared';


@Component({
  template: `
    <users [users]="followers"></users>
  `
})
export class FollowersComponent implements OnInit {
  private followers: Array<User> = [];

  constructor(
    private userService: UserService,
    private store: Store
  ) {}

  ngOnInit(): void {
    this.userService.getFollowers()
      .subscribe((users: Array<User>) => this.followers = users);

    this.store.changes("my_followers")
      .subscribe((users: Array<User>) => this.followers = users);
  }
}
