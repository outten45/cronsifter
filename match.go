package main

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var (
	excludesRegexps = make([]*regexp.Regexp, 0)
)

func matchAnyExcludes(line string) bool {
	res := false
	for _, r := range excludesRegexps {
		if r.MatchString(line) {
			return true
		}
	}
	return res
}

func setExcludesRegexps(filename string) []*regexp.Regexp {
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
