package utils_test

import (
	"github.com/bigcubecat/netsnake/internal/utils"
	"testing"
)

func TestIntToHexColor(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Black", 0x000000, "#000000"},
		{"White", 0xFFFFFF, "#FFFFFF"},
		{"Red", 0xFF0000, "#FF0000"},
		{"Green", 0x00FF00, "#00FF00"},
		{"Blue", 0x0000FF, "#0000FF"},
		{"Cyan", 0x00FFFF, "#00FFFF"},
		{"Magenta", 0xFF00FF, "#FF00FF"},
		{"Yellow", 0xFFFF00, "#FFFF00"},
		{"Out of range (negative)", -1, "#000000"},
		{"Out of range (too large)", 0x1000000, "#FFFFFF"},
		{"Random color", 0x1A2B3C, "#1A2B3C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.IntToHexColor(tt.input)
			if result != tt.expected {
				t.Errorf("intToHexColor(%d) = %s; expected %s",
					tt.input,
					result,
					tt.expected,
				)
			}
		})
	}
}

func TestHexColorToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{"Black", "#000000", 0x000000, false},
		{"White", "#FFFFFF", 0xFFFFFF, false},
		{"Red", "#FF0000", 0xFF0000, false},
		{"Green", "#00FF00", 0x00FF00, false},
		{"Blue", "#0000FF", 0x0000FF, false},
		{"Cyan", "#00FFFF", 0x00FFFF, false},
		{"Magenta", "#FF00FF", 0xFF00FF, false},
		{"Yellow", "#FFFF00", 0xFFFF00, false},
		{"Random color", "#1A2B3C", 0x1A2B3C, false},
		{"Invalid length (short)", "#12345", 0, true},
		{"Invalid length (long)", "#1234567", 0, true},
		{"Invalid characters", "#GHIJKL", 0, true},
		{"No # prefix", "FFFFFF", 0xFFFFFF, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.HexColorToInt(tt.input)

			if tt.hasError {
				// Expecting an error
				if err == nil {
					t.Errorf("hexColorToInt(%s) expected an error, but got none", tt.input)
				}
			} else {
				// Not expecting an error
				if err != nil {
					t.Errorf("hexColorToInt(%s) returned an error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("hexColorToInt(%s) = %d; expected %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}
