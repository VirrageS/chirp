import { Http, Headers, Response } from '@angular/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import 'rxjs/Rx';
import 'rxjs/add/observable/throw';


@Injectable()
export class ApiService {
  retry: number = 2000;
  timeout: number = 5000;
  headers: Headers = new Headers({
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  });
  apiUrl: string;

  constructor(private _http: Http) {
    this.apiUrl = 'http://0.0.0.0:8080'
    if (process.env.ENV === 'production') {
      this.apiUrl = '' // TODO: we should set some apiUrl when on production
    }
  }

  get(path: string): Observable<any> {
    console.log(this.headers)
    return this._http.get(`${this.apiUrl}${path}`, { headers: this.headers, body: {} })
      .retryWhen(error => error.delay(this.retry))
      .timeout(this.timeout, new Error('Delay exceeded'))
      .map(this._checkForError)
      .catch(err => Observable.throw(err))
      .map(this._getJson)
  }

  post(path: string, body): Observable<any> {
    return this._http.post(
        `${this.apiUrl}${path}`,
        JSON.stringify(body),
        { headers: this.headers }
      )
      .map(this._checkForError)
      .catch(err => Observable.throw(err))
      .map(this._getJson)
  }

  delete(path: string): Observable<any> {
    return this._http.delete(`${this.apiUrl}${path}`, { headers: this.headers })
      .map(this._checkForError)
      .catch(err => Observable.throw(err))
      .map(this._getJson)
  }

  setHeaders(headers) {
    Object.keys(headers).forEach(header => this.headers.set(header, headers[header]))
  }

  private _getJson(response: Response) {
    return response.json();
  }

  private _checkForError(response: Response): Response {
    if (response.status >= 200 && response.status < 300) {
      return response;
    } else {
      var error = new Error(response.statusText);
      error['response'] = response;
      throw error;
    }
  }
}
