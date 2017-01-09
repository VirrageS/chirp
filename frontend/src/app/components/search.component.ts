import { Component, OnInit } from '@angular/core';
import { Subject } from 'rxjs/Subject';

import { User, Tweet, SearchService } from '../shared';


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
            error => {
              console.log(error)
            }
          )
      )
  }

  search(term: string): void {
    this._searchTerms.next(term)
  }
}
