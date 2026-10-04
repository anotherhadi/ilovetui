// @bun
// src/context/help.ts
import { createSignal } from "solid-js";
var [helpOpen2, setHelpOpen] = createSignal(false);
function toggleHelp2() {
  setHelpOpen((open) => !open);
}

export { helpOpen2, toggleHelp2 };
