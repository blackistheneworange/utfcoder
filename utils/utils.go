package utils

import "utfcoder/types"

func GenerateUnknownCharacter(charset string) uint32 {
	// replacement character (U+fffd) is used for representing uknown character
	switch charset {
	case types.UTF_8:
		return 0xefbfbd
	case types.UTF_16, types.UTF_16BE, types.UTF_16LE:
		return 0xfffd
	case types.UTF_32, types.UTF_32BE, types.UTF_32LE:
		return 0x0000fffd
	}

	return 0xfffd
}

func IsValidUnicodeRange(bits uint32) bool {
	// check if the unicode goes beyong U+10FFFF
	if bits > 0x10FFFF {
		return false
	}

	// // check if unicode is within utf-16 surrogate range of U+D800 to U+DFFF
	if bits >= 0xD800 && bits <= 0xDFFF {
		return false
	}

	return true
}
