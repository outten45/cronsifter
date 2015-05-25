package main

import (
	"bufio"
	"log"
	"os/exec"
)

func handleOutput(s *bufio.Scanner, out chan<- string) {
	defer close(out)
	for s.Scan() {
		text := s.Text()
		out <- text
	}
}

// ExecCommand take a slice of string to be executed as the
// command.
func ExecCommand(a []string, stdout chan<- string, stderr chan<- string) {
	cmd := exec.Command(a[0], a[1:]...)

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
