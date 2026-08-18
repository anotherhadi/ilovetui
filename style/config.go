package style

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
	Colors    colorsYAML `yaml:"colors"`
	NerdFonts bool       `yaml:"nerd_fonts"`
	Border    string     `yaml:"border"`
}

func pickString(base, user string) string {
	if user != "" {
		return user
	}
	return base
}

func mergeConfig(base, user configYAML) configYAML {
	return configYAML{
		Colors: mergeColors(base.Colors, user.Colors),

		NerdFonts: base.NerdFonts || user.NerdFonts,
		Border:    pickString(base.Border, user.Border),
	}
}

func mergeColors(base, user colorsYAML) colorsYAML {
	return colorsYAML{
		Base00: pickString(base.Base00, user.Base00),
		Base01: pickString(base.Base01, user.Base01),
		Base02: pickString(base.Base02, user.Base02),
		Base03: pickString(base.Base03, user.Base03),
		Base04: pickString(base.Base04, user.Base04),
		Base05: pickString(base.Base05, user.Base05),
		Base06: pickString(base.Base06, user.Base06),
		Base07: pickString(base.Base07, user.Base07),
		Base08: pickString(base.Base08, user.Base08),
		Base09: pickString(base.Base09, user.Base09),
		Base0A: pickString(base.Base0A, user.Base0A),
		Base0B: pickString(base.Base0B, user.Base0B),
		Base0C: pickString(base.Base0C, user.Base0C),
		Base0D: pickString(base.Base0D, user.Base0D),
		Base0E: pickString(base.Base0E, user.Base0E),
		Base0F: pickString(base.Base0F, user.Base0F),
	}
}
