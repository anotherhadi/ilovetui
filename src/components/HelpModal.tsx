import { type ColorInput, TextAttributes } from "@opentui/core";
import { formatCommandBindings } from "@opentui/keymap/extras";
import { useKeymapSelector } from "@opentui/keymap/solid";
import { helpOpen, toggleHelp } from "../context/help.ts";
import { Modal } from "./Modal.tsx";

const KEY_DISPLAY = { up: "↑", down: "↓", left: "←", right: "→", enter: "Enter", escape: "Esc" };

export interface HelpModalProps {
  accentColor?: ColorInput;
  backgroundColor?: ColorInput;
  backdropColor?: ColorInput;
  width?: number;
}

interface HelpEntry {
  command: string;
  keys: string;
  label: string;
}

export function HelpModal(props: HelpModalProps = {}) {
  const entries = useKeymapSelector((keymap): HelpEntry[] =>
    keymap.getCommandEntries({ visibility: "active" }).map((entry) => ({
      command: entry.command.name,
      keys: formatCommandBindings(entry.bindings, { keyNameAliases: KEY_DISPLAY }) ?? "",
      label: typeof entry.command.desc === "string" ? entry.command.desc : entry.command.name,
    })),
  );

  return (
    <Modal
      open={helpOpen()}
      onDismiss={toggleHelp}
      accentColor={props.accentColor}
      backgroundColor={props.backgroundColor}
      backdropColor={props.backdropColor}
      width={props.width}
    >
      <text attributes={TextAttributes.BOLD}>Keybindings</text>
      <text> </text>
      <text>{entries().map((entry) => `${entry.keys.padEnd(10)}${entry.label}`).join("\n")}</text>
    </Modal>
  );
}
