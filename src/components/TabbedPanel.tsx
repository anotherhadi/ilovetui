import type { ParentProps } from "solid-js";
import { theme } from "../index.ts";
import { Box } from "./Box.tsx";
import { type TabItem, Tabs } from "./Tabs.tsx";

export interface TabbedPanelProps extends ParentProps {
  items: TabItem[];
  value: string;
  onChange: (value: string) => void;
  focused?: boolean;
}

/** <Tabs/> on top of a bordered panel whose top edge is drawn by the tabs. */
export function TabbedPanel(props: TabbedPanelProps) {
  return (
    <box flexDirection="column" flexGrow={1}>
      <Tabs items={props.items} value={props.value} onChange={props.onChange} focused={props.focused} />
      <Box
        flexGrow={1}
        border={["left", "right", "bottom"]}
        borderColor={props.focused ? theme.primary : theme.muted}
        flexDirection="column"
      >
        {props.children}
      </Box>
    </box>
  );
}
