package UTF16

import (
	"utfcoder/logger"
	"utfcoder/types"
	"utfcoder/utils"
)

func isSurrogateCodepoint(high byte, low byte) bool {
	if (high >= 0xD8 && high <= 0xDF) || (low >= 0xD8 && low <= 0xDF) {
		return true
	}
	return false
}

// returns Endianness string "le" or "be", has_BOM boolean
func checkUTF16Endianness(bytes []byte) (types.Endianness, bool) {
	if bytes[0] == 0xFF && bytes[1] == 0xFE {
		logger.Log("UTF-16 Little Endian format detected")
		return types.LITTLE_ENDIAN, true
	} else if bytes[0] == 0xFE && bytes[1] == 0xFF {
		logger.Log("UTF-16 Big Endian format detected")
		return types.BIG_ENDIAN, true
	} else {
		for i := 0; i+1 < len(bytes); i += 2 {
			// if the 4th byte is set or the leading 3 bits of 3rd byte is set, consider as big endian
			if bytes[i] >= 0xD8 && bytes[i] <= 0xDF {
				logger.Log("UTF-16 Big Endian format detected")
				return types.BIG_ENDIAN, false
			} else if bytes[i+1] >= 0xD8 && bytes[i+1] <= 0xDF {
				logger.Log("UTF-16 Little Endian format detected")
				return types.LITTLE_ENDIAN, false
			}
		}
	}

	logger.Log("No UTF-16 Byte Order Mark (BOM) detected. Considering Big Endian format as default")
	return types.BIG_ENDIAN, false
}

func extractBits(hByte byte, lByte byte) uint32 {
	return uint32(hByte)<<8 | uint32(lByte)
}

func extractBitsFromSurrogate(hSurrF, hSurrL, lSurrF, lSurrL byte) uint32 {
	var bits uint32

	leadingTenBits := (uint16(hSurrF)<<8 | uint16(hSurrL)) - 0xD800
	trailingTenBits := (uint16(lSurrF)<<8 | uint16(lSurrL)) - 0xDC00

	bits = uint32(leadingTenBits)<<10 | uint32(trailingTenBits&0x03ff)
	bits = bits + 0x10000

	return bits
}

func ConvertToUTF8(input []byte, addBOM bool) ([]byte, error) {
	logger.Log("\nConvert UTF-16", input, "To UTF-8")

	var output = make([]byte, 0, len(input))

	endianness, hasBOM := checkUTF16Endianness(input)
	isSourceBigEndian := endianness == types.BIG_ENDIAN

	startIdx := 0
	if hasBOM {
		startIdx = 2
	}

	if addBOM {
		// byte order mark for utf8 - 0xEFBBBF
		output = append(output, 0xEF, 0xBB, 0xBF)
	}

	for i := startIdx; i < len(input)-1; i += 2 {
		var bits uint32

		isSurrogate := isSurrogateCodepoint(input[i], input[i+1])

		if isSurrogate && i+3 < len(input) {
			if isSourceBigEndian {
				bits = extractBitsFromSurrogate(input[i], input[i+1], input[i+2], input[i+3])
			} else {
				bits = extractBitsFromSurrogate(input[i+1], input[i], input[i+3], input[i+2])
			}
			i += 2
		} else {
			if isSourceBigEndian {
				bits = extractBits(input[i], input[i+1])
			} else {
				bits = extractBits(input[i+1], input[i])
			}
		}

		if !utils.IsValidUnicodeRange(bits) {
			bits = utils.GenerateUnknownCharacter(types.UTF_8)
		} else if bits >= 0x10000 {
			// Mark with prefix 1111 0xxx 10xx xxxx 10xx xxxx 10xx xxxx and fill the x's with the available bits
			bits = (((bits & 0x1c0000) << 6) | ((bits & 0x30000) << 4)) | ((bits & 0xf000) << 4) | ((bits & 0xfc0) << 2) | (bits & 0x3f) | 0xf0808080
		} else if bits >= 0x800 {
			// Mark with prefix 1110 xxxx 10xx xxxx 10xx xxxx and fill the x's with the available bits
			bits = ((bits & 0xf000) << 4) | ((bits & 0xfc0) << 2) | (bits & 0x3f) | 0xe08080
		} else if bits >= 0x80 {
			// Mark with prefix 110x xxxx 10xx xxxx and fill the x's with the available bits
			bits = ((bits & 0x7c0) << 2) | (bits & 0x3f) | 0xc080
		}

		for range 4 {
			b := byte(bits >> 24)
			if b != 0 {
				output = append(output, b)
			}
			bits = bits << 8
		}
	}

	logger.Log("\nConverted to UTF-8", output)

	return output, nil
}

func ConvertToUTF32(input []byte, targetEncoding string, addBOM bool) ([]byte, error) {
	logger.Log("\nConvert UTF-16", input, "To UTF-32")

	var output = make([]byte, 0, len(input))

	endianness, hasBOM := checkUTF16Endianness(input)

	isTargetBigEndian := targetEncoding == types.UTF_32BE || targetEncoding == types.UTF_32
	isSourceBigEndian := endianness == types.BIG_ENDIAN

	startIdx := 0
	if hasBOM {
		startIdx = 2
	}

	if addBOM {
		if isTargetBigEndian {
			// byte order mark for utf32 BE - 0x0000FEFF
			output = append(output, 0, 0, 0xFE, 0xFF)
		} else {
			// byte order mark for utf16 LE - 0xFFFE0000
			output = append(output, 0xFF, 0xFE, 0, 0)
		}
	}

	for i := startIdx; i < len(input)-1; i += 2 {
		var bits uint32

		isSurrogate := isSurrogateCodepoint(input[i], input[i+1])

		if isSurrogate && i+3 < len(input) {
			if isSourceBigEndian {
				bits = extractBitsFromSurrogate(input[i], input[i+1], input[i+2], input[i+3])
			} else {
				bits = extractBitsFromSurrogate(input[i+1], input[i], input[i+3], input[i+2])
			}
			i += 2
		} else {
			if isSourceBigEndian {
				bits = extractBits(input[i], input[i+1])
			} else {
				bits = extractBits(input[i+1], input[i])
			}
		}

		if !utils.IsValidUnicodeRange(bits) {
			bits = utils.GenerateUnknownCharacter(targetEncoding)
		}

		for range 4 {
			if isTargetBigEndian {
				b := byte(bits >> 24)
				output = append(output, b)
				bits = bits << 8
			} else {
				b := byte(bits)
				output = append(output, b)
				bits = bits >> 8
			}
		}
	}

	logger.Log("\nConverted to", targetEncoding, output)

	return output, nil
}
