package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	outF     string
	errF     string
	stdOut   io.Writer
	stdErr   io.Writer
	logO     *log.Logger
	logE     *log.Logger
	logFlags = log.LstdFlags | log.Lshortfile
)

func init() {
	flag.StringVar(&outF, "out", "", "Command's output file")
	flag.StringVar(&errF, "err", "", "Command's error output file")
}

func checkOutFiles(outF, errF string) error {
	if outF == "" && errF == "" {
		stdOut = ioutil.Discard
		stdErr = ioutil.Discard
		return nil
	}
	if outF != "" && errF == "" {
		errF = outF
	} else if errF != "" && outF == "" {
		outF = errF
	}
	var err error
	stdOut, err = os.Create(outF)
	if err != nil {
		return err
	}
	if outF != errF {
		stdErr, err = os.Create(errF)
	} else {
		stdErr = stdOut
	}
	return err
}

func initLog() {
	logO = log.New(os.Stdout, "", logFlags)
	logE = log.New(os.Stderr, "", logFlags)
}

func initLogWithFiles(stdO, stdE io.Writer) {
	outW := io.MultiWriter(os.Stdout, stdO)
	errW := io.MultiWriter(os.Stderr, stdE)
	logO = log.New(outW, "", logFlags)
	logE = log.New(errW, "", logFlags)
}

func main() {
	flag.Parse()
	err := checkOutFiles(outF, errF)
	if err != nil {
		initLog()
		logE.Fatalln("Error:", err)
	}
	initLogWithFiles(stdOut, stdErr)
	logO.Println("Printed to stdout")
	logE.Println("Printed to stderr")
}
