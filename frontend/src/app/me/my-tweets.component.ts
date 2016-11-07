import { Component, OnInit } from '@angular/core';

import { UserService, Tweet } from '../shared';
import { Store } from '../store';


@Component({
  template: `
    <create-tweet></create-tweet>
    <tweets [tweets]="tweets"></tweets>
  `
})
export class MyTweetsComponent implements OnInit {
  tweets: Tweet[] = [
    {id: 1, author: {id: 2, name: "Name", username: "Username", email: "", password: "", created_at: ""}, likes: 1, retweets: 1, liked: false, retweeted: false, created_at: "", content: "Hello"}
  ]

  constructor(
    private _userService: UserService,
    private _store: Store
  ) {

  }

 ngOnInit(): void {
   this._userService.getFeed()

  //  this._store.changes.pluck('tweets')
  //    .subscribe((tweets: any) => this.tweets = tweets)
 }
}
