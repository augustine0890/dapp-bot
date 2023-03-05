package rankcard

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// parseColor parses a color from a string.
func parseColor(c interface{}) (color.Color, error) {
	switch c := c.(type) {
	case string:
		if c[0] == '#' {
			return parseHexColor(c)
		}
		return parseColor(c)
	case color.Color:
		return c, nil
	default:
		return nil, fmt.Errorf("invalid color type %T", c)
	}
}

// parseGradientColors parses gradient colors from an interface{}.
func parseGradientColors(colors interface{}) ([]color.Color, error) {
	switch colors := colors.(type) {
	case []string:
		result := make([]color.Color, len(colors))
		for i, c := range colors {
			if c[0] == '#' {
				col, err := parseHexColor(c)
				if err != nil {
					return nil, err
				}
				result[i] = col
			} else {
				col, err := parseColor(c)
				if err != nil {
					return nil, err
				}
				result[i] = col
			}
		}
		return result, nil
	case []color.Color:
		return colors, nil
	default:
		return nil, fmt.Errorf("invalid gradient color type %T", colors)
	}
}

// This function takes a hex color string as input (e.g. "#FF0000" for red) and
// returns a color.Color object that can be used in the gg package to set colors.
func parseHexColor(hex string) (color.Color, error) {
	if !strings.HasPrefix(hex, "#") {
		return nil, fmt.Errorf("invalid hex color format: %s", hex)
	}
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = hex + hex
	}
	if len(hex) != 6 {
		return nil, fmt.Errorf("invalid hex color format: %s", hex)
	}
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid hex color format: %s", hex)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid hex color format: %s", hex)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid hex color format: %s", hex)
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF}, nil
}

// Converts numbers into human-readable units like "1K", "1M", "1B", etc.
// The returned string is rounded to 1 decimal place and includes the appropriate unit.
func convertNumberToUnits(number int) string {
	units := []string{"", "K", "M", "B", "T"}
	i := 0
	for number >= 1000 {
		number /= 1000
		i++
	}
	return fmt.Sprintf("%d%s", number, units[i])
}

// ShortenText shortens the given text to the given maximum length and adds ellipsis if necessary.
func ShortenText(text string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	if maxLength >= len(text) {
		return text
	}
	return text[:maxLength] + "..."
}
