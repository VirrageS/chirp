import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

import { ApiService } from './api.service';
import { StoreHelper } from './store-helper';
import { User } from './user.model';


@Injectable()
export class AuthService implements CanActivate {
  TOKEN_KEY: string = "AUTH_TOKEN"
  USER_KEY: string = "USER_TOKEN"

  constructor(
      private _apiService: ApiService,
      private _router: Router,
      private _storeHelper: StoreHelper
  ) {
    this.refreshAuthorization()
  }

  setAuthorization(token: string, user: User) {
    window.localStorage.setItem(this.TOKEN_KEY, token)
    this._apiService.setHeaders({
      Authorization: `Bearer ${token}`
    })

    window.localStorage.setItem(this.USER_KEY, JSON.stringify(user))
    this._storeHelper.update('user', user)
  }

  refreshAuthorization() {
    const token = window.localStorage.getItem(this.TOKEN_KEY)
    const user = window.localStorage.getItem(this.USER_KEY)

    if (token && user) {
      this.setAuthorization(token, JSON.parse(user))
    } else {
      this.removeAuthorization()
    }
  }

  removeAuthorization() {
    window.localStorage.removeItem(this.TOKEN_KEY)
    this._apiService.setHeaders({
      Authorization: `Bearer `
    })

    window.localStorage.removeItem(this.USER_KEY)
    this._storeHelper.update('user', null)
  }

  isAuthenticated(): boolean {
    const token = window.localStorage.getItem(this.TOKEN_KEY)
    const user = window.localStorage.getItem(this.USER_KEY)
    return (!!token && !!user)
  }

  canActivate(): boolean {
    const canActivate = this.isAuthenticated()
    this.onCanActivate(canActivate)
    return canActivate
  }

  onCanActivate(canActivate: boolean) {
    if (!canActivate) {
      this._router.navigate(['', 'login']);
    }
  }

  signup(body) {
    return this._apiService.post("/signup", body);
  }

  login(body) {
    return this._apiService.post("/login", body)
      .do((res: any) => console.log(res.user))
      .do((res: any) => this.setAuthorization(res.auth_token, res.user))
  }

  logout() {
    this.removeAuthorization()
  }
}
