import { Injectable } from '@angular/core';
import { CanActivate, Router } from '@angular/router';

import { ApiService } from '../../shared';

@Injectable()
export class SignupService {
  constructor(
    private _apiService: ApiService
  ) {}

  signup(body) {
    return this._apiService.post("/signup", body)
  }
}
