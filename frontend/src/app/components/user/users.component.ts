import { Component, Input } from '@angular/core';

import { User } from '../../shared';


@Component({
  selector: 'users',
  template: `
    <div class="users shadow-2">
      <user
        *ngFor="let user of users"
        [user]="user"
        class="user"
      >
      </user>
    </div>
  `,
  styleUrls: ['./users.component.scss']
})
export class UsersComponent {
  @Input() users: User[]
}
