import { AbstractControl, ValidatorFn } from '@angular/forms';

export class CustomValidators {
  // static forbiddenCharactersValidator(nameRe: RegExp): ValidatorFn {
  //   return (control: AbstractControl): {[key: string]: any} => {
  //     const name = control.value;
  //     const no = nameRe.test(name);
  //     return no ? {'forbiddenCharacters': {name}} : null;
  //   };
  // }

  static email(control: AbstractControl): {[key: string]: any} {
    const emailRegExp: RegExp = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

    const name = control.value;

    // we allow empty email
    if ((name == null) || (name == ""))
      return null;

    const ok = emailRegExp.test(name);
    return !ok ? {'email': {name}} : null;
  }
}
