package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
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

	cout := make(chan string)
	cerr := make(chan string)

	go ExecCommand(cmdArgs, cout, cerr)

	var buffout bytes.Buffer
	for s := range cout {
		buffout.WriteString(fmt.Sprintln(s))
	}
	var bufferr bytes.Buffer
	for s := range cerr {
		bufferr.WriteString(fmt.Sprintln(s))
	}

	// stdoutLog.Info("STDOUT\n%s", stdout)
	// stderrLog.Info("STDERR\n%s", stderr)
	output(buffout.String(), bufferr.String(), regexIn, regexOut)
	// time.Sleep(time.Second * 2)
}

func output(stdout, stderr string, regexIn, regexOut chan string) {
	slines := strings.Split(stdout, "\n")
	for i := 0; i < len(slines); i++ {
		s := slines[i]
		if CheckRegex(regexIn, regexOut, s) == true {
			fmt.Fprintf(os.Stdout, "%s\n", s)
		}
	}

	slines = strings.Split(stderr, "\n")
	for i := 0; i < len(slines)-1; i++ {
		s := slines[i]
		if CheckRegex(regexIn, regexOut, s) {
			fmt.Fprintf(os.Stderr, "%s\n", s)
		}
	}
}

func cmd(a []string) (o, e string) {
	cmd := exec.Command(a[0], a[1:]...)
	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
	return out.String(), err.String()
}

func getCmdName() string {
	p1 := os.Args[:1]
	parts := strings.Split(p1[0], "/")
	return parts[len(parts)-1]
}
