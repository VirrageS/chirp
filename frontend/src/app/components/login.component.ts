import { Component } from '@angular/core';

import { User, UserService } from '../shared';


@Component({
  template: `
    <div class="container">
      <h1>Login</h1>
      <form (ngSubmit)="onSubmit()" #loginForm="ngForm">
        <div class="form-group">
          <label for="username">Username</label>
          <input type="username" class="form-control" id="username"
                 required
                 [(ngModel)]="user.username" name="username"
                 #username="ngModel" >
        </div>
        <div [hidden]="username.valid || username.pristine"
             class="alert alert-danger">
          Username is required
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
  user: User;

  constructor(private _userService: UserService) {
    this.user = new User();
  }

  onSubmit() {
    this._userService.loginUser(this.user)
      .subscribe(
        result => {
          console.log(result);
        },
        error => {
          console.log(error);
        }
      );
  }
}
