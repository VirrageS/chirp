import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';
import { HttpModule }     from '@angular/http';
import { RouterModule }   from '@angular/router';

import { AppComponent }   from './app.component';

import { ApiService } from './shared';
import { NavigationBarComponent } from './layout';

import {
  HomeComponent, SearchComponent,
  LoginComponent, LogoutComponent
} from './components';

import { MeComponent, TweetsComponent, FollowingComponent, FollowersComponent } from './me';

@NgModule({
  imports: [
    HttpModule,
    BrowserModule,
    FormsModule,
    RouterModule.forRoot([
      { path: '', redirectTo: '/home', pathMatch: 'full' },
      { path: 'home', component: HomeComponent },
      { path: 'me', component: MeComponent,
        children: [
          { path: '', redirectTo: 'tweets' },
          { path: 'tweets', component: TweetsComponent },
          { path: 'following', component: FollowingComponent },
          { path: 'followers', component: FollowersComponent },
        ]
      },
      { path: 'search', component: SearchComponent },

      { path: 'login', component: LoginComponent },
      { path: 'logout', component: LogoutComponent },
    ])
  ],
  declarations: [
    AppComponent,
    NavigationBarComponent,
    HomeComponent,
    SearchComponent,
    LoginComponent,
    LogoutComponent,

    MeComponent,
    TweetsComponent,
    FollowingComponent,
    FollowersComponent,
  ],
  providers: [
    ApiService
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
