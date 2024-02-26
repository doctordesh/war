package war

import (
	"time"

	"github.com/doctordesh/war/colors"
)

type runner struct {
	runnableTemplate RunnableTemplate
	delay            time.Duration

	command *runnable
}

func (r *runner) Stop() {
	if r.command == nil {
		return
	}

	err := r.command.Stop()
	if err != nil {
		panic(err)
	}

}

func (r *runner) Run(changesHappened <-chan string) {
	var err error
	lastEventAt := time.Time{}

	for {
		select {
		case filename, ok := <-changesHappened:
			if !ok {
				colors.Blue("event channel closed. quitting...")
				return
			}

			colors.Blue("file changed: %s", filename)
			if lastEventAt.Add(r.delay).After(time.Now()) {
				colors.Yellow("ignoring")
				continue
			}

			lastEventAt = time.Now()

			// If there's already a running command, kill it and start over
			if r.command != nil {
				colors.Yellow("restarting command")
				err = r.command.Stop()
				if err != nil {
					panic(err)
				}
			}

			r.run()
		default:

			time.Sleep(time.Millisecond * 500)

			if r.command != nil && r.command.State() == RunningStateStopped {
				code, err := r.command.ExitCode()
				if err != nil {
					panic(err)
				}

				if code < 0 {
					colors.Red("command killed")
				} else if code > 0 {
					colors.Red("command exited with code %d", code)
				} else {
					colors.Green("command succesful")
				}

				r.command = nil
			}
		}
	}
}

// run ...
func (r *runner) run() {
	var err error

	r.command = r.runnableTemplate.Build()
	colors.Blue("running command: %s", r.command.cmd.String())
	err = r.command.Start()
	if err != nil {
		panic(err)
	}

}
