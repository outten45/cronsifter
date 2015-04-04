package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "cronsifter"
	app.Usage = "filter stdout and stderr"
	app.Version = "0.0.1"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "excludes-file",
			Usage: "file of regular expressions to exclude from the stdout/stderr",
		},
	}

	app.Action = func(c *cli.Context) {
		run(c)
	}

	app.Run(os.Args)

}

func regexs(filename string) []*regexp.Regexp {
	rs := make([]*regexp.Regexp, 0)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		r, err := regexp.Compile(line)
		if err != nil {
			log.Printf("Regex didn't compile: %s\n", line)
		} else {
			rs = append(rs, r)
		}
	}

	return rs
}

func run(c *cli.Context) {
	excludes := make([]*regexp.Regexp, 0)

	ef := c.String("excludes-file")
	if len(ef) != 0 {
		excludes = regexs(ef)
	}

	stdout, stderr := cmd(c.Args())
	output(stdout, stderr, excludes)
}

func output(stdout, stderr string, excludes []*regexp.Regexp) {

	slines := strings.Split(stdout, "\n")
	for i := 0; i < len(slines)-1; i++ {
		s := slines[i]
		for _, r := range excludes {
			if !r.MatchString(s) {
				fmt.Fprintf(os.Stdout, "%s\n", s)
			}
		}
	}

	slines = strings.Split(stderr, "\n")
	for i := 0; i < len(slines)-1; i++ {
		s := slines[i]
		for _, r := range excludes {
			if !r.MatchString(s) {
				fmt.Fprintf(os.Stderr, "%s\n", s)
			}
		}
	}
}

func cmd(a []string) (o, e string) {
	cmd := exec.Command(a[0], a[1:]...)
	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
	return out.String(), err.String()
}
