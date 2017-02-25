import { Component, Input, Output, EventEmitter } from '@angular/core';

import { User } from './shared';

@Component({
  selector: 'users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss']
})
export class UsersComponent {
  @Input() users: User[];
  @Output() userUpdated = new EventEmitter();

  // propagate change
  private handleUserUpdated(user: User) {
    this.userUpdated.emit(user);
  }
}
