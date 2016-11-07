import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';
import { HttpModule }     from '@angular/http';
import { RouterModule }   from '@angular/router';
import { APP_BASE_HREF }    from '@angular/common';

import { AppComponent }   from './app.component';

import {
  ApiService, AuthService, StoreHelper,
  TweetService, UserService
} from './shared';
import { Store } from './store';

import { NavigationBarComponent } from './layout';
import {
  HomeComponent, SearchComponent,
  SignupComponent, LoginComponent, LogoutComponent,
  TweetsModule
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
      { path: 'logout', component: LogoutComponent, canActivate: [AuthService] },
    ]),

    MeModule,
    TweetsModule,
  ],
  declarations: [
    AppComponent,
    NavigationBarComponent,
    HomeComponent,
    SearchComponent,

    SignupComponent,
    LoginComponent,
    LogoutComponent,
  ],
  providers: [
    { provide: APP_BASE_HREF, useValue: '/' },
    ApiService,
    UserService,
    TweetService,
    AuthService,
    StoreHelper,
    Store,
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
