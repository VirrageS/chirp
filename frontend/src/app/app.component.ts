import { Component, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app',
  styles: [
    `
      .wrapper {
        max-width: 1000px;
        margin: 0 auto;
      }
    `,
    require('../../public/scss/main.scss'),
  ],
  encapsulation: ViewEncapsulation.None,
  template: `
    <div>
      <navigation-bar></navigation-bar>
      <div class="wrapper">
        <router-outlet></router-outlet>
      </div>
    </div>
  `
})
export class AppComponent {}
