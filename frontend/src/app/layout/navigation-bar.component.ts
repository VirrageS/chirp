import { Component } from '@angular/core';

import { AuthService, User } from '../shared';
import { Store } from '../store';


@Component({
  selector: 'navigation-bar',
  styleUrls: ['./navigation-bar.component.scss'],
  template: `
    <div class="navigation-bar shadow-1">
      <div class="menu">
        <a routerLink="/home" routerLinkActive="menu__link--active" class="menu__link">Home</a>
        <a routerLink="/me" *ngIf="_authenticated()" routerLinkActive="menu__link--active" class="menu__link">Me</a>
        <a routerLink="/search" *ngIf="_authenticated()" routerLinkActive="menu__link--active" class="menu__link">Find / Search</a>
      </div>

      <div class="">
        <div *ngIf="_authenticated()">
          <span>{{ user.name }}</span>
          <a routerLink="/logout" routerLinkActive="active" class="menu__link">Logout</a>
        </div>

        <div *ngIf="!_authenticated()">
          <a routerLink="/login" routerLinkActive="active" class="menu__link">Login</a>
        </div>
    </div>
  `
})
export class NavigationBarComponent {
  user?: User

  constructor(
    private _authService: AuthService,
    private _store: Store
  ) {
    this._store.changes.pluck('user')
      .subscribe((user: any) => this.user = user)
  }

  private _authenticated() {
    return this._authService.isAuthenticated()
  }
}
