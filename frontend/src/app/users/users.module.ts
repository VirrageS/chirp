import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { FormsModule }    from '@angular/forms';

import { UserService } from './shared';
import { UserComponent }       from './user/user.component';
import { UsersComponent }      from './users.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
  ],
  declarations: [
    UserComponent,
    UsersComponent,
  ],
  exports: [
    UsersComponent,
    UserService,
  ],
  providers: [
    UserService,
  ]
})
export class UsersModule {}
