import { Component, Input } from '@angular/core';

import { User, UserService } from '../../shared';
import { Store } from '../../store';


@Component({
  selector: 'user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent {
  @Input() user: User
  loggedUser: User

  constructor(
    private _userService: UserService,
    private _store: Store
  ) {
    this._store.changes.pluck("user")
      .subscribe((user: any) => this.loggedUser = user)
  }

  private _toggleFollow() {
    this.user.following = !this.user.following

    // send real request
    let toggleFunc = this._userService.follow(this.user.id)
    if (!this.user.following) {
      toggleFunc = this._userService.unfollow(this.user.id)
    }

    toggleFunc
      .subscribe(user => this.user = user)

    this._userService.getFollowers()
  }
}
