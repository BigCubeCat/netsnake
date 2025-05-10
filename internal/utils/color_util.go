package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// intToHexColor converts an integer to a hex color code.
func IntToHexColor(color int) string {
	// Ensure the color is within the valid range (0x000000 to 0xFFFFFF)
	if color < 0 {
		color = 0
	} else if color > 0xFFFFFF {
		color = 0xFFFFFF
	}

	// Format the integer as a 6-digit hex string, prefixed with "#"
	return fmt.Sprintf("#%06X", color)
}

func HexColorToInt(hex string) (int, error) {
	// Remove the "#" prefix if present
	hex = strings.TrimPrefix(hex, "#")

	// Ensure the hex string is exactly 6 characters long
	if len(hex) != 6 {
		return 0, fmt.Errorf("invalid hex color length: must be 6 characters")
	}

	// Parse the hex string to an integer
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hex color: %v", err)
	}

	return int(value), nil
}
