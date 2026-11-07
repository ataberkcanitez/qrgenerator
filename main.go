package main

import (
	"fmt"
	"os"

	"ataberkcanitez/qrgenerator/internal/prompt"
	"ataberkcanitez/qrgenerator/internal/qr"
)

func main() {
	content := prompt.Text("Enter URL or text to encode", "")
	format := prompt.Select("Output format", []string{"PNG", "SVG"})

	size := 256
	if format == "PNG" {
		size = prompt.PickSize()
	}

	mode := prompt.Select("QR code style", []string{"Normal", "Custom (with logo image)"})

	if mode == "Custom (with logo image)" && format == "SVG" {
		fmt.Println("Logo overlay is only supported for PNG. Switching to PNG output.")
		format = "PNG"
		size = prompt.PickSize()
	}

	var logoPath string
	if mode == "Custom (with logo image)" {
		logoPath = prompt.Text("Path to logo image (PNG or JPEG)", "")
		if _, err := os.Stat(logoPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: file %q does not exist\n", logoPath)
			os.Exit(1)
		}
	}

	colors := prompt.PickColors()

	defaultName := "qrcode.png"
	if format == "SVG" {
		defaultName = "qrcode.svg"
	}
	outputFile := prompt.Text("Output filename", defaultName)

	opts := qr.Options{
		Content:  content,
		Size:     size,
		Colors:   colors,
		LogoPath: logoPath,
	}

	var err error
	switch {
	case format == "SVG":
		err = qr.GenerateSVG(opts, outputFile)
	case logoPath != "":
		err = qr.GeneratePNGWithLogo(opts, outputFile)
	default:
		err = qr.GeneratePNG(opts, outputFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("QR code saved to %s\n", outputFile)
}
