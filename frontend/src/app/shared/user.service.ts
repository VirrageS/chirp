import { Injectable } from '@angular/core';
import { ApiService } from './api.service';

@Injectable()
export class UserService {
  user_path: string = "/user"
  user_id: string = ""

  constructor(private apiService: ApiService) {
    // TODO: get user_id
  }

  getUser() {
    return this.apiService.get(this.user_path + this.user_id);
  }

  singupUser(body) {
    let path: string = "/singup";
    return this.apiService.post(path, body);
  }

  getTweets() {
    let path: string = "/tweets";
    return this.apiService.get(this.user_path + "/" + this.user_id + path);
  }
}
