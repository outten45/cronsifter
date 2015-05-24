package main

import (
	"regexp"
	"testing"
)

func createRegexps() []*regexp.Regexp {
	var regexps = make([]*regexp.Regexp, 0)

	matchStrings := []string{"^abc", "^xyz", "lmn$"}
	// matchStrings := []string{"^abc"}

	for _, s := range matchStrings {
		r, _ := regexp.Compile(s)
		regexps = append(regexps, r)
	}
	return regexps
}

func TestRegexChannel(t *testing.T) {
	regexps := createRegexps()
	left, right := RegexChannels(regexps)
	var res string

	go func(c chan string) { c <- "unk" }(right)
	res = <-left
	if res != "unk" {
		t.Errorf("Expected the call to return \"unk\", but got [%s]", res)
	}

	go func(c chan string) { c <- "abcd" }(right)
	res = <-left
	if res != "" {
		t.Errorf("Expected the call to return an empty string, but we got [%s]", res)
	}

	go func(c chan string) { c <- "aoeuhaouhaeotuhalmn" }(right)
	res = <-left
	if res != "" {
		t.Errorf("Expected the call to return an empty string, but we got [%s]", res)
	}

	if CheckRegex(left, right, "abc") == true {
		t.Errorf("Expected CheckRegex to return true, but it was false")
	}
}
