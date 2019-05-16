package ec

const (
	GF_ZERO uint8 = 0x00
	GF_ONE  uint8 = 0x01

	MAX_K = 16
	MAX_M = 8

	// generator polynomial: x^8 + x^4 + x^3 + x^2 + 1
	POLY    uint8 = 0x1d
)