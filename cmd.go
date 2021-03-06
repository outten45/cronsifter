package cronsifter

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"time"
)

var execCommand = exec.Command

func handleOutput(s *bufio.Scanner, out chan<- string) {
	defer close(out)
	for s.Scan() {
		text := s.Text()
		out <- text
	}
}

// ExecCommand take a slice of string to be executed as the
// command. stdout and stderr are passed to the channels handed
// to the function.
func ExecCommand(a []string, stdout chan<- string, stderr chan<- string) {
	cmd := execCommand(a[0], a[1:]...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("RunCommand: cmd.StdoutPipe(): %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("RunCommand: cmd.StderrPipe(): %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("RunCommand: cmd.Start(): %v", err)
	}

	scanner := bufio.NewScanner(stdoutPipe)
	go handleOutput(scanner, stdout)

	errScanner := bufio.NewScanner(stderrPipe)
	go handleOutput(errScanner, stderr)

	if err := cmd.Wait(); err != nil {
		log.Fatalf("RunCommand: cmd.Wait(): %v", err)
	}
}

// Process representing the process that is going to be run.
type Process struct {
	Command  string
	Args     []string
	Delay    string
	Pid      int
	Status   string
	OsP      *os.Process
	respawns int
}

// RunIt runs the Process and makes sure it keeps running.
func RunIt(p *Process) chan *Process {
	ch := make(chan *Process)
	go func() {
		p.run()
		go p.monitor()
		ch <- p
	}()

	return ch
}

func (p *Process) run() {
	wd, _ := os.Getwd()
	proc := &os.ProcAttr{
		Dir: wd,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}

	process, err := os.StartProcess(p.Command, p.Args, proc)
	if err != nil {
		log.Fatalf("%s failed. %s\n", p.Command, err)
		return
	}
	p.OsP = process
	p.Pid = process.Pid
	p.Status = "started"

}

func (p *Process) monitor() {
	if p.OsP == nil {
		return
	}
	status := make(chan *os.ProcessState)
	died := make(chan error)
	go func() {
		state, err := p.OsP.Wait()
		if err != nil {
			died <- err
			return
		}
		status <- state
	}()
	select {
	case s := <-status:
		log.Printf("%s exit=%s, success=%#v, exited=%#v, respawn_count=%#v\n", p.Command, s, s.Success(), s.Exited(), p.respawns)
		p.respawns++
		if p.Delay != "" {
			log.Printf("%s sleeping for %#v before restarting\n", p.Command, p.Delay)
			t, _ := time.ParseDuration(p.Delay)
			time.Sleep(t)
		}
		RunIt(p)
		p.Status = "restarted"
	case err := <-died:
		log.Printf("%d %s killed = %#v\n", p.OsP.Pid, p.Command, err)
	}
}
