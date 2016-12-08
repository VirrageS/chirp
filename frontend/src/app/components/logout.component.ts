import { Component } from '@angular/core';

import { UserService } from '../shared';


@Component({
  template: `
    <h2>LogoutComponent</h2>
    <p>Get your heroes here</p>
  `
})
export class LogoutComponent {
  constructor(private _userService: UserService) {
    this._userService.logout()
  }
}
