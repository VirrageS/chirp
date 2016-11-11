import { Component, Input } from '@angular/core';

import { Tweet } from '../../shared';


@Component({
  selector: 'tweets',
  template: `
    <div class="tweets shadow-2">
      <tweet
        *ngFor="let tweet of tweets"
        [tweet]="tweet"
        class="tweet"
      >
      </tweet>
    </div>
  `,
  styleUrls: ['./tweets.component.scss']
})
export class TweetsComponent {
  @Input() tweets: Tweet[]
}
