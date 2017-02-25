import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { FormsModule }    from '@angular/forms';

import { TweetService } from './shared';
import { TweetComponent }       from './tweet/tweet.component';
import { TweetsComponent }      from './tweets.component';
import { CreateTweetComponent } from './create-tweet/create-tweet.component';


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
  providers: [
    TweetService,
  ]
})
export class TweetsModule {}
