// @bun
// src/context/notifications.ts
import { createSignal } from "solid-js";
var DEFAULT_DURATION = 3000;
var [notifications2, setNotifications] = createSignal([]);
var nextId = 0;
function notify2(message, options = {}) {
  const id = options.id ?? `toast-${nextId++}`;
  const toast = { id, title: options.title, message, kind: options.kind ?? "info" };
  setNotifications((prev) => [...prev.filter((n) => n.id !== id), toast]);
  const duration = options.duration ?? DEFAULT_DURATION;
  if (duration > 0) {
    setTimeout(() => dismiss2(id), duration);
  }
  return id;
}
function dismiss2(id) {
  setNotifications((prev) => prev.filter((n) => n.id !== id));
}

export { notifications2, notify2, dismiss2 };
