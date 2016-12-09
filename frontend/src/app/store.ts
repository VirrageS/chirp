import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Injectable } from '@angular/core';
import 'rxjs/Rx';

import { User, Tweet } from './shared';


export interface State {
  feed: Array<Tweet>
  my_tweets: Array<Tweet>
  user: User
  following: Array<User>
  followers: Array<User>
  auth_token: string
  refresh_token: string
}

const defaultState = {
  feed: [],
  my_tweets: [],
  user: {},
  following: [],
  followers: [],
  auth_token: '',
  refresh_token: '',
}

const _store = new BehaviorSubject<State>(defaultState);

@Injectable()
export class Store {
  private _store = _store;
  changes = this._store.asObservable().distinctUntilChanged()

  setState(state: State) {
    this._store.next(state);
  }

  getState(): State {
    return this._store.value;
  }

  purge() {
    this._store.next(defaultState);
  }
}
