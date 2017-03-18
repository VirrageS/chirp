import { Component, Input, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';

import { Tweet, TweetService } from '../shared';
import { AlertType } from '../../core/alerts';
import { StoreHelper } from '../../shared';


@Component({
  selector: 'create-tweet',
  templateUrl: './create-tweet.component.html',
  styleUrls: ['./create-tweet.component.scss']
})
export class CreateTweetComponent implements OnInit {
  private tweet: Tweet = {
    content: ""
  }

  private tweetForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private storeHelper: StoreHelper,
    private tweetService: TweetService,
  ) {}

  ngOnInit(): void {
    this.buildForm();
  }

  onSubmit(): void {
    this.tweet = this.tweetForm.value;
    this.tweetService.createTweet(this.tweet)
      .subscribe(
        result => {
          this.storeHelper.add("alerts", {message: "Tweet has been added successfully", type: AlertType.success});
          this.tweet = { content: "" }
        },
        error => {
          this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
        }
      )
  }

  private buildForm(): void {
    this.tweetForm = this.fb.group({
      'content': [null, [
          Validators.required,
          Validators.maxLength(140),
        ]
      ],
    });
  }
}
