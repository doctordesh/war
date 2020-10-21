package war

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type watchAndRun struct {
	watcher Watcher
	runner  Runner
	command string
	args    []string
	env     []string

	Verbose bool
}

func New(directory string, match []string, exclude []string, commandString string, env []string, delay int, timeout time.Duration) *watchAndRun {
	w := NewWatcher(directory, match, exclude)
	r := NewRunner()

	parts := strings.Split(commandString, " ")

	return &watchAndRun{
		watcher: w,
		runner:  r,
		command: parts[0],
		args:    parts[1:],
		env:     env,
		Verbose: false,
	}
}

func (w *watchAndRun) WatchAndRun() error {
	w.watcher.SetVerboseLogging(w.Verbose)
	w.runner.SetVerboseLogging(w.Verbose)

	c, err := w.watcher.Watch()
	if err != nil {
		return err
	}

	for file := range c {
		fmt.Printf("WatchAndRun with file: %s\n", file)
		if w.Verbose {
			fmt.Printf("file(s) changed: %s\n", file)
		}

		runnable := NewRunnable(
			w.command,
			w.args,
			os.Stdout,
			os.Stderr,
		)

		w.runner.Run(runnable)
	}

	return nil
}
