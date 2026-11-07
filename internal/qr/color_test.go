package qr

import (
	"image/color"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    color.RGBA
		wantErr bool
	}{
		{"black with hash", "#000000", color.RGBA{0, 0, 0, 255}, false},
		{"white without hash", "FFFFFF", color.RGBA{255, 255, 255, 255}, false},
		{"red", "#FF0000", color.RGBA{255, 0, 0, 255}, false},
		{"mixed case", "#aaBBcc", color.RGBA{170, 187, 204, 255}, false},
		{"too short", "#FFF", color.RGBA{}, true},
		{"too long", "#FFFFFFF", color.RGBA{}, true},
		{"invalid hex", "#GGGGGG", color.RGBA{}, true},
		{"empty", "", color.RGBA{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHexColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHexColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseHexColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatHexColor(t *testing.T) {
	tests := []struct {
		input color.RGBA
		want  string
	}{
		{color.RGBA{0, 0, 0, 255}, "#000000"},
		{color.RGBA{255, 255, 255, 255}, "#ffffff"},
		{color.RGBA{255, 87, 51, 255}, "#ff5733"},
	}

	for _, tt := range tests {
		got := FormatHexColor(tt.input)
		if got != tt.want {
			t.Errorf("FormatHexColor(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestPresetsHaveValidColors(t *testing.T) {
	for name, pair := range Presets {
		if pair.Foreground.A != 255 {
			t.Errorf("preset %q foreground alpha should be 255, got %d", name, pair.Foreground.A)
		}
		if pair.Background.A != 255 {
			t.Errorf("preset %q background alpha should be 255, got %d", name, pair.Background.A)
		}
	}
}
