import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { User, AuthService } from '../shared';


@Component({
  template: `
    <div class="container">
      <h1>Login</h1>
      <form (ngSubmit)="onSubmit()" #loginForm="ngForm">
        <div class="form-group">
          <label for="email">Email</label>
          <input type="email" class="form-control" id="email"
                 required
                 [(ngModel)]="user.email" name="email"
                 #email="ngModel"
                 pattern="[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}" >
        </div>
        <div [hidden]="email.valid || email.pristine"
             class="alert alert-danger">
          Email is required
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input type="password" class="form-control" id="password"
                 required
                 [(ngModel)]="user.password" name="password"
                 #password="ngModel" >
        </div>
        <div [hidden]="password.valid || password.pristine"
             class="alert alert-danger">
          Password is required
        </div>

        <button type="submit" class="btn btn-default" [disabled]="!loginForm.form.valid">Submit</button>
      </form>
    </div>
  `
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
