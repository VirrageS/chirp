import { Injectable } from '@angular/core';
import * as _ from 'lodash';

import { Store } from '../store';
import { ApiService } from './api.service';
import { StoreHelper } from './store-helper';
import { User } from './user.model';


@Injectable()
export class UserService {
  user?: User

  constructor(
    private _apiService: ApiService,
    private _store: Store,
    private _storeHelper: StoreHelper
  ) {
    this._store.changes.pluck('user')
      .subscribe((user: any) => this.user = user)
  }

  getUser() {
    return this._apiService.get("/users/" + this.user.id)
      .do(user => this._storeHelper.add("user", user))
  }

  getTweets(path) {
    return this._apiService.get("/users/" + this.user.id + path)
      .do(tweets => this._storeHelper.add("tweets", tweets))
  }

  getFeed() {
    return this._apiService.get("/home_feed")
      .do(tweets => {
        _.each(tweets, (tweet) => {
          this._storeHelper.add("tweets", tweet)
        })
      })
  }
}
