package main

import "utfcoder/logger"

var fatal = logger.Fatal

func isValidEncoding(pEncoding string) bool {
	for _, encoding := range validEncodings {
		if encoding == pEncoding {
			return true
		}
	}
	return false
}

func RunPrechecks() {
	if len(sourceFile) == 0 {
		fatal("no source file path mentioned. use '-s filepath/filename' to mention source file path")
	}

	if len(fromEncoding) == 0 || !isValidEncoding(fromEncoding) {
		fatal("no (or) invalid source encoding provided. use '-from utf-8/utf-16/utf-32'")
	}

	if len(toEncoding) == 0 || !isValidEncoding(toEncoding) {
		fatal("no (or) invalid target encoding provided. use '-to utf-8/utf-16/utf-32'")
	}

	if fromEncoding == toEncoding {
		fatal("incorrect source/target encoding provided. cannot encode", fromEncoding, "again to", toEncoding)
	}
}
