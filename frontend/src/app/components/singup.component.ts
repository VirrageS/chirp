import { Component } from '@angular/core';

import { User, UserService } from '../shared';


@Component({
  template: `
    {{diagnostic}}
    <div class="container">
      <h1>Singup</h1>
      <form (ngSubmit)="onSubmit()" #singupForm="ngForm">
        <div class="form-group">
          <label for="name">Name</label>
          <input type="name" class="form-control" id="name"
                 optional
                 [(ngModel)]="user.name" name="name"
                 #name="ngModel"
                 pattern="[a-zA-Z ]*" >
        </div>
        <div [hidden]="name.valid || name.pristine"
             class="alert alert-danger">
          Name is required
        </div>

        <div class="form-group">
          <label for="username">Username</label>
          <input type="username" class="form-control" id="username"
                 required
                 [(ngModel)]="user.username" name="username"
                 #username="ngModel"
                 pattern="[a-zA-Z0-9]*" >
        </div>
        <div [hidden]="username.valid || username.pristine"
             class="alert alert-danger">
          Username is required
        </div>

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

        <button type="submit" class="btn btn-default" [disabled]="!singupForm.form.valid">Submit</button>
      </form>
    </div>
  `
})
export class SingupComponent {
  user: User;

  constructor(private _userService: UserService) {
    this.user = new User();
  }

  onSubmit() {
    this._userService.singupUser(this.user)
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
