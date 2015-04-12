package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	gl "github.com/siddontang/go-log/log"
)

var (
	excludes = flag.String("excludes", "", "the file containing the regular expressions to remove from stdout/stderr")
	logDir   = flag.String("dir", ".", "directory to write standard out and err log files to")
	logSize  = flag.Int("size", 20, "file size in megabytes to rotate stdout and stderr files")
	logCount = flag.Int("count", 20, "rotate log file count for stdout and stderr files")
	cmdArgs  = []string{}
)

func main() {
	flag.Parse()
	cmdArgs = flag.Args()
	// fmt.Println("cmdArgs>>>>>> ", cmdArgs)
	run()
}

func run() {
	setExcludesRegexps(*excludes)

	of := fmt.Sprintf("%s/%s.out.log", *logDir, getCmdName())
	stdoutHandler, err := gl.NewRotatingFileHandler(of, *logSize*1024*1024, *logCount)
	if err != nil {
		panic(err)
	}
	ef := fmt.Sprintf("%s/%s.err.log", *logDir, getCmdName())
	stderrHandler, err := gl.NewRotatingFileHandler(ef, *logSize*1024*1024, *logCount)
	if err != nil {
		panic(err)
	}

	stdoutLog := gl.New(stdoutHandler, gl.Llevel|gl.Ltime)
	stderrLog := gl.New(stderrHandler, gl.Llevel|gl.Ltime)
	defer stdoutHandler.Close()
	defer stderrHandler.Close()

	stdout, stderr := cmd(cmdArgs)
	stdoutLog.Info("STDOUT\n%s", stdout)
	stderrLog.Info("STDERR\n%s", stderr)
	output(stdout, stderr, excludesRegexps)
	time.Sleep(time.Second * 2)
}

func output(stdout, stderr string, excludes []*regexp.Regexp) {

	slines := strings.Split(stdout, "\n")
	for i := 0; i < len(slines)-1; i++ {
		s := slines[i]
		if !matchAnyExcludes(s) {
			fmt.Fprintf(os.Stdout, "%s\n", s)
		}
	}

	slines = strings.Split(stderr, "\n")
	for i := 0; i < len(slines)-1; i++ {
		s := slines[i]
		if !matchAnyExcludes(s) {
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
