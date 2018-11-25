package jupyterbroker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(w io.Writer)
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

func (pr *ProcessRunner) Run(w io.Writer) {
	cmd := exec.Command(pr.cmd, pr.args...)
	cmd.Env = append(os.Environ(), pr.envs...)
	cmd.Stdout = w
	cmd.Stderr = w
	if e := cmd.Run(); e != nil {
		fmt.Fprintf(w, "ProcessRunner failed to run %s: %v", pr, e)
	}
}
