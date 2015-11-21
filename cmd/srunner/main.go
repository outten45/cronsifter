package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/outten45/cronsifter"
)

func main() {
	flag.CommandLine.Usage = func() {
		u := `srunner [command] [command args]

Will run the arguments after the srunner command. For example:

  srunner myGreatCommand -s

"myGreatCommand -s" is assume to run in the foreground and will be
restarted if it stops.
		`
		fmt.Println(u)
		return
	}
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
