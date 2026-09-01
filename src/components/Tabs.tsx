import { type BorderCharacters, type BoxRenderable, type ColorInput, BorderChars, CliRenderEvents, TextAttributes } from "@opentui/core";
import { onResize, useRenderer } from "@opentui/solid";
import { createMemo, createSignal, For } from "solid-js";
import { theme } from "../index.ts";

export interface TabItem {
  value: string;
  label: string;
}

export interface TabsProps {
  items: TabItem[];
  value: string;
  onChange: (value: string) => void;
  focused?: boolean;
  accentColor?: ColorInput;
  mutedColor?: ColorInput;
}

export function stepTabValue(items: TabItem[], value: string, delta: number): string {
  const index = items.findIndex((t) => t.value === value);
  if (index === -1) return value;
  return items[(index + delta + items.length) % items.length]!.value;
}

function innerWidth(label: string): number {
  return label.length + 2;
}

function segmentWidth(label: string): number {
  return innerWidth(label) + 2;
}

function topEdge(chars: BorderCharacters, label: string): string {
  return `${chars.topLeft}${chars.horizontal.repeat(innerWidth(label))}${chars.topRight}`;
}

function bottomEdge(chars: BorderCharacters, label: string, isFirst: boolean, isLast: boolean, isActive: boolean): string {
  const left = isFirst ? (isActive ? chars.vertical : chars.leftT) : isActive ? chars.bottomRight : chars.bottomT;
  const right = isLast ? (isActive ? chars.vertical : chars.rightT) : isActive ? chars.bottomLeft : chars.bottomT;
  const fill = isActive ? " " : chars.horizontal;
  return left + fill.repeat(innerWidth(label)) + right;
}

function buildBottomRow(
  chars: BorderCharacters,
  segments: { label: string; isActive: boolean }[],
  available: number,
): string {
  const tabsWidth = segments.reduce((sum, s) => sum + segmentWidth(s.label), 0);
  const remainder = available - tabsWidth;
  const extending = remainder > 0;

  const tabs = segments
    .map((s, i) => bottomEdge(chars, s.label, i === 0, i === segments.length - 1 && !extending, s.isActive))
    .join("");

  return extending ? tabs + chars.horizontal.repeat(remainder - 1) + chars.topRight : tabs;
}

function fitCount(items: TabItem[], start: number, available: number): number {
  const hiddenBefore = start;
  let used = 0;
  let count = 0;
  for (let i = start; i < items.length; i++) {
    const width = segmentWidth(items[i]!.label);
    const hiddenAfter = items.length - start - count - 1;
    const totalHidden = hiddenBefore + hiddenAfter;
    const reserve = totalHidden > 0 ? segmentWidth(`+${totalHidden}`) : 0;
    if (used + width + reserve > available) break;
    used += width;
    count++;
  }
  return Math.max(count, 1);
}

export function Tabs(props: TabsProps) {
  const renderer = useRenderer();
  const [width, setWidth] = createSignal(0);
  let container: BoxRenderable | undefined;
  const measure = () => {
    if (!container) return;
    setWidth(container.width);
  };
  const remeasureNextFrame = () => renderer.once(CliRenderEvents.FRAME, measure);
  onResize(remeasureNextFrame);

  const segments = createMemo(() => {
    const items = props.items;
    const activeIndex = Math.max(
      0,
      items.findIndex((t) => t.value === props.value),
    );

    let start = 0;
    let count = fitCount(items, start, width());
    if (activeIndex >= start + count) {
      start = activeIndex;
      count = fitCount(items, start, width());
    }

    const shownItems = items.slice(start, start + count);
    const shown = shownItems.map((item) => ({
      value: item.value as string | undefined,
      label: item.label,
      isActive: item.value === props.value,
    }));

    const hidden = items.length - count;
    if (hidden > 0) {
      const used = shownItems.reduce((sum, item) => sum + segmentWidth(item.label), 0);
      const tag = `+${hidden}`;
      if (shown.length === 0 || used + segmentWidth(tag) <= width()) {
        shown.push({ value: undefined, label: tag, isActive: false });
      }
    }
    return shown;
  });

  const accent = () => props.accentColor ?? theme.primary;
  const muted = () => props.mutedColor ?? theme.muted;

  const borderColor = () => (props.focused ? accent() : muted());
  const chars = BorderChars[theme.borderStyle];

  return (
    <box
      flexDirection="column"
      ref={(el) => {
        container = el;
        remeasureNextFrame();
      }}
    >
      <text selectable={false} fg={borderColor()}>
        {segments()
          .map((s) => topEdge(chars, s.label))
          .join("")}
      </text>
      <box flexDirection="row">
        <For each={segments()}>
          {(s) => (
            <>
              <text selectable={false} fg={borderColor()}>
                {chars.vertical}
              </text>
              <text
                selectable={false}
                fg={s.isActive ? accent() : muted()}
                attributes={s.isActive ? TextAttributes.BOLD : undefined}
                onMouseDown={theme.mouse && s.value !== undefined ? () => props.onChange(s.value!) : undefined}
              >
                {` ${s.label} `}
              </text>
              <text selectable={false} fg={borderColor()}>
                {chars.vertical}
              </text>
            </>
          )}
        </For>
      </box>
      <text selectable={false} fg={borderColor()}>
        {buildBottomRow(chars, segments(), width())}
      </text>
    </box>
  );
}
