import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { FormsModule }    from '@angular/forms';

import { UserComponent }       from './user.component';
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
  ],
  providers: []
})
export class UsersModule {}
