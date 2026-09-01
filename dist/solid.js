// @bun
import {
  presets,
  theme
} from "./chunk-57mam43k.js";
import"./chunk-ergc99dh.js";
import"./chunk-7f8jagy5.js";
import"./chunk-hcq62p48.js";

// src/solid.ts
import { extend } from "@opentui/solid";
import {
  ASCIIFontRenderable,
  BoxRenderable,
  InputRenderable,
  SelectRenderable,
  TabSelectRenderable,
  TextareaRenderable,
  TextRenderable
} from "@opentui/core";
function withThemeDefaults(Ctor, defaults, after) {
  class Themed extends Ctor {
    constructor(...args) {
      super(...args);
      Object.assign(this, defaults);
      after?.(this);
    }
  }
  return Themed;
}
var { borderStyle, borderColor, focusedBorderColor } = presets.box;
var boxDefaults = { borderStyle, borderColor, focusedBorderColor };
extend({
  box: withThemeDefaults(BoxRenderable, boxDefaults, (instance) => {
    instance.border = false;
  }),
  select: withThemeDefaults(SelectRenderable, presets.select),
  tab_select: withThemeDefaults(TabSelectRenderable, presets.tabSelect),
  input: withThemeDefaults(InputRenderable, presets.input),
  textarea: withThemeDefaults(TextareaRenderable, presets.textarea),
  text: withThemeDefaults(TextRenderable, { fg: theme.text }),
  ascii_font: withThemeDefaults(ASCIIFontRenderable, { color: theme.primary })
});
