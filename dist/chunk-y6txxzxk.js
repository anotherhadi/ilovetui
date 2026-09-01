// @bun
// src/context/notifications.ts
import { createSignal } from "solid-js";
var DEFAULT_DURATION = 3000;
var [notifications, setNotifications] = createSignal([]);
var nextId = 0;
function notify(message, options = {}) {
  const id = options.id ?? `toast-${nextId++}`;
  const toast = { id, title: options.title, message, kind: options.kind ?? "info" };
  setNotifications((prev) => [...prev.filter((n) => n.id !== id), toast]);
  const duration = options.duration ?? DEFAULT_DURATION;
  if (duration > 0) {
    setTimeout(() => dismiss(id), duration);
  }
  return id;
}
function dismiss(id) {
  setNotifications((prev) => prev.filter((n) => n.id !== id));
}

export { notifications, notify, dismiss };
