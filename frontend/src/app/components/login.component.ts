import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { User, AuthService } from '../shared';


@Component({
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  user: User
  errors: string[]

  constructor(private _authService: AuthService, private _router: Router) {
    this.user = {
      email: "",
      password: ""
    }
  }

  onSubmit() {
    this._authService.login(this.user)
      .subscribe(
        result => {
          this._router.navigateByUrl('home')
        },
        error => {
          this.errors = error["errors"]
        }
      );
  }
}
