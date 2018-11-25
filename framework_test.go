package jupyterbroker

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ProcessRunner struct {
	cmd        string
	args, envs []string
}

func (pr *ProcessRunner) Run(w io.Writer) {
	cmd := exec.Command(pr.cmd, pr.args...)
	cmd.Env = append(os.Environ(), pr.envs...)
	cmd.Stdout = w
	cmd.Stderr = w
	if e := cmd.Run(); e != nil {
		log.Panicf("ProcessRunner.Run: %v", e)
	}
}

func processRunnerHandler(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	pr := ProcessRunner{
		cmd:  req.Form["cmd"][0], // MakeSSEHandler will guard panics.
		args: req.Form["args"],
		envs: req.Form["envs"]}
	pr.Run(rw)
}

func TestProcessBroker(t *testing.T) {
	// Start the server.
	http.HandleFunc("/myproc", MakeSSEHandler(processRunnerHandler))
	addr, e := Start(":0")
	assert.NoError(t, e)

	// Ask the server to run echo.
	res, e := http.Get(fmt.Sprintf(
		"http://%s/myproc?cmd=echo&args=hello&args=world", addr))
	assert.NoError(t, e)
	txt, e := ioutil.ReadAll(res.Body)
	assert.NoError(t, e)
	assert.Equal(t, "hello world\n", string(txt))
	res.Body.Close()

	// Ask the server to run bash.
	res, e = http.Get(fmt.Sprintf(
		"http://%s/myproc?cmd=sh&args=%%2Dc&args=echo%%20hello%%20%%24foo&envs=foo%%3dbar", addr))
	assert.NoError(t, e)
	txt, e = ioutil.ReadAll(res.Body)
	assert.NoError(t, e)
	assert.Equal(t, "hello bar\n", string(txt))
	res.Body.Close()

	// Ask the server to run nothing.
	res, e = http.Get(fmt.Sprintf("http://%s/myproc", addr))
	assert.NoError(t, e) // NOTE: http.Error doesn't signal the error code.
	txt, e = ioutil.ReadAll(res.Body)
	assert.NoError(t, e)
	assert.True(t, strings.HasPrefix(string(txt), "runtime error:"))
	res.Body.Close()
}
