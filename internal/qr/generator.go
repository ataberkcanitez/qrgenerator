package qr

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
	qrcode "github.com/skip2/go-qrcode"
)

type Options struct {
	Content  string
	Size     int
	Colors   ColorPair
	LogoPath string
}

func GeneratePNG(opts Options, outputFile string) error {
	qr, err := qrcode.New(opts.Content, qrcode.Highest)
	if err != nil {
		return err
	}
	qr.ForegroundColor = opts.Colors.Foreground
	qr.BackgroundColor = opts.Colors.Background
	return qr.WriteFile(opts.Size, outputFile)
}

func GenerateSVG(opts Options, outputFile string) error {
	bitmap, err := buildBitmap(opts.Content)
	if err != nil {
		return err
	}
	svg := renderSVG(bitmap, opts.Colors.Foreground, opts.Colors.Background)
	return os.WriteFile(outputFile, []byte(svg), 0644)
}

func GeneratePNGWithLogo(opts Options, outputFile string) error {
	qr, err := qrcode.New(opts.Content, qrcode.Highest)
	if err != nil {
		return err
	}
	qr.ForegroundColor = opts.Colors.Foreground
	qr.BackgroundColor = opts.Colors.Background
	qrImg := qr.Image(opts.Size)

	logoImg, err := loadImage(opts.LogoPath)
	if err != nil {
		return err
	}

	result := overlayLogo(qrImg, logoImg, opts.Size)

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, result)
}

func buildBitmap(content string) ([][]bool, error) {
	qr, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		return nil, err
	}
	return qr.Bitmap(), nil
}

func renderSVG(bitmap [][]bool, fg, bg color.RGBA) string {
	modules := len(bitmap)
	cellSize := 10
	totalSize := modules * cellSize

	fgStr := FormatHexColor(fg)
	bgStr := FormatHexColor(bg)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d">`, totalSize, totalSize, totalSize, totalSize))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="%s"/>`, totalSize, totalSize, bgStr))
	sb.WriteString("\n")

	for y, row := range bitmap {
		for x, cell := range row {
			if cell {
				sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x*cellSize, y*cellSize, cellSize, cellSize, fgStr))
				sb.WriteString("\n")
			}
		}
	}

	sb.WriteString("</svg>\n")
	return sb.String()
}

func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}
	return img, nil
}

func overlayLogo(qrImg image.Image, logo image.Image, size int) *image.RGBA {
	logoSize := uint(size / 5)
	resizedLogo := resize.Resize(logoSize, logoSize, logo, resize.Lanczos3)

	canvas := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(canvas, canvas.Bounds(), qrImg, image.Point{}, draw.Src)

	offset := (size - int(logoSize)) / 2
	logoRect := image.Rect(offset, offset, offset+int(logoSize), offset+int(logoSize))
	draw.Draw(canvas, logoRect, resizedLogo, image.Point{}, draw.Over)

	return canvas
}
