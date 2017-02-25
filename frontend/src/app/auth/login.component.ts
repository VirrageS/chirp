import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { User } from '../users';
import { AuthService } from './auth.service';


@Component({
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent {
  user: User

  constructor(private _authService: AuthService, private _router: Router) {
    this.user = {
      email: "",
      password: ""
    }
  }

  private loginWithGoogle() {
    this._authService.authorizeWithGoogle() // firstly we have to authorize then login
      .subscribe(
        result => { window.location.href = result }, // redirect to Google AUTH
        error => {} // TODO
      )
  }

  private onSubmit() {
    this._authService.login(this.user)
      .subscribe(
        result => {
          // TODO
          this._router.navigateByUrl('home')
        },
        error => { // TODO
          console.log(error["errors"])
        }
      )
  }
}
