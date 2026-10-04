// @bun
import"./chunk-njvj322f.js";
import"./chunk-6ggk5fg7.js";
import"./chunk-yyg967q8.js";
import"./chunk-xykm10y9.js";
import {
  theme2,
  presets2
} from "./chunk-p6jvscj3.js";

// src/solid.ts
import { extend } from "@opentui/solid";
import {
  ASCIIFontRenderable,
  BoxRenderable,
  InputRenderable,
  SelectRenderable,
  SliderRenderable,
  TabSelectRenderable,
  TextareaRenderable,
  TextRenderable
} from "@opentui/core";
function withThemeDefaults(Ctor, defaults, after) {
  class Themed extends Ctor {
    constructor(...args) {
      super(...args);
      Object.assign(this, defaults());
      after?.(this);
    }
  }
  return Themed;
}
var boxDefaults = () => {
  const { borderStyle, borderColor, focusedBorderColor } = presets2.box;
  return { borderStyle, borderColor, focusedBorderColor };
};
extend({
  box: withThemeDefaults(BoxRenderable, boxDefaults, (instance) => {
    instance.border = false;
  }),
  select: withThemeDefaults(SelectRenderable, () => presets2.select),
  tab_select: withThemeDefaults(TabSelectRenderable, () => presets2.tabSelect),
  input: withThemeDefaults(InputRenderable, () => presets2.input),
  textarea: withThemeDefaults(TextareaRenderable, () => presets2.textarea),
  slider: withThemeDefaults(SliderRenderable, () => presets2.slider),
  text: withThemeDefaults(TextRenderable, () => ({ fg: theme2.text })),
  ascii_font: withThemeDefaults(ASCIIFontRenderable, () => ({ color: theme2.primary }))
});
