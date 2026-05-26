package ilovetui

type colorsYAML struct {
	Base00 string `yaml:"base00"`
	Base01 string `yaml:"base01"`
	Base02 string `yaml:"base02"`
	Base03 string `yaml:"base03"`
	Base04 string `yaml:"base04"`
	Base05 string `yaml:"base05"`
	Base06 string `yaml:"base06"`
	Base07 string `yaml:"base07"`
	Base08 string `yaml:"base08"`
	Base09 string `yaml:"base09"`
	Base0A string `yaml:"base0a"`
	Base0B string `yaml:"base0b"`
	Base0C string `yaml:"base0c"`
	Base0D string `yaml:"base0d"`
	Base0E string `yaml:"base0e"`
	Base0F string `yaml:"base0f"`
}

type configYAML struct {
	Colors colorsYAML `yaml:"colors"`
}

func mergeColors(base, user colorsYAML) colorsYAML {
	pick := func(b, u string) string {
		if u != "" {
			return u
		}
		return b
	}
	return colorsYAML{
		Base00: pick(base.Base00, user.Base00),
		Base01: pick(base.Base01, user.Base01),
		Base02: pick(base.Base02, user.Base02),
		Base03: pick(base.Base03, user.Base03),
		Base04: pick(base.Base04, user.Base04),
		Base05: pick(base.Base05, user.Base05),
		Base06: pick(base.Base06, user.Base06),
		Base07: pick(base.Base07, user.Base07),
		Base08: pick(base.Base08, user.Base08),
		Base09: pick(base.Base09, user.Base09),
		Base0A: pick(base.Base0A, user.Base0A),
		Base0B: pick(base.Base0B, user.Base0B),
		Base0C: pick(base.Base0C, user.Base0C),
		Base0D: pick(base.Base0D, user.Base0D),
		Base0E: pick(base.Base0E, user.Base0E),
		Base0F: pick(base.Base0F, user.Base0F),
	}
}
