import { Injectable } from '@angular/core';

import { ApiService } from '../shared';


@Injectable()
export class SearchService {
  private search_path: string = "/search";

  constructor(
    private _apiService: ApiService
  ) {}

  search(text: string) {
    return this._apiService.get(this.search_path + "?querystring=" + text);
  }
}
