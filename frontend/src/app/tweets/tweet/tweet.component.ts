import { Component, Input, OnDestroy, OnInit, Output, EventEmitter } from '@angular/core';

import { Tweet, TweetService } from '../shared';
import { User, UserService } from '../../users/shared';
import { Store } from '../../shared';
import * as _ from 'lodash';
import * as moment from 'moment';


@Component({
  selector: 'tweet',
  templateUrl: './tweet.component.html',
  styleUrls: ['./tweet.component.scss']
})
export class TweetComponent implements OnInit, OnDestroy {
  @Input() tweet: Tweet
  @Output() tweetChange = new EventEmitter()
  private loggedUser?: User
  private time: any = {
    difference: "",
    handler: {}
  }

  constructor(
    private _tweetService: TweetService,
    private _userService: UserService,
    private _store: Store
  ) {
    this._store.changes("user")
      .subscribe((user: any) => this.loggedUser = user)
  }

  ngOnInit(): void {
    this.updateTimeDifferrence()
  }

  ngOnDestroy(): void {
    clearTimeout(this.time.handler)
  }

  private updateTimeDifferrence(): void {
    this.time.difference = moment(this.tweet.created_at).fromNow()

    // update each minute plus some random to make tweets update
    // not fire at the same time (because this could lead to performance issues)
    // TODO: maybe consider changing 60 seconds to interval depending on actual
    // difference, so when difference is 1 hour, next update should be like every 20 minutes
    this.time.handler = setTimeout(() => {
      this.updateTimeDifferrence()
    }, (60 + Math.random() * 10) * 1000)
  }

  private _toggleFollow() {
    this.tweet.author.following = !this.tweet.author.following

    // send real request
    let toggleFunc = this._userService.follow(this.tweet.author.id)
    if (!this.tweet.author.following) {
      toggleFunc = this._userService.unfollow(this.tweet.author.id)
    }

    this.tweetChange.emit(this.tweet)
    toggleFunc
      .subscribe(author => {
        _.assign(this.tweet.author, author)
        this.tweetChange.emit(this.tweet)
      })
  }

  private _toggleLike() {
    this.tweet.liked = !this.tweet.liked

    let toggleFunc = this._tweetService.like(this.tweet.id)
    if (!this.tweet.liked) {
      this.tweet.like_count -= 1
      toggleFunc = this._tweetService.unlike(this.tweet.id)
    } else {
      this.tweet.like_count += 1
    }

    toggleFunc
      .subscribe(
        result => _.assign(this.tweet, result),
        error => {}
      )
  }

  // TODO:
  // private _retweet() {
  //   this.tweet.retweeted = true
  //   this._tweetService.retweet(this.tweet.id)
  //     .subscribe(
  //       result => _.assign(this.tweet, result),
  //       error => {}
  //     )
  // }
}
