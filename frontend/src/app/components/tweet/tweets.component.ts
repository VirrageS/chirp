import { Component, Input } from '@angular/core';

import { Tweet } from '../../shared';


@Component({
  selector: 'tweets',
  styles: [`.tweets { max-width: 800px; margin: 0 auto; background-color: white; } .tweets .tweet { margin: 10px 0px; }`],
  template: `
    <div class="tweets">
      <tweet
        *ngFor="let tweet of tweets"
        [tweet]="tweet"
        class="tweet"
      >
      </tweet>
    </div>
  `
})
export class TweetsComponent {
  @Input() tweets: Tweet[]
}
