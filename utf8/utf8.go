package UTF8

import (
	"utfcoder/logger"
	"utfcoder/types"
	"utfcoder/utils"
)

func ConvertToUTF32(input []byte, targetEncoding string, addBOM bool) ([]byte, error) {
	logger.Log("\nConvert UTF-32", input, "To UTF-8")

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
