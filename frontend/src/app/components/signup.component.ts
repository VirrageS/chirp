import { Component } from '@angular/core';
import { Router }    from '@angular/router';

import { User, UserService } from '../shared';


@Component({
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent {
  user: User
  errors: string[]

  constructor(private _userService: UserService, private _router: Router) {
    this.user = {
      name: "",
      username: "",
      email: "",
      password: ""
    }
  }

  private onSubmit() {
    this._userService.signup(this.user)
      .subscribe(
        result => this._router.navigate(['', 'login']),
        error => this.errors = error["errors"]
      )
  }
}
