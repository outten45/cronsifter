package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// Matcher the left and right channels for the Matcher
// to use.
type Matcher struct {
	Left  chan string
	Right chan string
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

// CheckRegex takes an in and out channel and returns a true if
// none of the regex match.
func CheckRegex(matcher *Matcher, line string) bool {
	// left <-chan string, right chan<- string
	go func(c chan<- string) { c <- line }(matcher.Right)
	res := <-matcher.Left
	return (res == line)
}
