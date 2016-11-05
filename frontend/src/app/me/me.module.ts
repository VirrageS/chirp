import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { RouterModule }   from '@angular/router';
import { FormsModule }    from '@angular/forms';

import { MeComponent }        from './me.component';
import { MyTweetsComponent }  from './my-tweets.component';
import { FollowingComponent } from './following.component';
import { FollowersComponent } from './followers.component';

import { TweetsModule } from '../components';


@NgModule({
  imports: [
    CommonModule,
    FormsModule,

    RouterModule.forChild([
      { path: 'me', component: MeComponent,
        children: [
          { path: '', redirectTo: 'tweets' },
          { path: 'tweets', component: MyTweetsComponent },
          { path: 'following', component: FollowingComponent },
          { path: 'followers', component: FollowersComponent },
        ]
      }
    ]),

    TweetsModule,
  ],
  declarations: [
    MeComponent,
    MyTweetsComponent,
    FollowingComponent,
    FollowersComponent,
  ],
  providers: []
})
export class MeModule {}
