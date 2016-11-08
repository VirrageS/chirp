import { Component } from '@angular/core';
import { Router }    from '@angular/router';

import { User, AuthService } from '../shared';


@Component({
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent {
  user: User
  errors: string[]

  constructor(private _authService: AuthService, private _router: Router) {
    this.user = {
      name: "",
      username: "",
      email: "",
      password: ""
    }
  }

  onSubmit() {
    this._authService.signup(this.user)
      .subscribe(
        result => {
          this._router.navigateByUrl('login')
        },
        error => {
          this.errors = error["errors"]
        }
      )
  }
}
