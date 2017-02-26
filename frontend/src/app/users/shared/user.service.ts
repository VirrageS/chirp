import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import * as _ from 'lodash';

import { ApiService, Store, StoreHelper } from '../../shared';
import { User } from './user.model';


@Injectable()
export class UserService {
  private user_path: string = "/users";
  private user: User = null;

  constructor(
    private apiService: ApiService,
    private router: Router,
    private store: Store,
    private storeHelper: StoreHelper
  ) {
    this.store.changes('user')
      .subscribe((user: any) => this.user = user)
  }

  getUser() {
    return this.apiService.get(this.user_path + "/" + this.user.id)
      .do(user => this.storeHelper.update("user", user))
  }

  getTweets() {
    return this.apiService.get(this.user_path + "/" + this.user.id + "/tweets")
      .do(tweets => this.storeHelper.update("my_tweets", tweets))
  }

  getFeed() {
    return this.apiService.get("/feed")
      .do(tweets => this.storeHelper.update("feed", tweets))
  }

  getFollowing() {
    // NOTE: only here should happen name rewrite from followees => following
    return this.apiService.get(this.user_path + "/" + this.user.id + "/followees")
      .do(followees => this.storeHelper.update("my_following", followees))
  }

  getFollowers() {
    return this.apiService.get(this.user_path + "/" + this.user.id + "/followers")
      .do(followers => this.storeHelper.update("my_followers", followers))
  }

  follow(user_id: number) {
    return this.apiService.post(this.user_path + "/" + user_id + "/follow", {})
  }

  unfollow(user_id: number) {
    return this.apiService.post(this.user_path + "/" + user_id + "/unfollow", {})
  }
}
