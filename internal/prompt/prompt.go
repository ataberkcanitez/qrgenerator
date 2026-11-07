package prompt

import (
	"fmt"
	"os"
	"strings"

	"ataberkcanitez/qrgenerator/internal/qr"

	"github.com/manifoldco/promptui"
)

func Text(label, defaultVal string) string {
	p := promptui.Prompt{
		Label:   label,
		Default: defaultVal,
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("input cannot be empty")
			}
			return nil
		},
	}
	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed: %v\n", err)
		os.Exit(1)
	}
	return result
}

func Select(label string, items []string) string {
	p := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed: %v\n", err)
		os.Exit(1)
	}
	return result
}

func PickColors() qr.ColorPair {
	choice := Select("QR code color", qr.PresetNames)

	if choice != "Custom (enter hex colors)" {
		return qr.Presets[choice]
	}

	fgHex := Text("Foreground color (hex, e.g. #FF5733)", "#000000")
	bgHex := Text("Background color (hex, e.g. #FFFFFF)", "#FFFFFF")

	fg, err := qr.ParseHexColor(fgHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid foreground color: %v\n", err)
		os.Exit(1)
	}
	bg, err := qr.ParseHexColor(bgHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid background color: %v\n", err)
		os.Exit(1)
	}
	return qr.ColorPair{Foreground: fg, Background: bg}
}

func PickSize() int {
	sizeStr := Select("Image size", []string{"256x256", "512x512", "1024x1024"})
	switch sizeStr {
	case "512x512":
		return 512
	case "1024x1024":
		return 1024
	default:
		return 256
	}
}
