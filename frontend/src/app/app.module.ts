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
import { NavComponent, PageNotFoundComponent } from './core';

import { TweetsModule} from './tweets';
import { UsersModule} from './users';
import { MeModule } from './me'


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
