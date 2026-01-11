package utils

import (
	"strconv"
	"strings"
)

func ParseDHash(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	// Try hex first if it's 16 chars or starts with 0x
	if len(s) == 16 || strings.HasPrefix(s, "0x") {
		cleanHex := strings.TrimPrefix(s, "0x")
		if len(cleanHex) == 16 {
			return strconv.ParseUint(cleanHex, 16, 64)
		}
	}

	// Try binary if it's 64 chars
	if len(s) == 64 {
		return strconv.ParseUint(s, 2, 64)
	}

	// Fallback: try as decimal number
	return strconv.ParseUint(s, 10, 64)
}

func HammingDistance(a, b uint64) int {
	x := a ^ b
	count := 0
	for x > 0 {
		x &= x - 1
		count++
	}
	return count
}
