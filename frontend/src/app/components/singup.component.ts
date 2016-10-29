import { Component } from '@angular/core';

import { User, UserService } from '../shared';


@Component({
  template: `
    {{diagnostic}}
    <div class="container">
      <h1>Singup</h1>
      <form *ngIf="active" (ngSubmit)="onSubmit()" #singupForm="ngForm">
        <div class="form-group">
          <label for="username">Username</label>
          <input type="text" class="form-control" id="username" [(ngModel)]="user.username" optional>
        </div>

        <div class="form-group">
          <label for="email">Email</label>
          <input type="email" class="form-control" id="email"
                 required
                 [(ngModel)]="user.email" name="email"
                 #email="ngModel" >
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

  constructor(private _userService: UserService) {}

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
