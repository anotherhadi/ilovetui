// @bun
// src/context/helper.tsx
import { createComponent as _$createComponent } from "@opentui/solid";
import { createContext, useContext } from "solid-js";
function createSimpleContext(input) {
  const ctx = createContext();
  return {
    provider: (props) => {
      const value = input.init(props);
      return _$createComponent(ctx.Provider, {
        value,
        get children() {
          return props.children;
        }
      });
    },
    use() {
      const value = useContext(ctx);
      if (!value)
        throw new Error(`${input.name} context must be used within its provider`);
      return value;
    }
  };
}

export { createSimpleContext };
