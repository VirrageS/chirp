import { Injectable } from '@angular/core';

import { ApiService } from '../../shared';
import { AuthService } from '../auth.service';


@Injectable()
export class LoginService {
  constructor(
    private apiService: ApiService,
    private authService: AuthService,
  ) {}

  authorizeWithGoogle() {
    return this.apiService.get("/authorize/google")
  }

  loginWithGoogle(code, state) {
    return this.apiService.post("/login/google", {code: code, state: state})
      .do((res: any) => this.authService.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }

  login(body) {
    return this.apiService.post("/login", body)
      .do((res: any) => this.authService.setAuthorization(res.user, res.auth_token, res.refresh_token))
  }
}
