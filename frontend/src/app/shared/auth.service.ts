import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

import { ApiService } from './api.service';


@Injectable()
export class AuthService implements CanActivate {
  TOKEN_KEY: string = "AUTH_TOKEN"

  constructor(private _apiService: ApiService, private _router: Router) {
    const token = window.localStorage.getItem(this.TOKEN_KEY)

    if (token) {
      this.setToken(token)
    }
  }

  setToken(token: string) {
    window.localStorage.setItem(this.TOKEN_KEY, token)
    this._apiService.setHeaders({
      Authorization: `Bearer ${token}`
    })
  }

  isAuthenticated(): boolean {
    const token = window.localStorage.getItem(this.TOKEN_KEY)
    return (token !== null)
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
      .do((res: any) => this.setToken(res.auth_token))
  }

  logout() {
    window.localStorage.removeItem(this.TOKEN_KEY)
  }
}
