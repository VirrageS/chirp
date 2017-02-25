import { Component } from '@angular/core';

import { AuthService } from '../../auth';
import { User } from '../../users';
import { Store } from '../../shared';


@Component({
  selector: 'nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss'],
})
export class NavComponent {
  user: User = null

  constructor(
    private _authService: AuthService,
    private _store: Store
  ) {
    this._store.changes("user")
      .subscribe((user: any) => this.user = user)
  }

  private _authenticated() {
    return this._authService.isAuthenticated()
  }
}
