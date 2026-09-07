import type { ColorInput } from "@opentui/core";
import { Portal } from "@opentui/solid";
import { For } from "solid-js";
import { dismiss, notifications, type NotificationKind } from "../context/notifications.ts";
import { theme } from "../index.ts";

const KIND_LABEL: Record<NotificationKind, string> = {
  info: "Info",
  success: "Success",
  warning: "Warning",
  error: "Error",
};

export interface NotificationHostProps {
  width?: number;
  colors?: Partial<Record<NotificationKind, ColorInput>>;
  backgroundColor?: ColorInput;
}

export function NotificationHost(props: NotificationHostProps = {}) {
  const defaultColor = (): Record<NotificationKind, ColorInput> => ({
    info: theme.primary,
    success: theme.success,
    warning: theme.warning,
    error: theme.error,
  });
  const colorFor = (kind: NotificationKind) => props.colors?.[kind] ?? defaultColor()[kind];

  return (
    <Portal>
      <box position="absolute" top={1} right={1} flexDirection="column" zIndex={1000}>
        <For each={notifications()}>
          {(toast) => (
            <box
              border
              borderStyle={theme.borderStyle}
              borderColor={colorFor(toast.kind)}
              backgroundColor={props.backgroundColor ?? theme.background}
              width={props.width ?? 32}
              marginBottom={1}
              onMouseDown={theme.mouse ? () => dismiss(toast.id) : undefined}
            >
              <text fg={colorFor(toast.kind)}>{toast.title ?? KIND_LABEL[toast.kind]}</text>
              <text>{toast.message}</text>
            </box>
          )}
        </For>
      </box>
    </Portal>
  );
}
