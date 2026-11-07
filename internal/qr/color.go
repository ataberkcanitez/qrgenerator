package qr

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type ColorPair struct {
	Foreground color.RGBA
	Background color.RGBA
}

var Presets = map[string]ColorPair{
	"Classic (Black on White)": {color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}},
	"Dark Blue on White":      {color.RGBA{0, 0, 139, 255}, color.RGBA{255, 255, 255, 255}},
	"Dark Green on White":     {color.RGBA{0, 100, 0, 255}, color.RGBA{255, 255, 255, 255}},
	"Dark Red on White":       {color.RGBA{139, 0, 0, 255}, color.RGBA{255, 255, 255, 255}},
	"White on Black":          {color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255}},
}

var PresetNames = []string{
	"Classic (Black on White)",
	"Dark Blue on White",
	"Dark Green on White",
	"Dark Red on White",
	"White on Black",
	"Custom (enter hex colors)",
}

func ParseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("expected 6-digit hex color, got %q", s)
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid red component: %w", err)
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid green component: %w", err)
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid blue component: %w", err)
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

func FormatHexColor(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
