import { type ColorInput, CliRenderEvents, TextAttributes, type BoxRenderable } from "@opentui/core";
import { onResize, useRenderer } from "@opentui/solid";
import { formatCommandBindings } from "@opentui/keymap/extras";
import { useKeymapSelector } from "@opentui/keymap/solid";
import { createMemo, createSignal, For } from "solid-js";
import { theme } from "../index.ts";

const KEY_DISPLAY = { up: "↑", down: "↓", left: "←", right: "→", enter: "Enter", escape: "Esc" };
const SEPARATOR = "   ";
const ELLIPSIS = "…";

export interface HelpBarProps {
  accentColor?: ColorInput;
  mutedColor?: ColorInput;
}

interface HelpEntry {
  command: string;
  keys: string;
  label: string;
}

export function HelpBar(props: HelpBarProps = {}) {
  const accent = () => props.accentColor ?? theme.primary;
  const muted = () => props.mutedColor ?? theme.muted;
  const renderer = useRenderer();
  const [width, setWidth] = createSignal(0);
  let container: BoxRenderable | undefined;
  const measure = () => {
    if (!container) return;
    setWidth(container.width);
  };
  const remeasureNextFrame = () => renderer.once(CliRenderEvents.FRAME, measure);
  onResize(remeasureNextFrame);

  const entries = useKeymapSelector((keymap): HelpEntry[] =>
    keymap
      .getCommandEntries({ visibility: "active" })
      .filter((entry) => entry.command.shortHelp === true)
      .map((entry) => ({
        command: entry.command.name,
        keys: formatCommandBindings(entry.bindings, { keyNameAliases: KEY_DISPLAY }) ?? "",
        label: typeof entry.command.desc === "string" ? entry.command.desc : entry.command.name,
      })),
  );

  const visible = createMemo((): HelpEntry[] => {
    const available = width() - 2;
    const shown: HelpEntry[] = [];
    let used = 0;
    for (const entry of entries()) {
      const sepWidth = shown.length > 0 ? SEPARATOR.length : 0;
      const entryWidth = entry.keys.length + 1 + entry.label.length;
      if (used + sepWidth + entryWidth <= available) {
        used += sepWidth + entryWidth;
        shown.push(entry);
        continue;
      }
      const remaining = available - used - sepWidth;
      const prefix = `${entry.keys} `;
      if (remaining > prefix.length + ELLIPSIS.length) {
        const labelBudget = remaining - prefix.length - ELLIPSIS.length;
        shown.push({ ...entry, label: `${entry.label.slice(0, labelBudget)}${ELLIPSIS}` });
      }
      break;
    }
    return shown;
  });

  return (
    <box
      flexDirection="row"
      paddingLeft={1}
      paddingRight={1}
      overflow="hidden"
      ref={(el) => {
        container = el;
        remeasureNextFrame();
      }}
    >
      <For each={visible()}>
        {(entry, i) => (
          <>
            {i() > 0 && <text fg={muted()}>{SEPARATOR}</text>}
            <text fg={accent()} attributes={TextAttributes.BOLD}>
              {entry.keys}
            </text>
            <text fg={muted()}>{` ${entry.label}`}</text>
          </>
        )}
      </For>
    </box>
  );
}
