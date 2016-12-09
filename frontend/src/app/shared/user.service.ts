import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import * as _ from 'lodash';

import { ApiService } from './api.service';
import { AuthService } from './auth.service';
import { Store } from '../store';
import { StoreHelper } from './store-helper';
import { User } from './user.model';


@Injectable()
export class UserService {
  user?: User

  constructor(
    private _apiService: ApiService,
    private _authService: AuthService,
    private _router: Router,
    private _store: Store,
    private _storeHelper: StoreHelper
  ) {
    this._store.changes.pluck('user')
      .subscribe((user: any) => this.user = user)
  }

  getUser() {
    return this._apiService.get("/users/" + this.user.id)
      .do(user => this._storeHelper.update("user", user))
  }

  getTweets() {
    // TODO: update path
    // return this._apiService.get("/users/" + this.user.id + path)
    return this._apiService.get("/home_feed")
      .do(tweets => this._storeHelper.update("my_tweets", tweets))
  }

  getFeed() {
    return this._apiService.get("/home_feed")
      .do(tweets => this._storeHelper.update("feed", tweets))
  }

  signup(body) {
    return this._apiService.post("/signup", body)
  }

  login(body) {
    return this._apiService.post("/login", body)
      .do((res: any) => this._authService.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }

  logout() {
    this._authService.removeAuthorization()
    this._router.navigate(['', 'home']);
  }
}
