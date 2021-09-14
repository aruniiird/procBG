package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

const (
	// ErrorLabel prefix label for error messages
	ErrorLabel = "[ ERROR ] "
	// VerboseLogFlags log flags for verbose mode
	VerboseLogFlags = log.LstdFlags | log.Lshortfile
)

var (
	outF     string
	errF     string
	stdOut   io.Writer
	stdErr   io.Writer
	logO     *log.Logger
	logE     *log.Logger
	logFlags = 0
	verbose  bool
)

func setOutErrFiles(outF, errF string) (stdOut, stdErr io.Writer, err error) {
	if outF == "" {
		stdOut = ioutil.Discard
	} else {
		if stdOut, err = os.Create(errF); err != nil {
			return nil, nil, err
		}
	}
	if errF == "" || errF == outF {
		stdErr = stdOut
	} else {
		if stdErr, err = os.Create(errF); err != nil {
			return nil, nil, err
		}
	}
	return
}

func initLog() {
	logO = log.New(os.Stdout, "", logFlags)
	logE = log.New(os.Stderr, ErrorLabel, VerboseLogFlags)
}

func initLogWithFiles(stdO, stdE io.Writer) {
	if !verbose {
		return
	}
	outW := io.MultiWriter(os.Stdout, stdO)
	errW := io.MultiWriter(os.Stderr, stdE)
	logO = log.New(outW, "", logFlags)
	logE = log.New(errW, ErrorLabel, VerboseLogFlags)
}

func init() {
	flag.StringVar(&outF, "out", "", "Command's output file")
	flag.StringVar(&errF, "err", "", "Command's error output file")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode")
	initLog()
}

func main() {
	flag.Parse()
	stdOut, stdErr, err := setOutErrFiles(outF, errF)
	if verbose {
		logFlags = VerboseLogFlags
	}
	if err != nil {
		logE.Fatalln("Error:", err)
	}
	initLogWithFiles(stdOut, stdErr)
	mainCmdStr := flag.Args()[0]
	otherArgs := flag.Args()[1:]
	if mCmd, err := exec.LookPath(mainCmdStr); err == nil {
		mainCmdStr = mCmd
	}
	mainCmd := exec.Command(mainCmdStr, otherArgs...)
	mainCmd.Stdout = stdOut
	mainCmd.Stderr = stdErr
	if verbose {
		logO.Println("Main Command    : ", mainCmdStr)
		logO.Println("Other Arguments : ", otherArgs)
		logO.Println("Arguments Length: ", len(otherArgs))
	}
	if err := mainCmd.Start(); err != nil {
		logE.Println("Failed to start the command: ", err)
		os.Exit(1)
	}
	logO.Println(mainCmd.Process.Pid)
}
