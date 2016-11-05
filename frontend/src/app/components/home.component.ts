import { Component, OnInit } from '@angular/core';

import { Tweet, UserService } from '../shared';


@Component({
  template: `
    <h2>HomeComponent</h2>
    <tweets [tweets]="tweets"></tweets>
  `
})
export class HomeComponent implements OnInit {
  tweets: Tweet[] = [
   {id: 1, author: {id: 2, name: "Name", username: "Username", email: "", password: "", created_at: ""}, likes: 1, retweets: 1, liked: false, retweeted: false, created_at: "", content: "Hello"}
 ]

 constructor(private _userService: UserService) {}

 ngOnInit(): void {
   // TODO: fetch user feed
 }
}
