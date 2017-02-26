import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { AlertType } from '../../core/alerts';
import { LoginService } from './login.service';
import { StoreHelper } from '../../shared';


@Component({
  template: `<h2>Logging in...</h2>`
})
export class LoginGoogleCallbackComponent implements OnInit {
  constructor(
    private activedRoute: ActivatedRoute,
    private loginService: LoginService,
    private router: Router,
    private storeHelper: StoreHelper,
  ) {}

  ngOnInit() {
    // subscribe to router event
    this.activedRoute.queryParams.subscribe(
      (param: any) => {
        let code = param["code"];
        let state = param["state"];
        if (code && state) {
          this.loginService.loginWithGoogle(code, state)
            .subscribe(
              result => this.router.navigate(["", "home"]),
              error => this.router.navigate(["", "login"])
            )
        } else {
          this.storeHelper.add("alerts", {message: "Failed to login with Google. Try one more time", type: AlertType.danger});
          this.router.navigate(["", "login"])
        }
      }
    )
  }
}
