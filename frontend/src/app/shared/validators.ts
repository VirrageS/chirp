import { AbstractControl, ValidatorFn } from '@angular/forms';

export class CustomValidators {
  // static forbiddenCharactersValidator(nameRe: RegExp): ValidatorFn {
  //   return (control: AbstractControl): {[key: string]: any} => {
  //     const name = control.value;
  //     const no = nameRe.test(name);
  //     return no ? {'forbiddenCharacters': {name}} : null;
  //   };
  // }

  static fullname(control: AbstractControl): {[key: string]: any} {
    const nameRegExp: RegExp = /^[A-Z]([-']?[a-z]+)*( [A-Z]([-']?[a-z]+)*)+$/;

    const name = control.value;

    // we allow empty email
    if ((name == null) || (name == ""))
      return null;

    const ok = nameRegExp.test(name);
    return !ok ? {'name': {name}} : null;
  }

  static username(control: AbstractControl): {[key: string]: any} {
    const usernameRegExp: RegExp = /^[a-zA-Z0-9]*$/;

    const name = control.value;

    // we allow empty email
    if ((name == null) || (name == ""))
      return null;

    const ok = usernameRegExp.test(name);
    return !ok ? {'username': {name}} : null;
  }

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
