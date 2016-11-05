import { Component, Input } from '@angular/core';

import { Tweet } from '../shared';


@Component({
  template: `
    <tweets [tweets]="tweets"></tweets>
  `
})
export class MyTweetsComponent {
  tweets: Tweet[] = [
   {id: 1, author: {id: 2, name: "Name", username: "Username", email: "", password: "", created_at: ""}, likes: 1, retweets: 1, liked: false, retweeted: false, created_at: "", content: "Hello"}
 ]
}
