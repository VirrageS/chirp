import { Component, OnInit } from '@angular/core';

import { User, UserService } from '../shared';
import { Store } from '../store';


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

    this._store.changes.pluck("my_following")
      .subscribe((users: any) => this.following = users)
  }
}
