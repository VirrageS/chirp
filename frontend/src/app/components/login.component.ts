import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { User, UserService } from '../shared';


@Component({
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  user: User
  errors: string[]

  constructor(private _userService: UserService, private _router: Router) {
    this.user = {
      email: "",
      password: ""
    }
  }

  private loginWithGoogle() {
    this._userService.authorizeWithGoogle() // firstly we have to authorize then login
      .subscribe(
        result => { window.location.href = result }, // redirect to Google AUTH
        error => {} // TODO
      )
  }

  private onSubmit() {
    this._userService.login(this.user)
      .subscribe(
        result => {
          this._router.navigateByUrl('home')
        },
        error => { // TODO
          this.errors = error["errors"]
        }
      );
  }
}
