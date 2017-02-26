import { Component, OnInit } from '@angular/core';
import { Subject } from 'rxjs/Subject';

import { AlertType } from '../core/alerts';
import { User } from '../users';
import { Tweet } from '../tweets';
import { SearchService } from './search.service';
import { StoreHelper } from '../shared';

import * as _ from 'lodash';


@Component({
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent implements OnInit {
  private searchTerms = new Subject<string>()
  private users: User[] = [];
  private tweets: Tweet[] = [];

  constructor(
    private searchService: SearchService,
    private storeHelper: StoreHelper,
  ) {}

  ngOnInit(): void {
    this.searchTerms
      .debounceTime(300)
      .distinctUntilChanged()
      .filter(term => { return term != "" })
      .subscribe(term =>
        this.searchService.search(term)
          .subscribe(
            result => {
              this.users = result.users
              this.tweets = result.tweets
            },
            error => {
              this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
            }
          )
      )
  }

  private search(term: string): void {
    this.searchTerms.next(term)
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
