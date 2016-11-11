import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { FormsModule }    from '@angular/forms';

import { TweetComponent }       from './tweet.component';
import { TweetsComponent }      from './tweets.component';
import { CreateTweetComponent } from './create-tweet.component';


@NgModule({
  imports: [
    CommonModule,
    FormsModule,
  ],
  declarations: [
    TweetComponent,
    TweetsComponent,
    CreateTweetComponent,
  ],
  exports: [
    TweetsComponent,
    CreateTweetComponent,
  ],
  providers: []
})
export class TweetsModule {}
