import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';

import { Alert, AlertType } from '../shared';


@Component({
  selector: 'alert',
  templateUrl: './alert.component.html',
  // styleUrls: ['./alert.component.scss'],
})
export class AlertComponent implements OnInit {
  @Input() alert: Alert;
  @Output() close = new EventEmitter();
  private alertTypeClass: string;

  ngOnInit() {
    this.alertTypeClass = AlertType[this.alert.type];
  }

  private dismiss() {
    this.close.emit(this.alert);
  }
}
