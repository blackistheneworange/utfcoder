package UTF32

import (
	"errors"
	"flag"
	"utfcoder/logger"
	"utfcoder/types"
)

var addBOM = flag.Bool("bom", false, "specifies whether to include or not include BOM prefix")

func generateUnknownCharacter(charset string) [4]byte {
	// replacement character (U+fffd) is used for representing uknown character

	switch charset {
	case types.UTF_8:
		return [4]byte{0xbd, 0xbf, 0xef, 0}
	case types.UTF_16, types.UTF_16BE:
		return [4]byte{0xff, 0xfd, 0, 0}
	case types.UTF_16LE:
		return [4]byte{0xfd, 0xff, 0, 0}
	}

	return [4]byte{0, 0, 0xff, 0xfd}
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
			bits = 0xefbfbd
		} else if bits >= 0x10000 {
			bits = (((bits & 0x1c0000) << 6) | ((bits & 0x30000) << 4)) | ((bits & 0xf000) << 4) | ((bits & 0xfc0) << 2) | (bits & 0x3f) | 0xf0808080
		} else if bits >= 0x800 {
			bits = ((bits & 0xf000) << 4) | ((bits & 0xfc0) << 2) | (bits & 0x3f) | 0xe08080
		} else if bits >= 0x80 {
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
			output = append(output, 0xFE, 0xFF)
		} else {
			output = append(output, 0xFF, 0xFE)
		}
	}

	for i := startIdx; i+3 < len(input); i += 4 {
		var bytes = [4]byte{}

		if endianness == types.LITTLE_ENDIAN {
			bytes[0], bytes[1], bytes[2], bytes[3] = input[i], input[i+1], input[i+2], input[i+3]
		} else {
			bytes[0], bytes[1], bytes[2], bytes[3] = input[i+3], input[i+2], input[i+1], input[i]
		}

		var bits uint32 = uint32(bytes[3])<<24 | uint32(bytes[2])<<16 | uint32(bytes[1])<<8 | uint32(bytes[0])

		if !isValidUnicodeRange(bits) {
			bytes = generateUnknownCharacter(targetEncoding)
		} else if bits >= 0x10000 {
			bytes[2] = bytes[2] - 0x01
			var surrogateBits uint32 = (uint32(bytes[2]&0x0f) << 16) | uint32(bytes[1])<<8 | uint32(bytes[0])

			// high surrogate
			var highSurrogate uint16 = 0xD800 + uint16(surrogateBits>>10)
			bytes[3] = byte(highSurrogate >> 8)
			bytes[2] = byte(highSurrogate)

			// low surrogate
			var lowSurrogate = 0xDC00 + uint16(surrogateBits&0x03ff)
			bytes[1] = byte(lowSurrogate >> 8)
			bytes[0] = byte(lowSurrogate)
		} else {
			bytes[2] = 0
			bytes[3] = 0
		}

		if isTargetBigEndian {
			if bytes[0] == 0 && bytes[1] == 0 {
				continue
			}
			if bytes[3] != 0 {
				output = append(output, bytes[3], bytes[2])
			}
			output = append(output, bytes[1], bytes[0])
		} else {
			if bytes[0] == 0 && bytes[1] == 0 {
				continue
			}
			if bytes[3] != 0 {
				output = append(output, bytes[2], bytes[3])
			}
			output = append(output, bytes[0], bytes[1])
		}
	}

	logger.Log("\nConverted to", targetEncoding, output)

	return output, nil
}
