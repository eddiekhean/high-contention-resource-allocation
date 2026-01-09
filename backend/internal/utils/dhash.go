package utils

func ParseDHash(binary string) (uint64, error) {
	var hash uint64
	for _, c := range binary {
		hash <<= 1
		if c == '1' {
			hash |= 1
		}
	}
	return hash, nil
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
