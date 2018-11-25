package jupyterbroker

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleProcessBroker() {
	http.HandleFunc("/process", MakeSSEHandler(ProcessRunnerHandler))
	addr, _ := Start(":0")

	res, _ := http.Get(fmt.Sprintf(
		"http://%s/process?cmd=echo&args=hello&args=world", addr))
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
	// Output:
	// hello world
}

func TestProcessBroker(t *testing.T) {
	// Start the server.
	http.HandleFunc("/myproc", MakeSSEHandler(ProcessRunnerHandler))
	addr, e := Start(":0")
	assert.NoError(t, e)

	// Ask the server to run bash.
	res, e := http.Get(fmt.Sprintf(
		"http://%s/myproc?cmd=sh&args=%%2Dc&args=echo%%20hello%%20%%24foo&envs=foo%%3dbar", addr))
	assert.NoError(t, e)
	txt, e := ioutil.ReadAll(res.Body)
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
