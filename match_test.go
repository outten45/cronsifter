package cronsifter

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

func TestRegexDaisyChain(t *testing.T) {
	regexps := createRegexps()
	matcher := RegexChannels(regexps)
	var res string

	go func(c chan string) { c <- "unk" }(matcher.Right)
	res = <-matcher.Left
	if res != "unk" {
		t.Errorf("Expected the call to return \"unk\", but got [%s]", res)
	}

	go func(c chan string) { c <- "abcd" }(matcher.Right)
	res = <-matcher.Left
	if res != "" {
		t.Errorf("Expected the call to return an empty string, but we got [%s]", res)
	}

	go func(c chan string) { c <- "aoeuhaouhaeotuhalmn" }(matcher.Right)
	res = <-matcher.Left
	if res != "" {
		t.Errorf("Expected the call to return an empty string, but we got [%s]", res)
	}

	if matcher.CheckRegex("abc") == true {
		t.Error("Expected CheckRegex to return true, but it was false")
	}

	if matcher.CheckRegex("cba") == false {
		t.Error("Expected CheckRegex to return false, but it was true")
	}
}

func TestReadingRegexpFromFile(t *testing.T) {
	regexps := GetExcludesRegexps("files/exclude.txt")
	if len(regexps) < 2 {
		t.Error("File did not parse correctly. There should be more that 2 regexs.")
	}

	if !regexps[0].MatchString("this should be ignored.") {
		t.Errorf("Line with \"ignore\" in it should match.")
	}

}
