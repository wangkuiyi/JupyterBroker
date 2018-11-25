package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessRunner(t *testing.T) {
	// Need shell for environment variable replacement.
	pr := NewProcessRunner("sh", []string{"-c", "echo hello $foo"}, []string{"foo=bar"})
	var buf bytes.Buffer
	assert.NoError(t, pr.Run(&buf))
	assert.Equal(t, "hello bar\n", buf.String())
}
