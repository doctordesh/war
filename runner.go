package war

import (
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
	verbose bool
}

type Runner interface {
	SetVerboseLogging(bool)
	Run(c <-chan string)
}

func NewRunner(cmd string, env []string, delay int) Runner {
	return &runner{cmd, env, delay, false}
}

func (r *runner) Run(c <-chan string) {

	now := func() int {
		return int(time.Now().UnixNano() / 1000000)
	}

	shouldRun := false
	timeSinceLastEvent := 0
	ds := 0

	for {

		ds = now()

		time.Sleep(time.Millisecond)

		select {
		case _, ok := <-c:
			if !ok {
				colors.Blue("event channel closed. quitting...")
				return
			}
			shouldRun = true
			timeSinceLastEvent = 0
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

			r.run()
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
	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.Env = append(cmd.Env, r.env...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		colors.Red("command '%s' exited with error: %v", r.cmd, err)
	} else {
		colors.Green("command '%s' successfull", r.cmd)
	}

	// Extra line between calls, helps with when skimming
	fmt.Println()
}
