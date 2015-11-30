package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/outten45/cronsifter"
)

var (
	excludesFile = flag.String("excludes", "", "the file containing the regular expressions to remove from stdout/stderr")
	logDir       = flag.String("dir", ".", "directory to write standard out and err log files to")
	logSize      = flag.Int("size", 20, "file size in megabytes to rotate stdout and stderr files")
	logCount     = flag.Int("count", 20, "rotate log file count for stdout and stderr files")
	name         = flag.String("name", "NA", "the name of the service which is used events sent to the collectURL")
	collectURL   = flag.String("collecturl", "", "the URL to post the results of the cron")
	collectToken = flag.String("collecttoken", "", "the API token to send with the request")
	cmdArgs      = []string{}
)

func main() {
	flag.Parse()
	cmdArgs = flag.Args()

	stdoutLog, err := cronsifter.NewSimpleLogger(getStdoutLogFile(), (*logSize * 1024), *logCount)
	if err != nil {
		log.Fatalf("Error with stdout file: %v", err)
	}
	stderrLog, err := cronsifter.NewSimpleLogger(getStderrLogFile(), (*logSize * 1024), *logCount)
	if err != nil {
		log.Fatalf("Error with stderr file: %v", err)
	}

	regexps := cronsifter.GetExcludesRegexps(*excludesFile)
	matcher := cronsifter.RegexChannels(regexps)
	cronsifter.RunMatch(matcher, cmdArgs, stdoutLog, stderrLog)
}

func getStdoutLogFile() string {
	return path.Join(*logDir, fmt.Sprintf("%s.out.log", getCmdName()))
}

func getStderrLogFile() string {
	return path.Join(*logDir, fmt.Sprintf("%s.err.log", getCmdName()))
}

func getCmdName() string {
	p1 := os.Args[:1]
	parts := strings.Split(p1[0], "/")
	return parts[len(parts)-1]
}
