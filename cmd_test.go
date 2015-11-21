package cronsifter

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const stdoutResult = "2015-01-05"
const stderrResult = "Error!"

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprintf(os.Stdout, stdoutResult)
	fmt.Fprintf(os.Stderr, stderrResult)

	os.Exit(0)
}

func TestExecCommandStdoutStderr(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	o := make(chan string)
	e := make(chan string)
	go ExecCommand([]string{"date"}, o, e)

	for s := range o {
		t.Logf("s: %v / %v", s, stdoutResult)
		if s != stdoutResult {
			t.Errorf("Stdout of [%s] doesn't match [%s]", s, stdoutResult)
		}
	}
	for s := range e {
		if s != stderrResult {
			t.Errorf("Stderr of [%s] doesn't match [%s]", s, stderrResult)
		}
	}
}
