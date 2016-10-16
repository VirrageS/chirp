import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { FormsModule }    from '@angular/forms';
import { HttpModule }     from '@angular/http';
import { RouterModule }   from '@angular/router';

import { AppComponent }   from './app.component';

import { ApiService } from './shared';
import { NavigationBarComponent } from './navigation';

@NgModule({
  imports: [
    HttpModule,
    BrowserModule,
    FormsModule,
    RouterModule.forRoot([
      // { path: 'me', component: MeComponent },
      // { path: 'search', component: SearchComponent },
      // { path: 'login', component: LoginComponent },
      // { path: 'logout', component: LogoutComponent }
    ])
  ],
  declarations: [
    AppComponent,
    NavigationBarComponent
  ],
  providers: [
    ApiService
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
