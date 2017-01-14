import { Component, OnInit } from '@angular/core';

import { User, UserService } from '../shared';
import { Store } from '../store';


@Component({
  template: `
    <users [users]="followers"></users>
  `
})
export class FollowersComponent implements OnInit {
  followers: User[] = []

  constructor(
    private _userService: UserService,
    private _store: Store
  ) {

  }

  ngOnInit(): void {
    this._userService.getFollowers()
      .subscribe((users: any) => this.followers = users)

    this._store.changes("my_followers")
      .subscribe((users: any) => this.followers = users)
  }
}
