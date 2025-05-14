package xqb

import (
	"encoding/json"
	"fmt"
	"os"
)

type Color string

// ANSI Color Codes
const (
	Reset    Color = "\033[0m"
	Blue     Color = "\033[1;34m"
	Green    Color = "\033[1;32m"
	Yellow   Color = "\033[1;33m"
	Red      Color = "\033[1;31m"
	BG_Black Color = "\033[40m"
)

func (c Color) Value() string {
	return string(c)
}

func DD(v ...any) {
	Dump(v...)
	Die()
}

func Dump(v ...any) {
	for _, val := range v {
		fmt.Println(formatOutput(val))
	}
}

func Die() {
	os.Exit(1)
}

func formatOutput(value any) string {
	red := Green.Value()
	reset := Reset.Value()
	blackBG := BG_Black.Value()

	// Handle different types
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s%s%s%s", blackBG, red, v, reset)
	case error:
		return fmt.Sprintf("%s%s%s%s", blackBG, red, v.Error(), reset)
	default:
		// Use JSON for complex structures
		jsonData, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%s%s[Error Formatting]%s", blackBG, red, reset)
		}
		return fmt.Sprintf("%s%s%s%s", blackBG, red, jsonData, reset)
	}
}
