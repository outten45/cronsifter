package main

import (
	"flag"
	"os"

	"github.com/outten45/cronsifter"
)

func main() {
	flag.Parse()
	cmdArgs := flag.Args()

	delay := os.Getenv("SPAWN_DELAY")
	if delay == "" {
		delay = "10s"
	}

	p := &cronsifter.Process{
		Command: cmdArgs[0],
		Args:    cmdArgs[1:],
		Delay:   delay,
	}
	ch := cronsifter.RunIt(p)
	<-ch
	<-ch

}
