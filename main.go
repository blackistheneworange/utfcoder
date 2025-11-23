package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"utfcoder/logger"
	"utfcoder/types"
	UTF32 "utfcoder/utf32"
)

var sourceFileFlag = flag.String("s", "", "source file to read")
var targetFileFlag = flag.String("t", "", "target file to write")

var fromEncodingFlag = flag.String("from", "", "source file encoding")
var toEncodingFlag = flag.String("to", "", "target file encoding")

var validEncodings = [5]string{types.UTF_8, types.UTF_16, types.UTF_16BE, types.UTF_16LE, types.UTF_32}
var sourceFile, targetFile, fromEncoding, toEncoding string

func main() {
	flag.Parse()

	sourceFile, targetFile, fromEncoding, toEncoding = *sourceFileFlag, *targetFileFlag, strings.ToLower(*fromEncodingFlag), strings.ToLower(*toEncodingFlag)

	RunPrechecks()

	sourceFilePath, sourceFilePathErr := filepath.Abs(sourceFile)
	targetFilePath, targetFilePathErr := filepath.Abs(targetFile)
	if sourceFilePathErr != nil {
		logger.Fatal(sourceFilePathErr)
	}

	data, readErr := os.ReadFile(sourceFilePath)
	if readErr != nil {
		logger.Fatal(readErr)
	}

	var output []byte
	var err error

	switch fromEncoding {
	case types.UTF_32:
		if toEncoding == types.UTF_8 {
			output, err = UTF32.ConvertToUTF8(data)
		} else if toEncoding == types.UTF_16 || toEncoding == types.UTF_16LE || toEncoding == types.UTF_16BE {
			output, err = UTF32.ConvertToUTF16(data, toEncoding)
		} else {
			logger.Fatal(strings.ToUpper(types.UTF_32), "to", strings.ToUpper(toEncoding), "not implemented")
		}
	default:
		logger.Fatal(strings.ToUpper(fromEncoding), "to", strings.ToUpper(toEncoding), "not implemented")
	}

	if err != nil {
		logger.Fatal(err)
	}

	if len(targetFile) == 0 || targetFilePathErr != nil {
		os.Stdout.Write(output)
	} else {
		os.WriteFile(targetFilePath, output, 0600)
	}
}
