package main

import (
	"fmt"
	"testing"
)

func MockFatal() bool {
	return true
}

func TestNoSourceFileRunPrechecks(t *testing.T) {
	var fatalMessage = ""
	sourceFile, targetFile, fromEncoding, toEncoding = "", "file2", "utf-32", "utf-8"
	fatal = func(items ...any) {
		fatalMessage = items[0].(string)
	}
	expectedFatalMessage := "no source file path mentioned. use '-s filepath/filename' to mention source file path"

	RunPrechecks()
	if fatalMessage != expectedFatalMessage {
		t.Errorf(`RunPrechecks() = error=%v, Expected = error=%v`, fatalMessage, expectedFatalMessage)
	}
}

func TestNoFromEncodingRunPrechecks(t *testing.T) {
	var fatalMessage = ""
	sourceFile, targetFile, fromEncoding, toEncoding = "file", "file2", "", "utf-8"
	fatal = func(items ...any) {
		fatalMessage = items[0].(string)
	}
	expectedFatalMessage := "no (or) invalid source encoding provided. use '-from utf-8/utf-16/utf-32'"

	RunPrechecks()
	if fatalMessage != expectedFatalMessage {
		t.Errorf(`RunPrechecks() = error=%v, Expected = error=%v`, fatalMessage, expectedFatalMessage)
	}
}

func TestNoToEncodingRunPrechecks(t *testing.T) {
	var fatalMessage = ""
	sourceFile, targetFile, fromEncoding, toEncoding = "file", "file2", "utf-32", ""
	fatal = func(items ...any) {
		fatalMessage = items[0].(string)
	}
	expectedFatalMessage := "no (or) invalid target encoding provided. use '-to utf-8/utf-16/utf-32'"

	RunPrechecks()
	if fatalMessage != expectedFatalMessage {
		t.Errorf(`RunPrechecks() = error=%v, Expected = error=%v`, fatalMessage, expectedFatalMessage)
	}
}

func TestSameFromToEncodingRunPrechecks(t *testing.T) {
	var fatalMessage = ""
	sourceFile, targetFile, fromEncoding, toEncoding = "file", "file2", "utf-32", "utf-32"
	fatal = func(items ...any) {
		for idx, item := range items {
			fatalMessage += item.(string)
			if idx != len(items)-1 {
				fatalMessage += " "
			}
		}
	}
	expectedFatalMessage := fmt.Sprintf("incorrect source/target encoding provided. cannot encode %v again to %v", fromEncoding, toEncoding)

	RunPrechecks()
	if fatalMessage != expectedFatalMessage {
		t.Errorf(`RunPrechecks() = error=%v, Expected = error=%v`, fatalMessage, expectedFatalMessage)
	}
}
