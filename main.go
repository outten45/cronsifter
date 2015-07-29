package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
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

	stdoutLog, err := NewSimpleLogger(getStdoutLogFile(), 1000000, 20)
	if err != nil {
		log.Fatalf("Error with stdout file: %v", err)
	}
	stderrLog, err := NewSimpleLogger(getStderrLogFile(), 1000000, 20)
	if err != nil {
		log.Fatalf("Error with stderr file: %v", err)
	}

	done := make(chan bool)
	cout := make(chan string)
	cerr := make(chan string)

	if len(cmdArgs) > 0 {
		go ExecCommand(cmdArgs, cout, cerr)
	} else {
		go readStdin(cout)
		close(cerr)
	}

	go func() {
		for s := range cout {
			matcher.PrintCheckRegexp(os.Stdout, s)
			stdoutLog.Write([]byte(s))
		}
		done <- true
	}()
	go func() {
		for s := range cerr {
			matcher.PrintCheckRegexp(os.Stderr, s)
			stderrLog.Write([]byte(s))
		}
		done <- true
	}()

	<-done
	<-done
	close(done)
}

func readStdin(stdin chan<- string) {
	defer close(stdin)
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			// close(stdin)
			break
		}
		stdin <- strings.TrimRight(input, "\n")
	}
}

func getStdoutLogFile() string {
	return fmt.Sprintf("%s.out.log", getCmdName())
}

func getStderrLogFile() string {
	return fmt.Sprintf("%s.err.log", getCmdName())
}

func getCmdName() string {
	p1 := os.Args[:1]
	parts := strings.Split(p1[0], "/")
	return parts[len(parts)-1]
}
