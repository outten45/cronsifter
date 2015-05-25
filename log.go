package main

import (
	"fmt"
	"io"
	"strings"
)

// PrintCheckRegexp takes a Writer to print the string to if the
// regex isn't in.
func PrintCheckRegexp(w io.Writer, s string, regexIn, regexOut chan string) {
	slines := strings.Split(s, "\n")
	for i := 0; i < len(slines); i++ {
		s := slines[i]
		if CheckRegex(regexIn, regexOut, s) == true {
			fmt.Fprintf(w, "%s\n", s)
		}
	}
}
