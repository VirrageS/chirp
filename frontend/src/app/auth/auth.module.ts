import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule }   from '@angular/router';

import { AuthService } from './auth.service';

import { LoginComponent, LoginGoogleCallbackComponent, LoginService } from './login';
import { LogoutComponent, LogoutService } from './logout';
import { SignupComponent, SignupService } from './signup';


@NgModule({
  imports: [
    CommonModule,
    FormsModule,

    RouterModule.forChild([
      { path: 'signup', component: SignupComponent },
      { path: 'login', component: LoginComponent },
      { path: 'login/google/callback', component: LoginGoogleCallbackComponent },
      { path: 'logout', component: LogoutComponent, canActivate: [AuthService] },
    ])
  ],
  declarations: [
    LoginComponent,
    LoginGoogleCallbackComponent,
    LogoutComponent,
    SignupComponent,
  ],
  exports: [],
  providers: [
    AuthService,
    LoginService,
    LogoutService,
    SignupService,
  ]
})
export class AuthModule {}
