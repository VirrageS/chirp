import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { RouterModule }   from '@angular/router';
import { FormsModule }    from '@angular/forms';

import { MeComponent }        from './me.component';
import { TweetsComponent }    from './tweets.component';
import { FollowingComponent } from './following.component';
import { FollowersComponent } from './followers.component';

// import { HeroService } from './hero.service';


@NgModule({
  imports: [
    CommonModule,
    FormsModule,

    RouterModule.forChild([
      { path: 'me', component: MeComponent,
        children: [
          { path: '', redirectTo: 'tweets' },
          { path: 'tweets', component: TweetsComponent },
          { path: 'following', component: FollowingComponent },
          { path: 'followers', component: FollowersComponent },
        ]
      }
    ])
  ],
  declarations: [
    MeComponent,
    TweetsComponent,
    FollowingComponent,
    FollowersComponent
  ],
  providers: [
    // HeroService
  ]
})
export class MeModule {}
