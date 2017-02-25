import { Component, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app',
  styles: [ // TODO
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
      <nav></nav>
      <div class="wrapper">
        <router-outlet></router-outlet>
      </div>
    </div>
  `
})
export class AppComponent {}
