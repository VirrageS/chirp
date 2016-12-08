import { Component } from '@angular/core';
import { Router } from '@angular/router';

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

  onSubmit() {
    this._userService.login(this.user)
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
