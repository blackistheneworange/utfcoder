package UTF16

import (
	"utfcoder/logger"
	"utfcoder/types"
	"utfcoder/utils"
)

func ConvertToUTF32(input []byte, targetEncoding string, addBOM bool) ([]byte, error) {
	logger.Log("\nConvert UTF-16", input, "To UTF-32")

	var output = make([]byte, 0, len(input))

	var leadingTenBits uint16
	var trailingTenBits uint16

	isTargetBigEndian := targetEncoding == types.UTF_32BE || targetEncoding == types.UTF_32
	isSourceBigEndian := true

	startIdx := 0
	// check if the first 2 bytes represent byte order mark for utf-16 little endian i.e. 0xFFFE
	if len(input) > 1 {
		if input[0] == 0xFF && input[1] == 0xFE {
			logger.Log("UTF-16 Little Endian format detected")
			startIdx = 2
			isSourceBigEndian = false
		} else if input[0] == 0xFE && input[1] == 0xFF {
			logger.Log("UTF-16 Big Endian format detected")
			startIdx = 2
		} else {
			logger.Log("No UTF-16 Byte Order Mark (BOM) detected. Considering Big Endian format as default")
		}
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

		if isSourceBigEndian {
			if i+3 < len(input) && input[i] >= 0xD8 && input[i] <= 0xDF {
				leadingTenBits = (uint16(input[i])<<8 | uint16(input[i+1])) - 0xD800
				trailingTenBits = (uint16(input[i+2])<<8 | uint16(input[i+3])) - 0xDC00

				bits = uint32(leadingTenBits)<<10 | uint32(trailingTenBits&0x03ff)
				bits = bits + 0x10000
				i += 2
			} else if i+1 < len(input) {
				bits = uint32(input[i])<<8 | uint32(input[i+1])
			}
		} else {
			if i+3 < len(input) && input[i+1] >= 0xD8 && input[i+1] <= 0xDF {
				leadingTenBits = (uint16(input[i+1])<<8 | uint16(input[i])) - 0xD800
				trailingTenBits = (uint16(input[i+3])<<8 | uint16(input[i+2])) - 0xDC00

				bits = uint32(leadingTenBits)<<10 | uint32(trailingTenBits&0x03ff)
				bits = bits + 0x10000
				i += 2
			} else if i+1 < len(input) {
				bits = uint32(input[i+1])<<8 | uint32(input[i])
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
