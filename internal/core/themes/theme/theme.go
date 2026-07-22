package theme

import "encoding/json"

type FontConfig struct {
	Family string  `json:"family"`
	Size   float64 `json:"size"`
}

type Fonts struct {
	System   FontConfig `json:"system"`
	Widget   FontConfig `json:"widget"`
	Terminal FontConfig `json:"terminal"`
}

func (f *Fonts) ResolveDefaults() {
	if f.System.Size == 0 {
		f.System.Size = 11.0
	}
	if f.Widget.Size == 0 {
		f.Widget.Size = f.System.Size
	}
	if f.Terminal.Size == 0 {
		f.Terminal.Size = 11.0
	}

	if f.Widget.Family == "" {
		f.Widget.Family = f.System.Family
	}
	if f.Terminal.Family == "" {
		f.Terminal.Family = f.System.Family
	}
}

type Theme struct {
	Theme struct {
		Name  string `json:"name"`
		Fonts Fonts  `json:"fonts"`
	} `json:"theme"`

	Tools map[string]json.RawMessage `json:"tools"`
}
