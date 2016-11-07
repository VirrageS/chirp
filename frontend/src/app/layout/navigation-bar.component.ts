import { Component } from '@angular/core';

@Component({
  selector: 'navigation-bar',
  styleUrls: ['./navigation-bar.component.scss'],
  template: `
    <div class="navigation-bar shadow-1">
      <div class="menu">
        <a routerLink="/home" routerLinkActive="menu__link--active" class="menu__link">Home</a>
        <a routerLink="/me" routerLinkActive="menu__link--active" class="menu__link">Me</a>
        <a routerLink="/search" routerLinkActive="menu__link--active" class="menu__link">Find / Search</a>
      </div>

      <div class="">
        <div *ngIf="user">
          <span>{{ user.name }}</span>
          <a routerLink="/logout" routerLinkActive="active" class="menu__link">Logout</a>
        </div>

        <div *ngIf="!user">
          <a routerLink="/login" routerLinkActive="active" class="menu__link">Login</a>
        </div>
    </div>
  `
})
export class NavigationBarComponent {
  user = {
    name: 'Chirp'
  }
}
