import { Component } from '@angular/core';
import { Router }    from '@angular/router';

import { User } from '../users';
import { AuthService } from './auth.service'


@Component({
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent {
  user: User

  constructor(private _authService: AuthService, private _router: Router) {
    this.user = {
      name: "",
      username: "",
      email: "",
      password: ""
    }
  }

  private onSubmit() {
    this._authService.signup(this.user)
      .subscribe(
        result => {
          // TODO set message that everything is okay
          this._router.navigate(['', 'login'])
        },
        error => { 
          // TODO: should show this message
          console.log(error["errors"])
        }
      )
  }
}
