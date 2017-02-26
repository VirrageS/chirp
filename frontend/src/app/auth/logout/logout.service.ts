import { Injectable } from '@angular/core';
import { Router } from '@angular/router';

import { AuthService } from '../auth.service';


@Injectable()
export class LogoutService {
  constructor(
    private authService: AuthService,
    private router: Router,
  ) {}

  logout() {
    this.authService.removeAuthorization()
    this.router.navigate(["", "home"]);
  }
}
