package UTF8

import (
	"bytes"
	"testing"
	"utfcoder/types"
)

func TestConvertToUTF32BE(t *testing.T) {
	for idx := 0; idx < len(utf8To32BigEndianTestInputs); idx += 2 {
		input := utf8To32BigEndianTestInputs[idx]
		expected := utf8To32BigEndianTestInputs[idx+1]
		output, err := ConvertToUTF32(input, types.UTF_32BE, false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestConvertToUTF32BEWithBOM(t *testing.T) {
	input := utf8To32BigEndianTestInputs[0]
	expected := []byte{0, 0, 0xfe, 0xff}
	expected = append(expected, utf8To32BigEndianTestInputs[1]...)

	output, err := ConvertToUTF32(input, types.UTF_32BE, true)

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
	}
}

func TestConvertToUTF32LE(t *testing.T) {
	for idx := 0; idx < len(utf8To32LittleEndianTestInputs); idx += 2 {
		input := utf8To32LittleEndianTestInputs[idx]
		expected := utf8To32LittleEndianTestInputs[idx+1]
		output, err := ConvertToUTF32(input, types.UTF_32LE, false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestConvertToUTF32LEWithBOM(t *testing.T) {
	input := utf8To32LittleEndianTestInputs[0]
	expected := []byte{0xff, 0xfe, 0, 0}
	expected = append(expected, utf8To32LittleEndianTestInputs[1]...)

	output, err := ConvertToUTF32(input, types.UTF_32LE, true)

	if !bytes.Equal(expected, output) || err != nil {
		t.Errorf(`ConvertToUTF32(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
	}
}

func TestConvertToUTF16LE(t *testing.T) {
	for idx := 0; idx < len(utf8To16LittleEndianTestInputs); idx += 2 {
		input := utf8To16LittleEndianTestInputs[idx]
		expected := utf8To16LittleEndianTestInputs[idx+1]
		output, err := ConvertToUTF16(input, types.UTF_16LE, false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF16(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

var utf8To32BigEndianTestInputs = [][]byte{
	{65}, {0, 0, 0, 65}, // 'A' (U+0041)
	{122}, {0, 0, 0, 122}, // 'z' (U+007A)
	{32}, {0, 0, 0, 32}, // space (U+0020)

	{195, 169}, {0, 0, 0, 233}, // é (U+00E9)
	{194, 169}, {0, 0, 0, 169}, // © (U+00A9)
	{194, 176}, {0, 0, 0, 176}, // ° (U+00B0)
	{195, 177}, {0, 0, 0, 241}, // ñ (U+00F1)

	{208, 144}, {0, 0, 4, 16}, // А (U+0410)
	{209, 129}, {0, 0, 4, 65}, // с (U+0441)

	{223, 191}, {0, 0, 7, 255}, // ߿ (U+07FF)
	{237, 159, 191}, {0, 0, 215, 255}, // ퟿ (U+D7FF)

	{240, 144, 128, 128}, {0, 1, 0, 0}, // 𐀀 (U+10000)
	{240, 144, 128, 129}, {0, 1, 0, 1}, // 𐀁 (U+10001)

	{240, 159, 152, 128}, {0, 1, 246, 0}, // 😀 (U+1F600)

	{244, 143, 191, 191}, {0, 16, 255, 255}, // U+10FFFF

	{239, 191, 189}, {0, 0, 255, 253}, // invalid surrogate → U+D83C? (kept BE bytes)
	{244, 129, 128, 128}, {0, 16, 16, 0}, // U+101000 (valid PUA)

	{225, 131, 191}, {0, 0, 16, 255}, // U+10FF

	{195, 191}, {0, 0, 0, 255}, // ÿ (U+00FF)

	{239, 191, 189}, {0, 0, 255, 253}, // U+FFFD replacement char

	{239, 191, 189}, {0, 0, 255, 253}, // invalid > U+10FFFF

	{224, 162, 128}, {0, 0, 8, 128}, // U+0880

	{226, 130, 172}, {0, 0, 32, 172}, // € (U+20AC)
	{226, 130, 185}, {0, 0, 32, 185}, // ₹ (U+20B9)

	{227, 128, 173}, {0, 0, 48, 45}, // ㌭ (U+302D)

	{}, {}, // Empty input

	{237, 156, 128}, {0, 0, 215, 0}, // U+D700 (valid BMP)
}

var utf8To32LittleEndianTestInputs = [][]byte{
	{65}, {65, 0, 0, 0}, // 'A' (U+0041)
	{122}, {122, 0, 0, 0}, // 'z' (U+007A)
	{32}, {32, 0, 0, 0}, // space (U+0020)

	{195, 169}, {233, 0, 0, 0}, // é (U+00E9)
	{194, 169}, {169, 0, 0, 0}, // © (U+00A9)
	{194, 176}, {176, 0, 0, 0}, // ° (U+00B0)
	{195, 177}, {241, 0, 0, 0}, // ñ (U+00F1)

	{208, 144}, {16, 4, 0, 0}, // А (U+0410)
	{209, 129}, {65, 4, 0, 0}, // с (U+0441)

	{223, 191}, {255, 7, 0, 0}, // ߿ (U+07FF)
	{237, 159, 191}, {255, 215, 0, 0}, // ퟿ (U+D7FF)

	{240, 144, 128, 128}, {0, 0, 1, 0}, // 𐀀 (U+10000)
	{240, 144, 128, 129}, {1, 0, 1, 0}, // 𐀁 (U+10001)

	{240, 159, 152, 128}, {0, 246, 1, 0}, // 😀 (U+1F600)

	{244, 143, 191, 191}, {255, 255, 16, 0}, // U+10FFFF

	{239, 191, 189}, {253, 255, 0, 0}, // U+FFFD (replacement)
	{244, 129, 128, 128}, {0, 16, 16, 0}, // U+101000 (valid PUA)

	{225, 131, 191}, {255, 16, 0, 0}, // U+10FF

	{195, 191}, {255, 0, 0, 0}, // ÿ (U+00FF)

	{239, 191, 189}, {253, 255, 0, 0}, // U+FFFD replacement

	{239, 191, 189}, {253, 255, 0, 0}, // invalid > U+10FFFF (→ U+FFFD)

	{224, 162, 128}, {128, 8, 0, 0}, // U+0880

	{226, 130, 172}, {172, 32, 0, 0}, // € (U+20AC)
	{226, 130, 185}, {185, 32, 0, 0}, // ₹ (U+20B9)

	{227, 128, 173}, {45, 48, 0, 0}, // ㌭ (U+302D)

	{}, {}, // Empty input

	{237, 156, 128}, {0, 215, 0, 0}, // U+D700 (valid BMP)
}

var utf8To16LittleEndianTestInputs = [][]byte{
	{65}, {65, 0}, // 'A' (U+0041)
	{122}, {122, 0}, // 'z' (U+007A)
	{32}, {32, 0}, // space (U+0020)

	{195, 169}, {233, 0}, // é (U+00E9)
	{194, 169}, {169, 0}, // © (U+00A9)
	{194, 176}, {176, 0}, // ° (U+00B0)
	{195, 177}, {241, 0}, // ñ (U+00F1)

	{208, 144}, {16, 4}, // А (U+0410)
	{209, 129}, {65, 4}, // с (U+0441)

	{223, 191}, {255, 7}, // ߿ (U+07FF)
	{237, 159, 191}, {255, 215}, // ퟿ (U+D7FF)

	{240, 144, 128, 128}, {0x00, 0xD8, 0x00, 0xDC}, // 𐀀 (U+10000)
	{240, 144, 128, 129}, {0x00, 0xD8, 0x01, 0xDC}, // 𐀁 (U+10001)

	{240, 159, 152, 128}, {0x3D, 0xD8, 0x00, 0xDE}, // 😀 (U+1F600)

	{244, 143, 191, 191}, {0xFF, 0xDB, 0xFF, 0xDF}, // U+10FFFF

	{239, 191, 189}, {253, 255}, // U+FFFD (replacement)

	{244, 129, 128, 128}, // U+101000 (valid PUA)
	{196, 219, 0, 220},   // surrogate pair for U+101000

	{225, 131, 191}, {255, 16}, // U+10FF

	{195, 191}, {255, 0}, // ÿ (U+00FF)

	{239, 191, 189}, {253, 255}, // U+FFFD replacement

	{239, 191, 189}, {253, 255}, // invalid → U+FFFD

	{224, 162, 128}, {128, 8}, // U+0880

	{226, 130, 172}, {172, 32}, // € (U+20AC)
	{226, 130, 185}, {185, 32}, // ₹ (U+20B9)

	{227, 128, 173}, {45, 48}, // ㌭ (U+302D)

	{}, {}, // empty input

	{237, 156, 128}, {0x00, 0xD7}, // U+D700 valid BMP
}
