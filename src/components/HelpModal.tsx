import { type ColorInput, type Renderable, TextAttributes } from "@opentui/core";
import { formatCommandBindings } from "@opentui/keymap/extras";
import { useKeymapSelector } from "@opentui/keymap/solid";
import { createSignal } from "solid-js";
import { helpOpen, toggleHelp } from "../context/help.ts";
import { Modal } from "./Modal.tsx";

const KEY_DISPLAY = { up: "↑", down: "↓", left: "←", right: "→", enter: "Enter", escape: "Esc" };

export interface HelpModalProps {
  accentColor?: ColorInput;
  backgroundColor?: ColorInput;
  backdropColor?: ColorInput;
  width?: number;
  /** See ModalProps.dismissKey. */
  dismissKey?: string | false;
}

interface HelpEntry {
  command: string;
  keys: string;
  label: string;
}

export function HelpModal(props: HelpModalProps = {}) {
  // The modal takes focus when it opens, which would make its own "close"
  // the only active command. List what was active where focus was *before*
  // it opened instead.
  const [scope, setScope] = createSignal<Renderable | null>(null);

  const entries = useKeymapSelector((keymap): HelpEntry[] =>
    keymap.getCommandEntries({ visibility: "active", focused: scope() }).map((entry) => ({
      command: entry.command.name,
      keys: formatCommandBindings(entry.bindings, { keyNameAliases: KEY_DISPLAY }) ?? "",
      label: typeof entry.command.desc === "string" ? entry.command.desc : entry.command.name,
    })),
  );

  return (
    <Modal
      open={helpOpen()}
      onDismiss={toggleHelp}
      onOpen={(previous) => setScope(() => previous)}
      dismissKey={props.dismissKey}
      accentColor={props.accentColor}
      backgroundColor={props.backgroundColor}
      backdropColor={props.backdropColor}
      width={props.width}
      maxHeight={entries().length + 6}
      scrollable
    >
      <text attributes={TextAttributes.BOLD}>Keybindings</text>
      <text> </text>
      <text>{entries().map((entry) => `${entry.keys.padEnd(10)}${entry.label}`).join("\n")}</text>
    </Modal>
  );
}
