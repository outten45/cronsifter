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

	regexps := GetExcludesRegexps(*excludesFile)
	matcher := RegexChannels(regexps)

	done := make(chan bool)
	cout := make(chan string)
	cerr := make(chan string)

	go ExecCommand(cmdArgs, cout, cerr)

	go func() {
		for s := range cout {
			matcher.PrintCheckRegexp(os.Stdout, s)
		}
		done <- true
	}()
	go func() {
		for s := range cerr {
			matcher.PrintCheckRegexp(os.Stderr, s)
		}
		done <- true
	}()

	<-done
	<-done
	close(done)
}

func getCmdName() string {
	p1 := os.Args[:1]
	parts := strings.Split(p1[0], "/")
	return parts[len(parts)-1]
}
