import { Component, OnInit } from '@angular/core';
import { FormGroup, FormBuilder, Validators } from '@angular/forms';
import { Router }    from '@angular/router';

import { AlertType } from '../../core/alerts';
import { User } from '../../users';
import { SignupService } from './signup.service'
import { StoreHelper, CustomValidators } from '../../shared';


@Component({
  templateUrl: './signup.component.html',
  // styleUrls: ['./signup.component.scss'],
})
export class SignupComponent implements OnInit {
  private user: User = {
    name: "",
    username: "",
    email: "",
    password: ""
  }

  private signupForm: FormGroup;

  private formErrors = {
    'name': [],
    'username': [],
    'email': [],
    'password': [],
  };

  private validationMessages = {
    'name': {
      'required':  'Name is required.',
      'maxlength': 'Username cannot be more than 100 characters long.',
      'name':      'Name can contain only letters, apostrophes, dashes and spaces. Name should also contain two segments.',
    },
    'username': {
      'required':  'Username is required.',
      'maxlength': 'Username cannot be more than 30 characters long.',
      'username':  'Username can contain English letters and numbers.',
    },
    'email': {
      'required': 'Email is required.',
      'email':    'Email is invalid.',

    },
    'password': {
      'required':  'Password is required.',
      'minlength': 'Password must be at least 8 characters long.',
      'maxlength': 'Password cannot be more than 40 characters long.',
    }
  };

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private signupService: SignupService,
    private storeHelper: StoreHelper,
  ) {}

  ngOnInit(): void {
    this.buildForm();
  }

  private onSubmit() {
    this.user = this.signupForm.value;
    this.signupService.signup(this.user)
      .subscribe(
        result => {
          this.storeHelper.add("alerts", {message: "Account has been created successfully", type: AlertType.success});
          this.router.navigate(['', 'login'])
        },
        error => {
          this.storeHelper.add("alerts", {message: error, type: AlertType.danger});
        }
      )
  }

  private buildForm(): void {
    this.signupForm = this.fb.group({
      'name': [null, [
          Validators.required,
          Validators.maxLength(100),
          CustomValidators.fullname,
        ]
      ],
      'username': [null, [
          Validators.required,
          Validators.maxLength(30),
          CustomValidators.username,
        ]
      ],
      'email': [null, [
          Validators.required,
          CustomValidators.email,
        ]
      ],
      'password': [null, [
          Validators.required,
          Validators.minLength(8),
          Validators.maxLength(40),
        ]
      ],
    });

    this.signupForm.valueChanges
      .subscribe(() => this.onValueChanged());

    this.onValueChanged(); // (re)set validation messages now
  }

  private onValueChanged(): void {
    if (!this.signupForm) { return; }
    const form = this.signupForm;
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
