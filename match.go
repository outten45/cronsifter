package cronsifter

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/outten45/cronsifter/collector"
)

type Match interface {
	CheckRegex(string) bool
	PrintCheckRegexp(io.Writer, string)
}

// Matcher the left and right channels for the Match
// to use.
type Matcher struct {
	Left  chan string
	Right chan string
}

// CheckRegex takes an in and out channel and returns a true if
// none of the regex match.
func (m *Matcher) CheckRegex(line string) bool {
	// left <-chan string, right chan<- string
	go func(c chan<- string) { c <- line }(m.Right)
	res := <-m.Left
	return (res == line)
}

// PrintCheckRegexp takes a Writer to print the string to if the
// regex isn't in.
func (m *Matcher) PrintCheckRegexp(w io.Writer, s string) {
	slines := strings.Split(s, "\n")
	for i := 0; i < len(slines); i++ {
		s := slines[i]
		if m.CheckRegex(s) == true {
			fmt.Fprintf(w, "%s\n", s)
		}
	}
}

func matchLine(r *regexp.Regexp, left chan<- string, right <-chan string) {
	for {
		line := <-right
		if !r.MatchString(line) {
			left <- line
		} else {
			left <- ""
		}
	}
}

// GetExcludesRegexps is given a filename and returns a slice of
// regular expressions for each line in the file.
func GetExcludesRegexps(filename string) []*regexp.Regexp {
	var excludesRegexps = make([]*regexp.Regexp, 0)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		r, err := regexp.Compile(line)
		if err != nil {
			log.Printf("Regex didn't compile: %s\n", line)
		} else {
			excludesRegexps = append(excludesRegexps, r)
		}
	}

	return excludesRegexps
}

// RegexChannels takes an array of regular expressions that create
// a chain of go routines. A channel is return that can be passed
// string that evaluates all the expressions.  If none of the
// expressions are matched, the string is returned; otherwise a
// blank string is return
func RegexChannels(regexps []*regexp.Regexp) *Matcher {
	leftmost := make(chan string)
	left := leftmost
	right := leftmost
	for _, r := range regexps {
		right = make(chan string)
		go matchLine(r, left, right)
		left = right
	}
	return &Matcher{Left: leftmost, Right: right}
}

// RunMatch take a matcher and cmd args, but run them.
func RunMatch(matcher *Matcher, cmdArgs []string, stdoutLog, stderrLog *SimpleLogger) {
	v := []string{}

	done := make(chan bool)
	cout := make(chan string)
	cerr := make(chan string)
	events := make(chan *collector.Event)

	if len(cmdArgs) > 0 {
		go ExecCommand(cmdArgs, cout, cerr, events)
	} else {
		go readStdin(cout)
		close(cerr)
	}

	go func() {
		for s := range cout {
			matcher.PrintCheckRegexp(os.Stdout, s)
			stdoutLog.Write([]byte(s))
			if len(v) >= 20 {
				_, v = v[0], v[1:]
			}
			v = append(v, s)
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
			break
		}
		stdin <- strings.TrimRight(input, "\n")
	}
}
