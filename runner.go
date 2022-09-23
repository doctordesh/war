package war

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/doctordesh/war/colors"
)

type runner struct {
	cmd     string
	env     []string
	delay   int
	timeout time.Duration
	verbose bool
}

func NewRunner(cmd string, env []string, delay int, timeout time.Duration) *runner {
	return &runner{cmd, env, delay, timeout, false}
}

func (r *runner) Run(c <-chan string) {

	now := func() int {
		return int(time.Now().UnixNano() / 1000000)
	}

	shouldRun := false
	timeSinceLastEvent := 0
	ds := 0

	doneChan := make(chan struct{})
	running := false

	for {

		ds = now()

		time.Sleep(time.Millisecond)

		select {
		case filename, ok := <-c:
			if !ok {
				colors.Blue("event channel closed. quitting...")
				return
			}
			if running {
				continue
			}
			shouldRun = true
			timeSinceLastEvent = 0
			colors.Blue("file changed: %s", filename)
		case <-doneChan:
			running = false
		default:
			if shouldRun == false {
				continue
			}

			// Waiting for delay to trigger
			if timeSinceLastEvent < r.delay {
				dt := now() - ds
				timeSinceLastEvent += dt
				continue
			}

			// run and reset
			if r.verbose {
				log.Println("runner - delay triggered, running")
			}

			running = true
			go func() {
				r.run()
				doneChan <- struct{}{}
			}()

			shouldRun = false
			timeSinceLastEvent = 0
		}
	}
}

func (r *runner) SetVerboseLogging(b bool) {
	r.verbose = b
}

func (r *runner) run() {
	parts := strings.Split(r.cmd, " ")

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	cmd.Env = append(cmd.Env, r.env...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	colors.Blue("running '%s'", r.cmd)

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		colors.Red("command took too long (%d s timeout), killed", r.timeout/time.Second)
	} else if err != nil {
		colors.Red("command exited with error: %v", err)
	} else {
		colors.Green("command successfull")
	}

	// Extra line between calls, helps with when skimming
	fmt.Println()
}
