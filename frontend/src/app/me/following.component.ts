import { Component, OnInit } from '@angular/core';

import { User, UserService } from '../users';
import { Store } from '../shared';


@Component({
  template: `
    <users [users]="following"></users>
  `
})
export class FollowingComponent implements OnInit {
  private following: Array<User> = [];

  constructor(
    private userService: UserService,
    private store: Store
  ) {}

  ngOnInit(): void {
    this.userService.getFollowing()
      .subscribe((users: Array<User>) => this.following = users)

    this.store.changes("my_following")
      .subscribe((users: Array<User>) => this.following = users)
  }
}
