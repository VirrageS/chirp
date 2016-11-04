import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';
import { HttpModule }     from '@angular/http';
import { RouterModule }   from '@angular/router';

import { AppComponent }   from './app.component';

import { ApiService, UserService } from './shared';
import { NavigationBarComponent } from './layout';

import {
  HomeComponent, SearchComponent,
  SingupComponent, LoginComponent, LogoutComponent
} from './components';

import { MeModule } from './me'


@NgModule({
  imports: [
    HttpModule,
    BrowserModule,
    FormsModule,
    MeModule,
    RouterModule.forRoot([
      { path: '', redirectTo: '/home', pathMatch: 'full' },
      { path: 'home', component: HomeComponent },
      { path: 'search', component: SearchComponent },

      { path: 'singup', component: SingupComponent },
      { path: 'login', component: LoginComponent },
      { path: 'logout', component: LogoutComponent },
    ])
  ],
  declarations: [
    AppComponent,
    NavigationBarComponent,
    HomeComponent,
    SearchComponent,

    SingupComponent,
    LoginComponent,
    LogoutComponent,
  ],
  providers: [
    ApiService,
    UserService,
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
