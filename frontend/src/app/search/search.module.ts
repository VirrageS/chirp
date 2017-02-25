import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { RouterModule } from '@angular/router';

import { TweetsModule } from '../tweets';
import { UsersModule } from '../users';

import { SearchService } from './search.service';
import { SearchComponent } from './search.component';

import { AuthService } from '../auth';

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forChild([
      {
        path: 'search',
        component: SearchComponent,
        canActivate: [AuthService]
      },
    ]),

    TweetsModule,
    UsersModule,
  ],
  declarations: [
    SearchComponent,
  ],
  exports: [],
  providers: [
    AuthService,
    SearchService,
  ]
})
export class SearchModule {}
