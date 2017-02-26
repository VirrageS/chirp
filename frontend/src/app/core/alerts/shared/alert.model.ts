export enum AlertType {
  success,
  info,
  warning,
  danger,
}

export interface Alert {
  message: string;
  type: AlertType;
}
