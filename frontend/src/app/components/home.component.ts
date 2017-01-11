import { Component, OnInit } from '@angular/core';

import { Tweet, UserService } from '../shared';


@Component({
  template: `
    <h2>Chirp - the real Twitter</h2>

  `
})
export class HomeComponent implements OnInit {
 constructor(private _userService: UserService) {}

 ngOnInit(): void {}
}
