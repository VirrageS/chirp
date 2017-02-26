import { Injectable } from '@angular/core';

import { ApiService } from '../../shared';


@Injectable()
export class SignupService {
  constructor(
    private apiService: ApiService
  ) {}

  signup(body) {
    return this.apiService.post("/signup", body)
  }
}
