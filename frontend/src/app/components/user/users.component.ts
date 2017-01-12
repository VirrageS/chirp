import { Component, Input, Output, EventEmitter } from '@angular/core';

import { User } from '../../shared';
import * as _ from 'lodash';


@Component({
  selector: 'users',
  template: `
    <div class="users shadow-2">
      <user
        *ngFor="let user of users"
        [user]="user"
        class="user"
        (userChange)="handleUserUpdated($event)"
      >
      </user>
    </div>
  `,
  styleUrls: ['./users.component.scss']
})
export class UsersComponent {
  @Input() users: User[]
  @Output() userUpdated = new EventEmitter()

  // propagate change
  private handleUserUpdated(user: User) {
    this.userUpdated.emit(user)
  }
}
