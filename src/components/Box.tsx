import type { BoxProps as OpenTUIBoxProps } from "@opentui/solid";
import { theme } from "../index.ts";

export interface BoxProps extends OpenTUIBoxProps {}

// The raw <box> intrinsic has no idea ilovetui's theme exists — its own
// borderStyle default ("single") is hardcoded in @opentui/core. Use <Box>
// instead of <box> wherever you want theme.borderStyle applied by default,
// same as Sidebar/Modal/Tabs already do internally.
export function Box(props: BoxProps) {
  return <box {...props} borderStyle={props.borderStyle ?? theme.borderStyle} />;
}
