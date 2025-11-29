package UTF16

import (
	"bytes"
	"testing"
	"utfcoder/types"
)

func TestConvertToUTF32LE(t *testing.T) {
	for idx := 0; idx < len(utf16LittleEndianTo32LittleEndianTestInputs); idx += 2 {
		input := utf16LittleEndianTo32LittleEndianTestInputs[idx]
		expected := utf16LittleEndianTo32LittleEndianTestInputs[idx+1]
		output, err := ConvertToUTF32(append([]byte{0xFF, 0xFE}, input...), types.UTF_32LE, false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestConvertToUTF32LEWithBOM(t *testing.T) {
	input := utf16LittleEndianTo32LittleEndianTestInputs[0]
	expected := append([]byte{0xFF, 0xFE, 0, 0}, utf16LittleEndianTo32LittleEndianTestInputs[1]...)
	output, err := ConvertToUTF32(append([]byte{0xFF, 0xFE}, input...), types.UTF_32LE, true)

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
	}
}

func TestConvertToUTF32BE(t *testing.T) {
	for idx := 0; idx < len(utf16LittleEndianTo32BigEndianTestInputs); idx += 2 {
		input := utf16LittleEndianTo32BigEndianTestInputs[idx]
		expected := utf16LittleEndianTo32BigEndianTestInputs[idx+1]
		output, err := ConvertToUTF32(append([]byte{0xFF, 0xFE}, input...), types.UTF_32BE, false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestConvertToUTF32BEWithBOM(t *testing.T) {
	input := utf16LittleEndianTo32BigEndianTestInputs[0]
	expected := append([]byte{0, 0, 0xFE, 0xFF}, utf16LittleEndianTo32BigEndianTestInputs[1]...)
	output, err := ConvertToUTF32(append([]byte{0xFF, 0xFE}, input...), types.UTF_32BE, true)

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
	}
}

func TestConvertToUTF8(t *testing.T) {
	for idx := 0; idx < len(utf16LittleEndianTo8TestInputs); idx += 2 {
		input := utf16LittleEndianTo8TestInputs[idx]
		expected := utf16LittleEndianTo8TestInputs[idx+1]
		output, err := ConvertToUTF8(append([]byte{0xFF, 0xFE}, input...), false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF8(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestConvertToUTF8WithBOM(t *testing.T) {
	input := utf16LittleEndianTo8TestInputs[0]
	expected := append([]byte{0xEF, 0xBB, 0xBF}, utf16LittleEndianTo8TestInputs[1]...)
	output, err := ConvertToUTF8(append([]byte{0xFF, 0xFE}, input...), true)

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf(`ConvertToUTF8(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
	}
}

var utf16LittleEndianTo32LittleEndianTestInputs = [][]byte{
	{65, 0}, {65, 0, 0, 0},
	{122, 0}, {122, 0, 0, 0},
	{32, 0}, {32, 0, 0, 0},

	{233, 0}, {233, 0, 0, 0},
	{169, 0}, {169, 0, 0, 0},
	{176, 0}, {176, 0, 0, 0},
	{241, 0}, {241, 0, 0, 0},

	{16, 4}, {16, 4, 0, 0},
	{65, 4}, {65, 4, 0, 0},

	{255, 7}, {255, 7, 0, 0},
	{255, 215}, {255, 215, 0, 0},

	{128, 8}, {128, 8, 0, 0},
	{172, 32}, {172, 32, 0, 0},
	{185, 32}, {185, 32, 0, 0},
	{45, 48}, {45, 48, 0, 0},

	{0, 215}, {0, 215, 0, 0},

	{0, 216, 0, 220}, {0, 0, 1, 0}, // U+10000
	{0, 216, 1, 220}, {1, 0, 1, 0}, // U+10001

	{61, 216, 0, 222}, {0, 246, 1, 0}, // 😀 U+1F600

	{255, 219, 255, 223}, {255, 255, 16, 0}, // U+10FFFF

	{196, 219, 0, 220}, {0, 16, 16, 0}, // U+101000

	{255, 16}, {255, 16, 0, 0}, // U+10FF

	{255, 0}, {255, 0, 0, 0}, // U+00FF

	{253, 255}, {253, 255, 0, 0}, // U+FFFD
	{254, 255}, {254, 255, 0, 0}, // U+FFFE

	// invalid surrogate cases → U+FFFD
	{60, 216}, {253, 255, 0, 0},
	{0, 216}, {253, 255, 0, 0},
	{17, 0}, {17, 0, 0, 0}, // U+0011 control character

	// UTF-16 BOM cases
	{255, 254}, {255, 254, 0, 0}, // U+FEFF
	{255, 253}, {255, 253, 0, 0}, // U+FDFF
}

var utf16LittleEndianTo32BigEndianTestInputs = [][]byte{
	{65, 0}, {0, 0, 0, 65},
	{122, 0}, {0, 0, 0, 122},
	{32, 0}, {0, 0, 0, 32},

	{233, 0}, {0, 0, 0, 233},
	{169, 0}, {0, 0, 0, 169},
	{176, 0}, {0, 0, 0, 176},
	{241, 0}, {0, 0, 0, 241},

	{16, 4}, {0, 0, 4, 16}, // U+0410
	{65, 4}, {0, 0, 4, 65}, // U+0441

	{255, 7}, {0, 0, 7, 255}, // U+07FF
	{255, 215}, {0, 0, 215, 255}, // U+D7FF

	{128, 8}, {0, 0, 8, 128}, // U+0880
	{172, 32}, {0, 0, 32, 172}, // €
	{185, 32}, {0, 0, 32, 185}, // ₹
	{45, 48}, {0, 0, 48, 45}, // ㌭

	{0, 215}, {0, 0, 215, 0}, // U+D700

	// Supplementary planes (surrogates)
	{0, 216, 0, 220}, {0, 1, 0, 0}, // U+10000
	{0, 216, 1, 220}, {0, 1, 0, 1}, // U+10001

	{61, 216, 0, 222}, {0, 1, 246, 0}, // 😀 U+1F600

	{255, 219, 255, 223}, {0, 16, 255, 255}, // U+10FFFF

	{196, 219, 0, 220}, {0, 16, 16, 0}, // U+101000

	{255, 16}, {0, 0, 16, 255}, // U+10FF

	{255, 0}, {0, 0, 0, 255}, // U+00FF

	{253, 255}, {0, 0, 255, 253}, // U+FFFD
	{254, 255}, {0, 0, 255, 254}, // U+FFFE

	// Invalid surrogate cases → U+FFFD
	{60, 216}, {0, 0, 255, 253},
	{0, 216}, {0, 0, 255, 253},

	// BMP control character
	{17, 0}, {0, 0, 0, 17}, // U+0011

	// BOM-as-character cases
	{255, 254}, {0, 0, 254, 255}, // U+FEFF
	{255, 253}, {0, 0, 253, 255}, // U+FDFF
}

var utf16LittleEndianTo8TestInputs = [][]byte{
	{65, 0}, {65}, // 'A'
	{122, 0}, {122}, // 'z'
	{32, 0}, {32}, // space

	{233, 0}, {195, 169}, // é
	{169, 0}, {194, 169}, // ©
	{176, 0}, {194, 176}, // °
	{241, 0}, {195, 177}, // ñ

	{16, 4}, {208, 144}, // А (U+0410)
	{65, 4}, {209, 129}, // с (U+0441)

	{255, 7}, {223, 191}, // ߿ (U+07FF)
	{255, 215}, {237, 159, 191}, // ퟿ (U+D7FF)

	{128, 8}, {224, 162, 128}, // U+0880
	{172, 32}, {226, 130, 172}, // €
	{185, 32}, {226, 130, 185}, // ₹
	{45, 48}, {227, 128, 173}, // ㌭

	{0, 215}, {237, 156, 128}, // U+D700

	//-----------------------------------------------------
	// Supplementary Planes (Surrogate Pairs)
	//-----------------------------------------------------

	{0, 216, 0, 220}, {240, 144, 128, 128}, // U+10000
	{0, 216, 1, 220}, {240, 144, 128, 129}, // U+10001

	{61, 216, 0, 222}, {240, 159, 152, 128}, // 😀 (U+1F600)

	{255, 219, 255, 223}, {244, 143, 191, 191}, // U+10FFFF

	{196, 219, 0, 220}, {244, 129, 128, 128}, // U+101000

	//-----------------------------------------------------
	// BMP special values
	//-----------------------------------------------------

	{255, 16}, {225, 131, 191}, // U+10FF
	{255, 0}, {195, 191}, // ÿ (U+00FF)

	//-----------------------------------------------------
	// Noncharacters & replacement
	//-----------------------------------------------------

	{253, 255}, {239, 191, 189}, // U+FFFD
	{254, 255}, {239, 191, 190}, // U+FFFE

	//-----------------------------------------------------
	// Invalid surrogate cases → replaced with U+FFFD
	//-----------------------------------------------------

	{60, 216}, {239, 191, 189}, // invalid high surrogate
	{0, 216}, {239, 191, 189}, // lone low surrogate

	//-----------------------------------------------------
	// Control character
	//-----------------------------------------------------

	{17, 0}, {17}, // U+0011

	//-----------------------------------------------------
	// BOM as character (not encoding marker here)
	//-----------------------------------------------------

	{255, 254}, {239, 187, 191}, // U+FEFF (ZERO WIDTH NO-BREAK SPACE)

	//-----------------------------------------------------
	// Noncharacter valid code point
	//-----------------------------------------------------

	{255, 253}, {239, 183, 191}, // U+FDFF
}
