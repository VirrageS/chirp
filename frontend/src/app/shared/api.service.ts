import { Http, Headers, Response, RequestOptions, Request, RequestMethod } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import 'rxjs/Rx';
import 'rxjs/add/observable/throw';
import { environment } from '../../environment';

import { AuthService } from '../auth/auth.service';
import { Store } from './store';
import { StoreHelper } from './store-helper';
import { User } from '../users';


@Injectable()
export class ApiService {
  retry: number = 2000;
  timeout: number = 5000;
  headers: Headers = new Headers({
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  });
  apiUrl: string = environment.api_url;
  refreshToken: string;
  user?: User;

  constructor(
    private _authService: AuthService,
    private _http: Http,
    private _store: Store,
    private _storeHelper: StoreHelper,
  ) {
    this._store.changes('user')
      .subscribe((user: any) => this.user = user)

    this._store.changes('auth_token')
      .subscribe((token: any) => this._setHeaders({
        Authorization: `Bearer ${token}`
      }))

    this._store.changes('refresh_token')
      .subscribe((token: any) => this.refreshToken = token)
  }

  get(path: string): Observable<any> {
    return this.request(RequestMethod.Get, path, {})
  }

  post(path: string, body): Observable<any> {
    return this.request(RequestMethod.Post, path, body)
  }

  delete(path: string): Observable<any> {
    return this.request(RequestMethod.Delete, path, {})
  }

  request(method: RequestMethod|string, path: string, body): Observable<any> {
    let options = new RequestOptions({
      method: method,
      headers: this.headers,
      url: `${this.apiUrl}${path}`,
      body: JSON.stringify(body)
    })
    return this._http.request(new Request(options))
      .catch((error: Response) => {
        if (error && error.status == 401) {
          // not authenticated so try to refresh token and send request again
          return this._refreshToken().flatMap((response: Response) => {
            options.headers = this.headers // update header
            return this._http.request(new Request(options))
          })
        }

        return Observable.throw(error)
      })
      .retryWhen(error => error.delay(this.retry))
      .timeout(this.timeout)
      .map(this._checkForError)
      .catch(err => Observable.throw(err))
      .map(this._getJson)
  }

  private _refreshToken() {
    return this._http.post(
        `${this.apiUrl}/token`,
        JSON.stringify({
          'user_id': this.user.id,
          'refresh_token': this.refreshToken
        }),
        { headers: this.headers }
      )
      .retryWhen(error => error.delay(this.retry))
      .timeout(this.timeout)
      .map(this._checkForError)
      .catch(err => Observable.throw(err))
      .map(this._getJson)
      .do(res => this._storeHelper.update('auth_token', res.auth_token))
  }

  private _setHeaders(headers) {
    Object.keys(headers).forEach(header => this.headers.set(header, headers[header]))
  }

  private _getJson(response: Response) {
    return response.json();
  }

  private _checkForError(response: Response): Response {
    if (response.status >= 200 && response.status < 300) {
      return response;
    } else if (response.status == 401) {
      this._authService.removeAuthorization();
      return response;
    } else {
      var error = new Error(response.statusText);
      error['response'] = response;
      throw error;
    }
  }
}
