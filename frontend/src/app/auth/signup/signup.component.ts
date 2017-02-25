import { Component } from '@angular/core';
import { Router }    from '@angular/router';

import { User } from '../../users';
import { SignupService } from './signup.service'


@Component({
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})
export class SignupComponent {
  user: User = {
    name: "",
    username: "",
    email: "",
    password: ""
  }

  constructor(
    private _signupService: SignupService,
    private _router: Router
  ) {}

  private onSubmit() {
    this._signupService.signup(this.user)
      .subscribe(
        result => {
          // TODO: set message that everything is okay
          this._router.navigate(['', 'login'])
        },
        error => {
          // TODO: should show this message
          console.log(error["errors"])
        }
      )
  }
}
