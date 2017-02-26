import { NgModule }       from '@angular/core';
import { CommonModule }   from '@angular/common';
import { FormsModule }    from '@angular/forms';

import { AlertComponent }       from './alert/alert.component';
import { AlertsComponent }      from './alerts.component';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
  ],
  declarations: [
    AlertComponent,
    AlertsComponent,
  ],
  exports: [
    AlertsComponent,
  ],
  providers: []
})
export class AlertsModule {}
