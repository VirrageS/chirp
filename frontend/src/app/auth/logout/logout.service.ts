import { Injectable } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService } from '../auth.service';

@Injectable()
export class LogoutService {
  constructor(
    private _authService: AuthService,
    private _router: Router,
  ) {}

  logout() {
    this._authService.removeAuthorization()
    this._router.navigate(['', 'home']);
  }
}
