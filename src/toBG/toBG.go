package main

import (
	"flag"
)

var (
	outF string
	errF string
)

func init() {
	flag.StringVar(&outF, "out", "", "Command's output file")
	flag.StringVar(&errF, "err", "", "Command's error output file")
}

func checkOutFiles() {}

func main() {
	flag.Parse()
}
