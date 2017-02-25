import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { User } from '../../users';
import { LoginService } from './login.service';


@Component({
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  user: User = {
    email: "",
    password: ""
  }

  constructor(
    private _loginService: LoginService,
    private _router: Router
  ) {}

  private loginWithGoogle() {
    this._loginService.authorizeWithGoogle() // firstly we have to authorize then login
      .subscribe(
        result => { window.location.href = result }, // redirect to Google AUTH
        error => {
          // TODO: bad message
          console.log("wtf")
        }
      );
  }

  private onSubmit() {
    this._loginService.login(this.user)
      .subscribe(
        result => {
          // TODO: good messgae
          this._router.navigateByUrl('home')
        },
        error => {
          // TODO: bad message
          console.log("wtf1")
        }
      );
  }
}
