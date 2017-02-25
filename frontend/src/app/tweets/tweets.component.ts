import { Component, Input, Output, EventEmitter } from '@angular/core';

import { Tweet } from './shared';


@Component({
  selector: 'tweets',
  template: `
    <div class="tweets shadow-2">
      <tweet
        *ngFor="let tweet of tweets"
        [tweet]="tweet"
        class="tweet"
        (tweetChange)="handleTweetUpdated($event)"
      >
      </tweet>
    </div>
  `,
  styleUrls: ['./tweets.component.scss']
})
export class TweetsComponent {
  @Input() tweets: Tweet[];
  @Output() tweetUpdated = new EventEmitter();

  // propagate change
  private handleTweetUpdated(tweet: Tweet) {
    this.tweetUpdated.emit(tweet)
  }
}
