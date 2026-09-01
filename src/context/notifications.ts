import { createSignal } from "solid-js";

export type NotificationKind = "info" | "success" | "warning" | "error";

export interface Notification {
  id: string;
  title?: string;
  message: string;
  kind: NotificationKind;
}

export interface NotifyOptions {
  id?: string;
  title?: string;
  kind?: NotificationKind;
  duration?: number;
}

const DEFAULT_DURATION = 3000;

const [notifications, setNotifications] = createSignal<Notification[]>([]);
export { notifications };

let nextId = 0;

export function notify(message: string, options: NotifyOptions = {}): string {
  const id = options.id ?? `toast-${nextId++}`;
  const toast: Notification = { id, title: options.title, message, kind: options.kind ?? "info" };

  setNotifications((prev) => [...prev.filter((n) => n.id !== id), toast]);

  const duration = options.duration ?? DEFAULT_DURATION;
  if (duration > 0) {
    setTimeout(() => dismiss(id), duration);
  }

  return id;
}

export function dismiss(id: string): void {
  setNotifications((prev) => prev.filter((n) => n.id !== id));
}
