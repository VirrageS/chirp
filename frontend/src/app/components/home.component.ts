import { Component } from '@angular/core';

import { User, Tweet, UserService } from '../shared';
import { Store } from '../store';


@Component({
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent {
  private feed: Array<Tweet> = []
  private loggedUser?: User

  constructor(
    private _userService: UserService,
    private _store: Store
  ) {
    this._store.changes("user")
      .subscribe((user: any) => {
        this.loggedUser = user

        if (this.loggedUser) {
          this._userService.getFeed()
            .subscribe((tweets: any) => this.feed = tweets)
          this._store.changes("feed")
            .subscribe((tweets: any) => this.feed = tweets)
        }
      })
  }
}
