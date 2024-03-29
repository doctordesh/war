package war

import (
	"time"
)

type watchAndRun struct {
	watcher Watcher
	runner  *runner

	Verbose bool
}

func New(directory string, match []string, exclude []string, commandString string, env []string, delay int, timeout time.Duration) *watchAndRun {
	w := NewWatcher(directory, match, exclude)
	r := NewRunner(commandString, env, delay, timeout)
	return &watchAndRun{w, r, false}
}

func (w *watchAndRun) WatchAndRun() error {
	w.watcher.SetVerboseLogging(w.Verbose)
	w.runner.SetVerboseLogging(w.Verbose)
	c, err := w.watcher.Watch()
	if err != nil {
		return err
	}

	go w.runner.Run(c)

	// Run the command initially
	w.runner.run()

	return nil
}
