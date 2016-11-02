import { Component } from '@angular/core';

@Component({
  styleUrls: ['./me.component.scss'],
  templateUrl: './me.component.html'
})
export class MeComponent {
  followers: number = 2934890;
  following: number = 384;
  tweets: number = 809;
}
