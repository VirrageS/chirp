import { Component } from '@angular/core';

import { LogoutService } from './logout.service';


@Component({
  template: ``,
})
export class LogoutComponent {
  constructor(
    private _logoutService: LogoutService
  ) {
    this._logoutService.logout()
  }
}
