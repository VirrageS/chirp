import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

import { ApiService, StoreHelper } from '../../shared';
import { AuthService } from '../auth.service';

@Injectable()
export class LoginService {
  constructor(
    private _router: Router,
    private _apiService: ApiService,
    private _authService: AuthService
  ) {}

  authorizeWithGoogle() {
    return this._apiService.get("/authorize/google")
  }

  loginWithGoogle(code, state) {
    return this._apiService.post("/login/google", {code: code, state: state})
      .do((res: any) => this._authService.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }

  login(body) {
    return this._apiService.post("/login", body)
      .do((res: any) => this._authService.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }
}
