import { Component, Input } from '@angular/core';

import { User, UserService } from '../../shared';

import '../../../../public/scss/abstract.scss';


@Component({
  selector: 'user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.scss']
})
export class UserComponent {
  @Input() user: User

  constructor(
    private _userService: UserService
  ) {

  }
}
