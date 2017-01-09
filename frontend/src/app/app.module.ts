import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';
import { HttpModule }     from '@angular/http';
import { RouterModule }   from '@angular/router';
import { APP_BASE_HREF }    from '@angular/common';

import { AppComponent }   from './app.component';

import {
  ApiService, AuthService, StoreHelper,
  TweetService, UserService, SearchService
} from './shared';
import { Store } from './store';

import { NavigationBarComponent } from './layout';
import {
  HomeComponent, SearchComponent,
  SignupComponent, LoginComponent, LoginGoogleCallbackComponent, LogoutComponent,
  TweetsModule, UsersModule,
} from './components';

import { MeModule } from './me'


@NgModule({
  imports: [
    HttpModule,
    BrowserModule,
    FormsModule,
    RouterModule.forRoot([
      { path: '', redirectTo: '/home', pathMatch: 'full' },
      { path: 'home', component: HomeComponent },
      { path: 'search', component: SearchComponent },

      { path: 'signup', component: SignupComponent },
      { path: 'login', component: LoginComponent },
      { path: 'login/google/callback', component: LoginGoogleCallbackComponent },
      { path: 'logout', component: LogoutComponent, canActivate: [AuthService] },
    ]),

    MeModule,
    TweetsModule,
    UsersModule,
  ],
  declarations: [
    AppComponent,
    NavigationBarComponent,
    HomeComponent,
    SearchComponent,

    SignupComponent,
    LoginComponent,
    LoginGoogleCallbackComponent,
    LogoutComponent,
  ],
  providers: [
    { provide: APP_BASE_HREF, useValue: '/' },
    ApiService,
    AuthService,
    StoreHelper,
    Store,
    TweetService,
    UserService,
    SearchService,
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
