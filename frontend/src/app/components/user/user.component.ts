import { Component, Input, Output, EventEmitter } from '@angular/core';

import { User, UserService } from '../../shared';
import { Store } from '../../store';
import * as _ from 'lodash';


@Component({
  selector: 'user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent {
  @Input() user: User
  @Output() userChange = new EventEmitter()
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

    this.userChange.emit(this.user)
    toggleFunc
      .subscribe(user => {
        // TODO: make this "this.user = user"
        // it involves changing whole binding system since `users.component`
        // will not update reference in table which will result is detached objects
        _.assign(this.user, user)
        this.userChange.emit(this.user)
        this._userService.getFollowers()
      })
  }
}
