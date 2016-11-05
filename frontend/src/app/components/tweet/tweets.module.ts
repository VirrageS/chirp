import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';

import { TweetComponent }  from './tweet.component';
import { TweetsComponent } from './tweets.component';


@NgModule({
  imports: [
    CommonModule,
  ],
  declarations: [
    TweetComponent,
    TweetsComponent,
  ],
  exports: [
    TweetsComponent,
  ],
  providers: []
})
export class TweetsModule {}
