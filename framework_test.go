package jupyterbroker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrameworkRun(t *testing.T) {
	http.Handle("/myproc", &ProcessBroker{})
	addr, e := Start(":0")
	assert.NoError(t, e)

	res, e := http.Get(fmt.Sprintf(
		"http://%s/myproc?cmd=echo&args=hello&args=world", addr))
	assert.NoError(t, e)

	txt, e := ioutil.ReadAll(res.Body)
	assert.NoError(t, e)
	assert.Equal(t, "hello world\n", string(txt))
	res.Body.Close()
}
