package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(w io.Writer) error
}

type ProcessRunner struct {
	cmd        string
	args, envs []string
}

func NewProcessRunner(cmd string, args []string, envs []string) *ProcessRunner {
	return &ProcessRunner{cmd, args, envs}
}

func (pr *ProcessRunner) String() string {
	return strings.Join(pr.envs, " ") +
		" " + pr.cmd + " " + strings.Join(pr.args, " ")
}

func streamError(w io.Writer, e error) error {
	w.Write([]byte(e.Error()))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return e
}

func (pr *ProcessRunner) Run(w io.Writer) error {
	cmd := exec.Command(pr.cmd, pr.args...)
	cmd.Env = append(os.Environ(), pr.envs...)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}
