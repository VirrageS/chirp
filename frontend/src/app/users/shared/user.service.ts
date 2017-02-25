import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import * as _ from 'lodash';

import { ApiService, Store, StoreHelper } from '../../shared';
import { User } from './user.model';


@Injectable()
export class UserService {
  private user_path: string = "/users";
  user?: User;

  constructor(
    private _apiService: ApiService,
    private _router: Router,
    private _store: Store,
    private _storeHelper: StoreHelper
  ) {
    this._store.changes('user')
      .subscribe((user: any) => this.user = user)
  }

  getUser() {
    return this._apiService.get(this.user_path + "/" + this.user.id)
      .do(user => this._storeHelper.update("user", user))
  }

  getTweets() {
    return this._apiService.get(this.user_path + "/" + this.user.id + "/tweets")
      .do(tweets => this._storeHelper.update("my_tweets", tweets))
  }

  getFeed() {
    return this._apiService.get("/feed")
      .do(tweets => this._storeHelper.update("feed", tweets))
  }

  getFollowing() {
    // NOTE: only here should happen name rewrite from followees => following
    return this._apiService.get(this.user_path + "/" + this.user.id + "/followees")
      .do(followees => this._storeHelper.update("my_following", followees))
  }

  getFollowers() {
    return this._apiService.get(this.user_path + "/" + this.user.id + "/followers")
      .do(followers => this._storeHelper.update("my_followers", followers))
  }

  follow(user_id: number) {
    return this._apiService.post(this.user_path + "/" + user_id + "/follow", {})
  }

  unfollow(user_id: number) {
    return this._apiService.post(this.user_path + "/" + user_id + "/unfollow", {})
  }
}
