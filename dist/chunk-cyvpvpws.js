// @bun
// src/context/help.ts
import { createSignal } from "solid-js";
var [helpOpen, setHelpOpen] = createSignal(false);
function toggleHelp() {
  setHelpOpen((open) => !open);
}

export { helpOpen, toggleHelp };
