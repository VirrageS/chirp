import { NgModule, forwardRef }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';
import { HttpModule }     from '@angular/http';
import { RouterModule }   from '@angular/router';
import { APP_BASE_HREF }    from '@angular/common';

import { AppComponent }   from './app.component';

import { ApiService, StoreHelper, Store } from './shared';
import { AuthService, AuthModule } from './auth';
import { HomeComponent } from './home';
import { SearchModule } from './search';
import { AlertsModule, NavComponent, PageNotFoundComponent } from './core';

import { MeModule } from './me';
import { TweetsModule } from './tweets';
import { UsersModule } from './users';


@NgModule({
  imports: [
    HttpModule,
    BrowserModule,
    FormsModule,
    RouterModule.forRoot([
      { path: '', redirectTo: '/home', pathMatch: 'full' },
      { path: 'home', component: HomeComponent },

      { path: '**', component: PageNotFoundComponent },
    ]),

    AlertsModule,
    AuthModule,
    MeModule,
    TweetsModule,
    UsersModule,
    SearchModule,
  ],
  declarations: [
    AppComponent,
    NavComponent,
    HomeComponent,

    PageNotFoundComponent,
  ],
  providers: [
    { provide: APP_BASE_HREF, useValue: '/' },
    AuthService,
    ApiService,
    StoreHelper,
    Store,
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
