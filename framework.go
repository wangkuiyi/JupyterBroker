package jupyterbroker

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

// Runner defines an action which writes output to io.Writer.
type Runner interface {
	Run(w io.Writer)
}

// ProcessRunner is an example implementation of Runner.
type ProcessRunner struct {
	cmd        string
	args, envs []string
}

// Run a command in a sub-process, writing stdout and stderr to w.
// With whatever error, just panic it.
func (pr *ProcessRunner) Run(w io.Writer) {
	cmd := exec.Command(pr.cmd, pr.args...)
	cmd.Env = append(os.Environ(), pr.envs...)
	cmd.Stdout = w
	cmd.Stderr = w
	if e := cmd.Run(); e != nil {
		log.Panicf("ProcessRunner.Run: %v", e)
	}
}

// ProcessRunnerHandler is an example broker handler, which creates a
// ProcessRunner by parsing the URL and runs the runner.
func ProcessRunnerHandler(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	pr := ProcessRunner{
		cmd:  req.Form["cmd"][0], // MakeSSEHandler will guard panics.
		args: req.Form["args"],
		envs: req.Form["envs"]}
	pr.Run(rw)
}

// SetServerSentEventHeader marks an http.ResponseWriter of SSEs.
func SetServerSentEventHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// MakeSSEHandler returns a handler that guards panics in the given handler.
func MakeSSEHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, fmt.Sprintf("%v", e),
					http.StatusInternalServerError)
			}
		}()
		SetServerSentEventHeader(w)
		f(w, r)
	}
}

func Start(addr string) (string, error) {
	lst, e := net.Listen("tcp", addr)
	if e != nil {
		return "", e
	}
	go func() {
		log.Fatal("HTTP server error: ", http.Serve(lst, nil))
	}()
	return lst.Addr().String(), nil
}
