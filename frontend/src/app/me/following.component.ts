import { Component, OnInit } from '@angular/core';

import { User, UserService } from '../users';
import { Store } from '../shared';


@Component({
  template: `
    <users [users]="following"></users>
  `
})
export class FollowingComponent implements OnInit {
  following: User[] = []

  constructor(
    private _userService: UserService,
    private _store: Store
  ) {

  }

  ngOnInit(): void {
    this._userService.getFollowing()
      .subscribe((users: any) => this.following = users)

    this._store.changes("my_following")
      .subscribe((users: any) => this.following = users)
  }
}
