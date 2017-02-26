import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { Store } from '../../shared';
import { Alert, AlertType } from './shared';
import { StoreHelper } from '../../shared';
import * as _ from 'lodash';


@Component({
  selector: 'alerts',
  templateUrl: './alerts.component.html',
  styleUrls: ['./alerts.component.scss']
})
export class AlertsComponent implements OnInit {
  private alerts: Array<Alert> = [];

  constructor(
    private router: Router,
    private store: Store,
    private storeHelper: StoreHelper,
  ) {}

  ngOnInit(): void {
    this.store.changes("alerts")
      .subscribe((alerts: Array<Alert>) => this.alerts = alerts);

    // when changing router we should clear alerts
    this.router.events.subscribe(
      () => {this.storeHelper.update("alerts", [])}
    );
  }

  private handleClose(alert: Alert): void {
    this.storeHelper.findAndDelete("alerts", alert);
  }
}
