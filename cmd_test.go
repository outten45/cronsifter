package main

import (
	"regexp"
	"testing"
)

func TestExecCommand(t *testing.T) {
	o := make(chan string)
	e := make(chan string)
	go ExecCommand([]string{"date"}, o, e)

	match := false
	for s := range o {
		match, _ = regexp.MatchString("\\s[\\d]{4}", s)
		if !match {
			t.Errorf("4 digits in a row not found in date: %s", s)
		}
	}
	for s := range e {
		t.Errorf("stderr found and it shouldn't be: %v", s)
	}
}
