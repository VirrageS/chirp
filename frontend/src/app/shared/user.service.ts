import { Injectable } from '@angular/core';

import { ApiService } from './api.service';
import { StoreHelper } from './store-helper';


@Injectable()
export class UserService {
  user_id: number = 0

  constructor(private _apiService: ApiService, private _storeHelper: StoreHelper) {}

  getUser() {
    return this._apiService.get("/user/" + this.user_id)
      .do(user => this._storeHelper.add('user', user))
  }

  getTweets(path) {
    return this._apiService.get("/user/" + this.user_id + path)
      .do(tweets => this._storeHelper.add('tweets', tweets))
  }
}
