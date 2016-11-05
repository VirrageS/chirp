import { Component } from '@angular/core';

import { Tweet } from '../shared';


@Component({
  selector: 'me',
  templateUrl: './me.component.html',
  styleUrls: ['./me.component.scss']
})
export class MeComponent {
  followers: number = 2934890
  following: number = 384
  tweets: number = 234
}
