import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Injectable } from '@angular/core';
import 'rxjs/Rx';

import { Alert } from '../core/alerts/shared';
import { User } from '../users/shared';
import { Tweet } from '../tweets/shared';


export interface State {
  user: User
  feed: Array<Tweet>
  my_tweets: Array<Tweet>
  my_following: Array<User>
  my_followers: Array<User>
  auth_token: string
  refresh_token: string
  alerts: Array<Alert>
}

const defaultState = {
  user: {},
  feed: [],
  my_tweets: [],
  my_following: [],
  my_followers: [],
  auth_token: "",
  refresh_token: "",
  alerts: [],
}

const _store = new BehaviorSubject<State>(defaultState);

@Injectable()
export class Store {
  private _store = _store
  private _changes = this._store.asObservable().distinctUntilChanged()

  setState(state: State) {
    this._store.next(state);
  }

  getState(): State {
    return this._store.value;
  }

  changes(name: string): any {
    return this._changes.pluck(name).distinctUntilChanged()
  }

  purge() {
    this._store.next(defaultState);
  }
}
