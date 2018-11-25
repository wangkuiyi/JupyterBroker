package jupyterbroker

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessRunner(t *testing.T) {
	pr := NewProcessRunner("sh",
		[]string{"-c", "echo hello $foo"}, []string{"foo=bar"})
	var buf bytes.Buffer
	pr.Run(&buf)
	assert.Equal(t, "hello bar\n", buf.String())

	pr = NewProcessRunner("sh", []string{"-c", "echo hello $foo"}, nil)
	buf.Reset()
	pr.Run(&buf)
	assert.Equal(t, "hello\n", buf.String())

	pr = NewProcessRunner("echo", nil, nil)
	buf.Reset()
	pr.Run(&buf)
	assert.Equal(t, "\n", buf.String())
}
