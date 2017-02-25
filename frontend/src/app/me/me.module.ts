import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { RouterModule }   from '@angular/router';
import { FormsModule }    from '@angular/forms';

import { MeComponent }        from './me.component';
import { MyTweetsComponent }  from './my-tweets.component';
import { FollowingComponent } from './following.component';
import { FollowersComponent } from './followers.component';

import { TweetsModule } from '../tweets';
import { UsersModule } from '../users';

import { AuthService } from '../auth';


@NgModule({
  imports: [
    CommonModule,
    FormsModule,

    RouterModule.forChild([
      {
        path: 'me',
        component: MeComponent,
        canActivate: [AuthService],
        canActivateChild: [AuthService],
        children: [
          { path: '', redirectTo: 'tweets', pathMatch: 'prefix' },
          { path: 'tweets', component: MyTweetsComponent },
          { path: 'following', component: FollowingComponent },
          { path: 'followers', component: FollowersComponent },
        ]
      }
    ]),

    TweetsModule,
    UsersModule,
  ],
  declarations: [
    MeComponent,
    MyTweetsComponent,
    FollowingComponent,
    FollowersComponent,
  ],
  providers: [
    AuthService,
  ]
})
export class MeModule {}
