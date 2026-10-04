import { TextAttributes } from "@opentui/core";
import { Show } from "solid-js";
import { theme } from "../index.ts";
import { Spinner } from "./Spinner.tsx";

export interface SectionHeaderProps {
  title?: string;
  loading?: boolean;
  /** Shown next to the spinner while loading (default "Loading…"). */
  loadingLabel?: string;
  /** Shown on the right when not loading, e.g. "r refresh". */
  hint?: string;
}

// Lives in its own entry point (ilovetui/components/section-header) because
// it pulls in the optional opentui-spinner peer dependency.
export function SectionHeader(props: SectionHeaderProps) {
  return (
    <box flexShrink={0} flexDirection="row" justifyContent="space-between" alignItems="center" height={1}>
      <box flexGrow={1} flexShrink={1} minWidth={0} marginRight={1} overflow="hidden">
        <text fg={theme.muted} attributes={TextAttributes.BOLD} wrapMode="none">
          {props.title ?? ""}
        </text>
      </box>
      <box flexShrink={0}>
        <Show
          when={!props.loading}
          fallback={
            <box flexDirection="row" alignItems="center">
              <Spinner />
              <text fg={theme.muted}> {props.loadingLabel ?? "Loading…"}</text>
            </box>
          }
        >
          <text fg={theme.muted}>{props.hint ?? ""}</text>
        </Show>
      </box>
    </box>
  );
}
