import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

import { ApiService, StoreHelper } from '../shared';
import { User } from '../users';


@Injectable()
export class AuthService implements CanActivate {
  AUTH_TOKEN_KEY: string = "AUTH_TOKEN"
  REFRESH_TOKEN_KEY: string = "REFRESH_TOKEN"
  USER_KEY: string = "USER_TOKEN"

  constructor(
      private _router: Router,
      private _apiService: ApiService,
      private _storeHelper: StoreHelper
  ) {
    this.initializeAuthorization()
  }

  setAuthorization(user: User, authToken: string, refreshToken: string) {
    window.localStorage.setItem(this.USER_KEY, JSON.stringify(user))
    this._storeHelper.update('user', user)

    window.localStorage.setItem(this.AUTH_TOKEN_KEY, authToken)
    this._storeHelper.update('auth_token', authToken)

    window.localStorage.setItem(this.REFRESH_TOKEN_KEY, refreshToken)
    this._storeHelper.update('refresh_token', refreshToken)
  }

  initializeAuthorization() {
    const user = window.localStorage.getItem(this.USER_KEY)
    const authToken = window.localStorage.getItem(this.AUTH_TOKEN_KEY)
    const refreshToken = window.localStorage.getItem(this.REFRESH_TOKEN_KEY)

    if (user && authToken && refreshToken) {
      this.setAuthorization(JSON.parse(user), authToken, refreshToken)
    } else {
      this.removeAuthorization()
    }
  }

  removeAuthorization() {
    window.localStorage.removeItem(this.USER_KEY)
    this._storeHelper.update('user', null)

    window.localStorage.removeItem(this.AUTH_TOKEN_KEY)
    this._storeHelper.update('auth_token', '')

    window.localStorage.removeItem(this.REFRESH_TOKEN_KEY)
    this._storeHelper.update('refresh_token', '')

    this._router.navigate(['', 'home']);
  }

  isAuthenticated(): boolean {
    const user = window.localStorage.getItem(this.USER_KEY)
    const authToken = window.localStorage.getItem(this.AUTH_TOKEN_KEY)
    const refreshToken = window.localStorage.getItem(this.REFRESH_TOKEN_KEY)
    return (!!user && !!authToken && !!refreshToken)
  }

  canActivate(): boolean {
    const canActivate = this.isAuthenticated()
    this.onCanActivate(canActivate)
    return canActivate
  }

  canActivateChild(): boolean {
    return this.canActivate()
  }

  onCanActivate(canActivate: boolean) {
    if (!canActivate) {
      this._router.navigate(['', 'login']);
    }
  }

  authorizeWithGoogle() {
    return this._apiService.get("/authorize/google")
  }

  loginWithGoogle(code, state) {
    return this._apiService.post("/login/google", {code: code, state: state})
      .do((res: any) => this.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }

  signup(body) {
    return this._apiService.post("/signup", body)
  }

  login(body) {
    return this._apiService.post("/login", body)
      .do((res: any) => this.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }

  logout() {
    this.removeAuthorization()
    this._router.navigate(['', 'home']);
  }
}
