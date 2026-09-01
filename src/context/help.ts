import { createSignal } from "solid-js";

const [helpOpen, setHelpOpen] = createSignal(false);
export { helpOpen };

export function toggleHelp(): void {
  setHelpOpen((open) => !open);
}
