import { Component, OnInit } from '@angular/core';
import { Subject } from 'rxjs/Subject';

import { User } from '../users';
import { Tweet } from '../tweets';
import { SearchService } from './search.service';
import * as _ from 'lodash';


@Component({
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit {
  private _searchTerms = new Subject<string>()
  users: User[]
  tweets: Tweet[]

  constructor(private _searchService: SearchService) {
    this.users = []
    this.tweets = []
  }

  ngOnInit(): void {
    this._searchTerms
      .debounceTime(300)
      .distinctUntilChanged()
      .filter(term => { return term != "" })
      .subscribe(term =>
        this._searchService.search(term)
          .subscribe(
            result => {
              this.users = result.users
              this.tweets = result.tweets
            },
            error => {} // TODO
          )
      )
  }

  private search(term: string): void {
    this._searchTerms.next(term)
  }

  private handleUserUpdated(updatedUser: User): void {
    let userTweets = _.filter(this.tweets, (tweet) => { return tweet.author.id == updatedUser.id })
    _.map(userTweets, (tweet) => _.assign(tweet.author, updatedUser))

    // NOTE: since there should be only one unique user
    // we should not do anything with `this.users`
  }

  private handleTweetUpdated(updatedTweet: Tweet) {
    let authorTweets = _.filter(this.tweets, (tweet) => { return tweet.author.id == updatedTweet.author.id })
    _.map(authorTweets, (tweet) => _.assign(tweet.author, updatedTweet.author))

    // NOTE: there should be only one user but this is more elegant
    let users = _.filter(this.users, (user) => { return user.id == updatedTweet.author.id })
    _.map(users, (user) => _.assign(user, updatedTweet.author))
  }
}
