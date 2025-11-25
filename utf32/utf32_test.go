package UTF32

import (
	"bytes"
	"testing"
	"utfcoder/types"
)

func TestConvertToUTF8(t *testing.T) {
	for idx := 0; idx < len(utf8TestInputs); idx += 2 {
		input := utf8TestInputs[idx]
		expected := utf8TestInputs[idx+1]
		output, err := ConvertToUTF8(input, false)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF8(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestInvalidInputConvertToUTF8(t *testing.T) {
	for idx := 0; idx < len(invalidTestInputs); idx += 2 {
		input := invalidTestInputs[idx]
		expected := invalidTestInputs[idx+1]
		output, err := ConvertToUTF8(input, false)

		if !bytes.Equal(expected, output) || err == nil {
			t.Errorf(`ConvertToUTF8(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), "invalid input")
		}
	}
}

func TestConvertToUTF16(t *testing.T) {
	for idx := 0; idx < len(utf16TestInputs); idx += 2 {
		input := utf16TestInputs[idx]
		expected := utf16TestInputs[idx+1]
		output, err := ConvertToUTF16(input, types.UTF_16LE, true)

		if !bytes.Equal(expected, output) || err != nil {
			t.Errorf(`ConvertToUTF16(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), nil)
		}
	}
}

func TestInvalidInputConvertToUTF16(t *testing.T) {
	for idx := 0; idx < len(invalidTestInputs); idx += 2 {
		input := invalidTestInputs[idx]
		expected := invalidTestInputs[idx+1]
		output, err := ConvertToUTF16(input, types.UTF_16LE, false)

		if !bytes.Equal(expected, output) || err == nil {
			t.Errorf(`ConvertToUTF8(%v) = output=%v (%v), error=%v, Expected = output=%v (%v), error=%v`, input, output, string(output), err, expected, string(expected), "invalid input")
		}
	}
}

// invalid UTF-32 inputs
var invalidTestInputs = [][]byte{
	{}, {},
	{0, 0, 0}, {},
}

// UTF-32 inputs (mixed endianness) and their correct UTF-16 LE outputs.
var utf16TestInputs = [][]byte{
	// BMP - Single 16-bit unit (0x0000 - 0xFFFF)
	{65, 0, 0, 0}, {255, 254, 65, 0}, // 'A' (U+0041)
	{0, 0, 0, 65}, {255, 254, 65, 0}, // 'A' BE
	{122, 0, 0, 0}, {255, 254, 122, 0}, // 'z' (U+007A)
	{0, 0, 0, 122}, {255, 254, 122, 0}, // 'z' BE
	{32, 0, 0, 0}, {255, 254, 32, 0}, // space U+0020
	{0, 0, 0, 32}, {255, 254, 32, 0}, // space BE
	{233, 0, 0, 0}, {255, 254, 233, 0}, // é U+00E9
	{0, 0, 0, 233}, {255, 254, 233, 0}, // é BE
	{169, 0, 0, 0}, {255, 254, 169, 0}, // © U+00A9
	{0, 0, 0, 169}, {255, 254, 169, 0}, // © BE
	{176, 0, 0, 0}, {255, 254, 176, 0}, // ° U+00B0
	{0, 0, 0, 176}, {255, 254, 176, 0}, // ° BE
	{241, 0, 0, 0}, {255, 254, 241, 0}, // ñ U+00F1
	{0, 0, 0, 241}, {255, 254, 241, 0}, // ñ BE

	{16, 4, 0, 0}, {255, 254, 16, 4}, // А U+0410
	{0, 0, 4, 16}, {255, 254, 16, 4}, // А BE
	{65, 4, 0, 0}, {255, 254, 65, 4}, // с U+0441
	{0, 0, 4, 65}, {255, 254, 65, 4}, // с BE
	{255, 7, 0, 0}, {255, 254, 255, 7}, // ߿ U+07FF
	{0, 0, 7, 255}, {255, 254, 255, 7}, // BE
	{255, 215, 0, 0}, {255, 254, 255, 215}, // ퟿ U+D7FF (last before surrogates)
	{0, 0, 215, 255}, {255, 254, 255, 215}, // BE
	{128, 8, 0, 0}, {255, 254, 128, 8}, // U+0880
	{0, 0, 8, 128}, {255, 254, 128, 8}, // BE
	{172, 32, 0, 0}, {255, 254, 172, 32}, // €
	{0, 0, 32, 172}, {255, 254, 172, 32}, // BE
	{185, 32, 0, 0}, {255, 254, 185, 32}, // ₹
	{0, 0, 32, 185}, {255, 254, 185, 32}, // BE
	{45, 48, 0, 0}, {255, 254, 45, 48}, // ㌭
	{0, 0, 48, 45}, {255, 254, 45, 48}, // BE
	{0, 0, 215, 0}, {255, 254, 0, 215}, // U+D700
	{0, 215, 0, 0}, {255, 254, 0, 215}, // U+D700 LE

	// supplementary (> 0xFFFF)
	{0, 0, 1, 0}, {255, 254, 0, 216, 0, 220}, // U+10000
	{0, 1, 0, 1}, {255, 254, 0, 216, 1, 220}, // U+10001
	{0, 1, 246, 0}, {255, 254, 61, 216, 0, 222}, // 😀 BE
	{0, 246, 1, 0}, {255, 254, 61, 216, 0, 222}, // 😀 LE
	{0, 16, 255, 255}, {255, 254, 255, 219, 255, 223}, // U+10FFFF (BE)
	{255, 255, 16, 0}, {255, 254, 255, 219, 255, 223}, // U+10FFFF (LE)
	{60, 216, 0, 0}, {255, 254, 253, 255}, // invalid surrogate
	{0, 0, 216, 60}, {255, 254, 253, 255}, // invalid surrogate BE
	{0, 0, 17, 0}, {255, 254, 253, 255}, // invalid (> U+10FFFF)
	{255, 254, 0, 0}, {255, 254}, // UTF-32LE BOM – ignored → only BOM
	{0, 0, 254, 255}, {255, 254}, // UTF-32BE BOM – ignored → only BOM
	{0, 16, 16, 0}, {255, 254, 196, 219, 0, 220}, // U+101000
	{255, 16, 0, 0}, {255, 254, 255, 16}, // U+10FF
	{0, 0, 0, 255}, {255, 254, 255, 0}, // ÿ U+00FF
	{253, 255, 0, 0}, {255, 254, 253, 255}, // U+FFFD
	{0, 0, 255, 253}, {255, 254, 253, 255}, // U+FFFD BE
}

// UTF-32 inputs (mixed endianness) and their correct UTF-8 outputs.
var utf8TestInputs = [][]byte{
	{65, 0, 0, 0}, {65}, // 'A' (U+0041, Latin Capital A)
	{0, 0, 0, 65}, {65}, // 'A' (U+0041, BE)
	{122, 0, 0, 0}, {122}, // 'z' (U+007A, Latin Small z)
	{0, 0, 0, 122}, {122}, // 'z' (U+007A, BE)
	{32, 0, 0, 0}, {32}, // space (U+0020)
	{0, 0, 0, 32}, {32}, // space (BE)

	{233, 0, 0, 0}, {195, 169}, // é (U+00E9, Latin small e with acute)
	{0, 0, 0, 233}, {195, 169}, // é (BE)
	{169, 0, 0, 0}, {194, 169}, // © (U+00A9, Copyright sign)
	{0, 0, 0, 169}, {194, 169}, // © (BE)
	{176, 0, 0, 0}, {194, 176}, // ° (U+00B0, Degree sign)
	{0, 0, 0, 176}, {194, 176}, // ° (BE)
	{241, 0, 0, 0}, {195, 177}, // ñ (U+00F1, Latin small n with tilde)
	{0, 0, 0, 241}, {195, 177}, // ñ (BE)

	{16, 4, 0, 0}, {208, 144}, // А (U+0410, Cyrillic Capital A)
	{0, 0, 4, 16}, {208, 144}, // А (BE)
	{65, 4, 0, 0}, {209, 129}, // с (U+0441, Cyrillic small es)
	{0, 0, 4, 65}, {209, 129}, // с (BE)
	{255, 7, 0, 0}, {223, 191}, // ߿ (U+07FF, last 2-byte codepoint)
	{0, 0, 7, 255}, {223, 191}, // ߿ (BE)
	{255, 215, 0, 0}, {237, 159, 191}, // ퟿ (U+D7FF, last BMP before surrogates)
	{0, 0, 215, 255}, {237, 159, 191}, // ퟿ (BE)

	{0, 0, 1, 0}, {240, 144, 128, 128}, // 𐀀 (U+10000, Linear B syllable)
	{0, 1, 0, 1}, {240, 144, 128, 129}, // 𐀁 (U+10001, Linear B Syllable B008A)
	{0, 1, 246, 0}, {240, 159, 152, 128}, // 😀 (U+1F600, grinning face) BE
	{0, 246, 1, 0}, {240, 159, 152, 128}, // 😀 (U+1F600, grinning face) LE
	{0, 16, 255, 255}, {244, 143, 191, 191}, // U+10FFFF (BE)
	{255, 255, 16, 0}, {244, 143, 191, 191}, // U+10FFFF (LE)

	{60, 216, 0, 0}, {239, 191, 189}, // invalid surrogate (LE)
	{0, 0, 216, 60}, {239, 191, 189}, // invalid surrogate (BE)
	{0, 16, 16, 0}, {244, 129, 128, 128}, // U+101000 (valid, Plane 16 Private-Use Area), UTF-8 = F4 81 80 80
	{255, 16, 0, 0}, {225, 131, 191}, // U+10FF valid (within BMP)
	{0, 0, 0, 255}, {195, 191}, // ÿ (U+00FF, Latin small y with diaeresis)
	{253, 255, 0, 0}, {239, 191, 189}, // U+FFFD replacement character
	{0, 0, 255, 253}, {239, 191, 189}, // U+FFFD (BE)
	{0, 0, 17, 0}, {239, 191, 189}, // U+110000 invalid (> U+10FFFF) → replaced � (LE)

	{128, 8, 0, 0}, {224, 162, 128}, // U+0880 (ࠀ) — UTF-8: [224,162,128]
	{0, 0, 8, 128}, {224, 162, 128}, // U+0880 (BE) same code point → same UTF-8
	{172, 32, 0, 0}, {226, 130, 172}, // € (U+20AC, Euro sign)
	{0, 0, 32, 172}, {226, 130, 172}, // € (BE)
	{185, 32, 0, 0}, {226, 130, 185}, // ₹ (U+20B9, Indian Rupee sign)
	{0, 0, 32, 185}, {226, 130, 185}, // ₹ (BE)
	{45, 48, 0, 0}, {227, 128, 173}, // ㌭ (U+302D)
	{0, 0, 48, 45}, {227, 128, 173}, // ㌭ (BE)
	{255, 254, 0, 0}, {}, // UTF-32LE BOM → ignored (encoding marker, no UTF-8 output)
	{0, 0, 254, 255}, {}, // UTF-32BE BOM → ignored (encoding marker, no UTF-8 output)

	{0, 0, 215, 0}, {237, 156, 128}, // U+D700 valid (just below surrogate range), UTF-8 = ED 9C 80
	{0, 215, 0, 0}, {237, 156, 128}, // U+D700 (LE)
}
