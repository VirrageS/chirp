import { Component } from '@angular/core';

import { AuthService } from '../shared';


@Component({
  template: `
    <h2>LogoutComponent</h2>
    <p>Get your heroes here</p>
  `
})
export class LogoutComponent {
  constructor(private _authService: AuthService) {
    this._authService.logout()
  }
}
