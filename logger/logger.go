package logger

import (
	"flag"
	"fmt"
	"os"
)

var verbose = flag.Bool("verbose", false, "")

func Log(items ...any) {
	if *verbose {
		fmt.Println(items...)
	}
}

func Fatal(items ...any) {
	fmt.Println(items...)
	os.Exit(1)
}
