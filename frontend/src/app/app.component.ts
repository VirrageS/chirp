import { Component } from '@angular/core';

import '../../public/css/styles.css';

@Component({
  selector: 'app',
  styles: [`
    .wrapper {
      max-width: 1000px;
      margin: 0 auto;
    }
  `],
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
