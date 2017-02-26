import { Component } from '@angular/core';
import { Router }    from '@angular/router';

import { AlertType } from '../../core/alerts';
import { User } from '../../users';
import { SignupService } from './signup.service'
import { StoreHelper } from '../../shared';


@Component({
  templateUrl: './signup.component.html',
  // styleUrls: ['./signup.component.scss'],
})
export class SignupComponent {
  private user: User = {
    name: "",
    username: "",
    email: "",
    password: ""
  }

  constructor(
    private router: Router,
    private signupService: SignupService,
    private storeHelper: StoreHelper,
  ) {}

  private onSubmit() {
    this.signupService.signup(this.user)
      .subscribe(
        result => {
          this.storeHelper.add("alerts", {message: "Account has been created successfully", type: AlertType.success});
          this.router.navigate(['', 'login'])
        },
        error => {
          this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
        }
      )
  }
}
