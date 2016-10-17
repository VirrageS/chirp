import { Component } from '@angular/core';

@Component({
  selector: 'navigation-bar',
  styles: [`
    .navigation-bar {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      align-items: center;

      height: 60px;
      margin-bottom: 30px;
      border-bottom: 1px;

      padding: 0px 5vw;
      border-bottom: 1px solid #c9c9c9;
    }

    .menu {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      align-items: center;

      font-weight: 500;
      color: #000000;
      width: 300px;
    }

    .menu__link {
      color: inherit;
      text-decoration: none;
    }

    .menu__link:hover, .menu__link--active {
      color: rgb(251, 104, 78);
    }
  `],
  template: `
    <div class="navigation-bar">
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
    name: 'Twutter'
  }
}
