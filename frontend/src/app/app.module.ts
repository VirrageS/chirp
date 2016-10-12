import { NgModule } from '@angular/core';
import { BrowserModule  } from '@angular/platform-browser';
import { AppComponent } from './app.component';
import { HttpModule } from '@angular/http';

import { ApiService } from './shared';
import { MainComponent } from './containers';

@NgModule({
  imports: [
    HttpModule,
    BrowserModule
  ],
  declarations: [
    AppComponent,
    MainComponent
  ],
  providers: [
    ApiService
  ],
  bootstrap: [AppComponent]
 })
 export class AppModule {}
