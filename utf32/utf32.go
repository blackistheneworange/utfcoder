package UTF32

import (
	"errors"
	"flag"
	"utfcoder/logger"
	"utfcoder/types"
)

var addBOM = flag.Bool("addbom", false, "specifies whether to include or not include BOM prefix")

func generateGarbageCharacter() [4]byte {
	return [4]byte{0xbd, 0xbf, 0xef, 0}
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

func isValidUnicodeRange(bytes [4]byte) bool {
	// check if the unicode goes beyong U+10FFFF
	if bytes[3]&0xff != 0 || bytes[2] > 0x10 {
		return false
	}

	// check if unicode is within utf-16 surrogate range of U+D800 to U+DFFF
	if bytes[1] >= 216 && bytes[1] <= 223 && bytes[2] == 0 && bytes[3] == 0 {
		return false
	}

	return true
}

func getWordLength(bytes [4]byte) uint8 {
	if bytes[3] != 0 || bytes[2] != 0 { // if 3rd or 4th byte is not 0, then it is going to take 4 bytes to represent in utf8
		return 4
	} else if bytes[1]&0xf8 != 0 { // if any of 2nd byte's leading 5 bits is set, then it is going to take 3 bytes to represent in utf8
		return 3
	} else if bytes[1] != 0 || bytes[0]&0x80 != 0 { // else if 2nd byte is not zero or 1st byte's leading bit is set, then it is going to take 2 bytes to represent in utf8
		return 2
	}

	return 1
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

	for i := startIdx; i+3 < len(input); i += 4 {
		var bytes = [4]byte{}

		if endianness == types.LITTLE_ENDIAN {
			bytes[0], bytes[1], bytes[2], bytes[3] = input[i], input[i+1], input[i+2], input[i+3]
		} else {
			bytes[0], bytes[1], bytes[2], bytes[3] = input[i+3], input[i+2], input[i+1], input[i]
		}

		wordLength := getWordLength(bytes)

		if !isValidUnicodeRange(bytes) {
			bytes = generateGarbageCharacter()
		} else if wordLength == 4 {
			bytes[3] = 0xf0 | ((bytes[2] & 0x1c) >> 2)           // mark 4th byte as 1111 0xxx to indicate 4 width byte character, and copy 3rd, 4th, 5th leading bits from 3rd byte
			bytes[2] = ((bytes[2]<<4)|bytes[1]>>4)&0x3f | 0x80   // copy leading 4 bits from 2nd byte, and set 10 prefix
			bytes[1] = ((bytes[1]<<2)|(bytes[0]>>6))&0x3f | 0x80 // copy leading 2 bits from 1st byte, and set 10 prefix
			bytes[0] = (bytes[0] | 0x80) & 0xbf                  // set 10 prefix
		} else if wordLength == 3 {
			bytes[2] = 0xe0 | (bytes[1] >> 4)                    // mark 3rd byte as 1110 xxxx to indicate 3 width byte character, and copy leading 4 bits from 2nd byte
			bytes[1] = ((bytes[1]<<2)|(bytes[0]>>6))&0x3f | 0x80 // copy leading 2 bits from 1st byte, and set 10 prefix
			bytes[0] = (bytes[0] & 0x3f) | 0x80                  // set 10 prefix
		} else if wordLength == 2 {
			bytes[1] = ((bytes[1]<<2)&0x1f | 0xc0) | (bytes[0] >> 6) // mark 2nd byte as 110x xxxx to indicate 2 width byte character, and copy leading 2 bits from 1st byte
			bytes[0] = (bytes[0] & 0x3f) | 0x80                      // set 10 prefix
		}

		for i := len(bytes) - 1; i >= 0; i-- {
			if bytes[i] != 0 {
				output = append(output, bytes[i])
			}
		}
	}

	if *addBOM {
		bomPrefixedOutput := append([]byte{0xEF, 0xBB, 0xBF}, output...)
		output = bomPrefixedOutput
	}

	logger.Log("\nConverted to UTF-8", output)

	return output, nil
}
