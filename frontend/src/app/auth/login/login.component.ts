import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Router } from '@angular/router';

import { AlertType } from '../../core/alerts';
import { User } from '../../users';
import { LoginService } from './login.service';
import { StoreHelper, CustomValidators } from '../../shared';


@Component({
  templateUrl: './login.component.html',
  // styleUrls: ['./login.component.scss'],
})
export class LoginComponent implements OnInit {
  private user: User = {
    email: "",
    password: ""
  }
  private loginForm: FormGroup;

  private formErrors = {
    'email': [],
    'password': [],
  };

  private validationMessages = {
    'email': {
      'required':  'Email is required.',
      'email':     'Email is invalid.',
    },
    'password': {
      'required':      'Password is required.',
      'minlength':     'Password must be at least 8 characters long.',
      'maxlength':     'Password cannot be more than 20 characters long.',
    }
  };

  constructor(
    private fb: FormBuilder,
    private loginService: LoginService,
    private router: Router,
    private storeHelper: StoreHelper,
  ) {}

  ngOnInit(): void {
    this.buildForm();
  }

  private loginWithGoogle(): void {
    this.loginService.authorizeWithGoogle() // firstly we have to authorize then login
      .subscribe(
        result => { window.location.href = result }, // redirect to Google AUTH
        error => {
          this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
        }
      );
  }

  private onSubmit(): void {
    this.user = this.loginForm.value;
    this.loginService.login(this.user)
      .subscribe(
        result => {
          this.router.navigateByUrl('home')
        },
        error => {
          this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
        }
      );
  }

  private buildForm(): void {
    this.loginForm = this.fb.group({
      'email': [null, [
          Validators.required,
          CustomValidators.email,
        ]
      ],
      'password': [null, [
          Validators.required,
          Validators.minLength(8),
          Validators.maxLength(20),
          // CustomValidators.forbiddenCharactersValidator(/bob/i),
        ]
      ],
    });

    this.loginForm.valueChanges
      .subscribe(() => this.onValueChanged());

    this.onValueChanged(); // (re)set validation messages now
  }

  private onValueChanged(): void {
    if (!this.loginForm) { return; }
    const form = this.loginForm;
    for (const field in this.formErrors) {
      // clear previous error message (if any)
      this.formErrors[field] = [];
      const control = form.get(field);

      if (control && control.dirty && !control.valid) {
        const messages = this.validationMessages[field];
        for (const key in control.errors) {
          this.formErrors[field].push(messages[key]);
        }
      }
    }
  }
}
