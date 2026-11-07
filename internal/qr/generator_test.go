package qr

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var defaultColors = ColorPair{
	Foreground: color.RGBA{0, 0, 0, 255},
	Background: color.RGBA{255, 255, 255, 255},
}

func TestGeneratePNG(t *testing.T) {
	out := filepath.Join(t.TempDir(), "test.png")
	opts := Options{
		Content: "https://example.com",
		Size:    256,
		Colors:  defaultColors,
	}

	if err := GeneratePNG(opts, out); err != nil {
		t.Fatalf("GeneratePNG failed: %v", err)
	}

	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("output file missing: %v", err)
	}
	if fi.Size() == 0 {
		t.Fatal("output file is empty")
	}
}

func TestGeneratePNGWithCustomColors(t *testing.T) {
	out := filepath.Join(t.TempDir(), "test_color.png")
	opts := Options{
		Content: "https://example.com",
		Size:    256,
		Colors: ColorPair{
			Foreground: color.RGBA{0, 0, 139, 255},
			Background: color.RGBA{255, 255, 255, 255},
		},
	}

	if err := GeneratePNG(opts, out); err != nil {
		t.Fatalf("GeneratePNG with custom colors failed: %v", err)
	}

	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("output file missing: %v", err)
	}
	if fi.Size() == 0 {
		t.Fatal("output file is empty")
	}
}

func TestGenerateSVG(t *testing.T) {
	out := filepath.Join(t.TempDir(), "test.svg")
	opts := Options{
		Content: "https://example.com",
		Colors:  defaultColors,
	}

	if err := GenerateSVG(opts, out); err != nil {
		t.Fatalf("GenerateSVG failed: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	svg := string(data)
	if !strings.HasPrefix(svg, "<svg") {
		t.Error("SVG output should start with <svg")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Error("SVG output should contain closing </svg> tag")
	}
	if !strings.Contains(svg, `fill="#000000"`) {
		t.Error("SVG should contain foreground color rects")
	}
}

func TestGenerateSVGWithColors(t *testing.T) {
	out := filepath.Join(t.TempDir(), "test_color.svg")
	opts := Options{
		Content: "hello",
		Colors: ColorPair{
			Foreground: color.RGBA{255, 0, 0, 255},
			Background: color.RGBA{0, 255, 0, 255},
		},
	}

	if err := GenerateSVG(opts, out); err != nil {
		t.Fatalf("GenerateSVG failed: %v", err)
	}

	data, _ := os.ReadFile(out)
	svg := string(data)
	if !strings.Contains(svg, `fill="#ff0000"`) {
		t.Error("SVG should use custom foreground color")
	}
	if !strings.Contains(svg, `fill="#00ff00"`) {
		t.Error("SVG should use custom background color")
	}
}

func TestGeneratePNGWithLogo(t *testing.T) {
	logoPath := filepath.Join(t.TempDir(), "logo.png")
	createTestLogo(t, logoPath)

	out := filepath.Join(t.TempDir(), "test_logo.png")
	opts := Options{
		Content:  "https://example.com",
		Size:     256,
		Colors:   defaultColors,
		LogoPath: logoPath,
	}

	if err := GeneratePNGWithLogo(opts, out); err != nil {
		t.Fatalf("GeneratePNGWithLogo failed: %v", err)
	}

	fi, err := os.Stat(out)
	if err != nil {
		t.Fatalf("output file missing: %v", err)
	}
	if fi.Size() == 0 {
		t.Fatal("output file is empty")
	}
}

func TestGeneratePNGInvalidContent(t *testing.T) {
	out := filepath.Join(t.TempDir(), "test.png")
	opts := Options{Content: "", Size: 256, Colors: defaultColors}

	err := GeneratePNG(opts, out)
	if err == nil {
		t.Error("expected error for empty content")
	}
}

func TestGeneratePNGWithLogoInvalidPath(t *testing.T) {
	out := filepath.Join(t.TempDir(), "test.png")
	opts := Options{
		Content:  "test",
		Size:     256,
		Colors:   defaultColors,
		LogoPath: "/nonexistent/logo.png",
	}

	err := GeneratePNGWithLogo(opts, out)
	if err == nil {
		t.Error("expected error for missing logo file")
	}
}

func TestOverlayLogoCentersCorrectly(t *testing.T) {
	qrImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	logo := image.NewRGBA(image.Rect(0, 0, 50, 50))

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			logo.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	result := overlayLogo(qrImg, logo, 100)

	center := result.At(50, 50)
	r, _, _, _ := center.RGBA()
	if r == 0 {
		t.Error("center pixel should have logo color (red), got black")
	}
}

func createTestLogo(t *testing.T, path string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}
