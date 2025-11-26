package UTF8

import (
	"utfcoder/logger"
	"utfcoder/types"
	"utfcoder/utils"
)

func ConvertToUTF32(input []byte, targetEncoding string, addBOM bool) ([]byte, error) {
	logger.Log("\nConvert UTF-8", input, "To UTF-32")

	var output = make([]byte, 0, len(input))

	isTargetBigEndian := targetEncoding == types.UTF_32BE || targetEncoding == types.UTF_32

	startIdx := 0
	// check if the first 3 bytes represent byte order mark for utf-8 i.e. 0xEFBBBF
	if len(input) > 2 && input[0] == 0xEF && input[1] == 0xBB && input[2] == 0xBF {
		startIdx = 3
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

	for i := startIdx; i < len(input); i += 1 {
		var bits uint32

		if input[i]&0xf8 == 240 && i+3 < len(input) {
			bits = uint32(input[i]&0x07)<<18 | uint32(input[i+1]&0x3f)<<12 | uint32(input[i+2]&0x3f)<<6 | uint32(input[i+3]&0x3f)
			i += 3
		} else if input[i]&0xf0 == 224 && i+2 < len(input) {
			bits = uint32(input[i]&0x0f)<<12 | uint32(input[i+1]&0x3f)<<6 | uint32(input[i+2]&0x3f)
			i += 2
		} else if input[i]&0xe0 == 192 && i+1 < len(input) {
			bits = uint32(input[i]&0x1f)<<6 | uint32(input[i+1]&0x3f)
			i += 1
		} else {
			bits = uint32(input[i])
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

func ConvertToUTF16(input []byte, targetEncoding string, addBOM bool) ([]byte, error) {
	logger.Log("\nConvert UTF-8", input, "To", targetEncoding)

	var output = make([]byte, 0, len(input))

	isTargetBigEndian := targetEncoding == types.UTF_16BE || targetEncoding == types.UTF_16

	startIdx := 0
	// check if the first 3 bytes represent byte order mark for utf-8 i.e. 0xEFBBBF
	if len(input) > 2 && input[0] == 0xEF && input[1] == 0xBB && input[2] == 0xBF {
		startIdx = 3
	}

	if addBOM {
		if isTargetBigEndian {
			// byte order mark for utf16 BE - 0xFEFF
			output = append(output, 0xFE, 0xFF)
		} else {
			// byte order mark for utf16 LE - 0xFFFE
			output = append(output, 0xFF, 0xFE)
		}
	}

	for i := startIdx; i < len(input); i += 1 {
		var bits uint32
		var highSurrogate, lowSurrogate uint16

		if input[i]&0xf8 == 240 && i+3 < len(input) {
			bits = uint32(input[i]&0x07)<<18 | uint32(input[i+1]&0x3f)<<12 | uint32(input[i+2]&0x3f)<<6 | uint32(input[i+3]&0x3f)
			i += 3
		} else if input[i]&0xf0 == 224 && i+2 < len(input) {
			bits = uint32(input[i]&0x0f)<<12 | uint32(input[i+1]&0x3f)<<6 | uint32(input[i+2]&0x3f)
			i += 2
		} else if input[i]&0xe0 == 192 && i+1 < len(input) {
			bits = uint32(input[i]&0x1f)<<6 | uint32(input[i+1]&0x3f)
			i += 1
		} else {
			bits = uint32(input[i])
		}

		if !utils.IsValidUnicodeRange(bits) {
			bits = utils.GenerateUnknownCharacter(targetEncoding)
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
