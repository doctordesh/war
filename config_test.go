package war

import (
	"bytes"
	"testing"

	"github.com/doctordesh/check"
)

func TestConfig(t *testing.T) {
	data := `
directory: "."
command: "go test './streamer'"
env:
    - "SOMETHING=asdf"
delay: 100
timeout: 10
`
	b := bytes.NewBufferString(data)
	c, err := ConfigFromFile(b)
	check.OK(t, err)
	check.Equals(t, c.Directory, ".")
	check.Equals(t, c.CommandString, "go test './streamer'")
	check.Equals(t, c.Env, []string{"SOMETHING=asdf"})
	check.Equals(t, c.Delay, 100)
	check.Equals(t, c.Timeout, 10)
}
