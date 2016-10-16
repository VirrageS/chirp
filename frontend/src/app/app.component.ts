import { Component } from '@angular/core';

import '../../public/css/styles.css';

@Component({
  selector: 'app',
  template: `
    <div>
      <navigation-bar></navigation-bar>
      <router-outlet></router-outlet>
    </div>
  `
})
export class AppComponent {}
