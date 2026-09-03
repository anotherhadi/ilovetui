import { extend } from "@opentui/solid";
import {
  ASCIIFontRenderable,
  BoxRenderable,
  InputRenderable,
  SelectRenderable,
  SliderRenderable,
  TabSelectRenderable,
  TextareaRenderable,
  TextRenderable,
} from "@opentui/core";
import { presets, theme } from "./index.ts";

declare module "@opentui/solid" {
  interface OpenTUIComponents {
    slider: typeof SliderRenderable;
  }
}

type AnyCtor = new (...args: any[]) => any;

function withThemeDefaults<T extends AnyCtor>(
  Ctor: T,
  defaults: object,
  after?: (instance: any) => void,
): T {
  class Themed extends Ctor {
    constructor(...args: any[]) {
      super(...args);
      Object.assign(this, defaults);
      after?.(this);
    }
  }
  return Themed as T;
}

const { borderStyle, borderColor, focusedBorderColor } = presets.box;
const boxDefaults = { borderStyle, borderColor, focusedBorderColor };

extend({
  box: withThemeDefaults(BoxRenderable, boxDefaults, (instance) => {
    instance.border = false;
  }),
  select: withThemeDefaults(SelectRenderable, presets.select),
  tab_select: withThemeDefaults(TabSelectRenderable, presets.tabSelect),
  input: withThemeDefaults(InputRenderable, presets.input),
  textarea: withThemeDefaults(TextareaRenderable, presets.textarea),
  slider: withThemeDefaults(SliderRenderable, presets.slider),
  text: withThemeDefaults(TextRenderable, { fg: theme.text }),
  ascii_font: withThemeDefaults(ASCIIFontRenderable, { color: theme.primary }),
});
