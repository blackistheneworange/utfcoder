package UTF32

import (
	"errors"
	"flag"
	"utfcoder/logger"
	"utfcoder/types"
)

var addBOM = flag.Bool("bom", false, "specifies whether to include or not include BOM prefix")

func generateUnknownCharacter(charset string) uint32 {
	// replacement character (U+fffd) is used for representing uknown character
	switch charset {
	case types.UTF_8:
		return 0xefbfbd
	case types.UTF_16, types.UTF_16BE, types.UTF_16LE:
		return 0xfffd
	}

	return 0xfffd
}

// returns Endianness string "le" or "be", has_BOM boolean
func checkUTF32Endianness(bytes []byte) (types.Endianness, bool) {
	if bytes[0] == 255 && bytes[1] == 254 && bytes[2] == 0 && bytes[3] == 0 {
		logger.Log("UTF-32 Little Endian format detected")
		return types.LITTLE_ENDIAN, true
	} else if bytes[0] == 0 && bytes[1] == 0 && bytes[2] == 254 && bytes[3] == 255 {
		logger.Log("UTF-32 Big Endian format detected")
		return types.BIG_ENDIAN, true
	} else {
		for i := 0; i+3 < len(bytes); i += 4 {
			// if the 4th byte is set or the leading 3 bits of 3rd byte is set, consider as big endian
			if bytes[i+3] > 0 || bytes[i+2]&0xe0 != 0 {
				logger.Log("UTF-32 Big Endian format detected")
				return types.BIG_ENDIAN, false
			} else if bytes[i] > 0 || bytes[i+1]&0xe0 != 0 {
				logger.Log("UTF-32 Little Endian format detected")
				return types.LITTLE_ENDIAN, false
			}
		}
	}

	logger.Log("No UTF-32 Byte Order Mark (BOM) detected. Considering Little Endian format as default")
	return types.LITTLE_ENDIAN, false
}

func isValidUnicodeRange(bits uint32) bool {
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

func isValidInput(input []byte) bool {
	return len(input) != 0 && len(input)%4 == 0
}

func ConvertToUTF8(input []byte) ([]byte, error) {
	logger.Log("\nConvert UTF-32", input, "To UTF-8")

	if !isValidInput(input) {
		return []byte{}, errors.New("invalid input")
	}

	endianness, hasBOM := checkUTF32Endianness(input)

	var output = make([]byte, 0, len(input))

	startIdx := 0
	if hasBOM {
		startIdx = 4
	}

	if *addBOM {
		// byte order mark for utf8 - 0xEFBBBF
		output = append(output, 0xEF, 0xBB, 0xBF)
	}

	for i := startIdx; i+3 < len(input); i += 4 {
		var bits uint32

		if endianness == types.LITTLE_ENDIAN {
			bits = uint32(input[i+3])<<24 | uint32(input[i+2])<<16 | uint32(input[i+1])<<8 | uint32(input[i])
		} else {
			bits = uint32(input[i])<<24 | uint32(input[i+1])<<16 | uint32(input[i+2])<<8 | uint32(input[i+3])
		}

		if !isValidUnicodeRange(bits) {
			bits = generateUnknownCharacter(types.UTF_8)
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

		for bits != 0 {
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

func ConvertToUTF16(input []byte, targetEncoding string) ([]byte, error) {
	logger.Log("\nConvert UTF-32", input, "To UTF-16")

	if !isValidInput(input) {
		return []byte{}, errors.New("invalid input")
	}

	endianness, hasBOM := checkUTF32Endianness(input)

	var output = make([]byte, 0, len(input))

	startIdx := 0
	if hasBOM {
		startIdx = 4
	}

	isTargetBigEndian := targetEncoding == types.UTF_16BE || targetEncoding == types.UTF_16

	if *addBOM {
		if isTargetBigEndian {
			// byte order mark for utf16 BE - 0xFEFF
			output = append(output, 0xFE, 0xFF)
		} else {
			// byte order mark for utf16 LE - 0xFFFE
			output = append(output, 0xFF, 0xFE)
		}
	}

	for i := startIdx; i+3 < len(input); i += 4 {
		var bits uint32

		if endianness == types.LITTLE_ENDIAN {
			bits = uint32(input[i+3])<<24 | uint32(input[i+2])<<16 | uint32(input[i+1])<<8 | uint32(input[i])
		} else {
			bits = uint32(input[i])<<24 | uint32(input[i+1])<<16 | uint32(input[i+2])<<8 | uint32(input[i+3])
		}

		var highSurrogate, lowSurrogate uint16

		if !isValidUnicodeRange(bits) {
			bits = generateUnknownCharacter(targetEncoding)
		} else if bits >= 0x10000 {
			bits = bits - 0x10000

			// high surrogate - add 0xD800 with the leading 10 bits
			highSurrogate = 0xD800 + uint16(bits>>10)
			// low surrogate - add 0xDC00 with the trailing 10 bits
			lowSurrogate = 0xDC00 + uint16(bits&0x03ff)
		}

		if isTargetBigEndian {
			if lowSurrogate != 0 {
				output = append(output, byte(highSurrogate>>8), byte(highSurrogate), byte(lowSurrogate>>8), byte(lowSurrogate))
			} else {
				output = append(output, byte(bits>>8), byte(bits))
			}
		} else {
			if lowSurrogate != 0 {
				output = append(output, byte(highSurrogate), byte(highSurrogate>>8), byte(lowSurrogate), byte(lowSurrogate>>8))
			} else {
				output = append(output, byte(bits), byte(bits>>8))
			}
		}
	}

	logger.Log("\nConverted to", targetEncoding, output)

	return output, nil
}
