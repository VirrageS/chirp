import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

import { UserService } from '../shared';


@Component({
  template: `<h2>Logging in...</h2>`
})
export class LoginGoogleCallbackComponent implements OnInit {
  constructor(
    private _userService: UserService,
    private _router: Router,
    private _activedRoute: ActivatedRoute
  ) {
  }

  ngOnInit() {
    // subscribe to router event
    this._activedRoute.queryParams.subscribe(
      (param: any) => {
        let code = param['code'];
        let state = param['state'];
        if (code && state) {
          this._userService.loginWithGoogle(code, state)
            .subscribe(
              result => this._router.navigate(['', 'home']),
              error => this._router.navigate(['', 'login'])
            )
        } else {
          this._router.navigate(['', 'login'])
        }
      }
    )
  }
}
