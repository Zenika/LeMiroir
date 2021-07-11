import { Notification } from "svelma";

const notify = (type, typeOptions) => (message, options) => Notification.create({
  type,
  message,
  ...typeOptions,
  ...options
});

export const success = notify('is-success')

export const warning = notify('is-warning')

export const error = notify('is-danger', { autoClose: false })

export const info = notify('is-info')
