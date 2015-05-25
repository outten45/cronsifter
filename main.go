package main

import (
	"flag"
	"os"
	"strings"
)

var (
	excludesFile = flag.String("excludes", "", "the file containing the regular expressions to remove from stdout/stderr")
	logDir       = flag.String("dir", ".", "directory to write standard out and err log files to")
	logSize      = flag.Int("size", 20, "file size in megabytes to rotate stdout and stderr files")
	logCount     = flag.Int("count", 20, "rotate log file count for stdout and stderr files")
	cmdArgs      = []string{}
)

func main() {
	flag.Parse()
	cmdArgs = flag.Args()
	// fmt.Println("cmdArgs>>>>>> ", cmdArgs)
	run()
}

func run() {
	regexps := GetExcludesRegexps(*excludesFile)
	regexIn, regexOut := RegexChannels(regexps)

	dout := make(chan string)
	derr := make(chan string)
	cout := make(chan string)
	cerr := make(chan string)

	go ExecCommand(cmdArgs, cout, cerr)

	go func() {
		for s := range cout {
			PrintCheckRegexp(os.Stdout, s, regexIn, regexOut)
		}
		dout <- "done"
	}()
	go func() {
		for s := range cerr {
			PrintCheckRegexp(os.Stderr, s, regexIn, regexOut)
		}
		derr <- "done"
	}()

	<-dout
	<-derr
}

func getCmdName() string {
	p1 := os.Args[:1]
	parts := strings.Split(p1[0], "/")
	return parts[len(parts)-1]
}
