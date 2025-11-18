package war

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/doctordesh/war/colors"
)

type watchAndRun struct {
	watcher *watcher
	runner  *runner

	Verbose bool
}

func New(directoryToWatch string, runnable RunnableTemplate, delay, ignoreChangesFor time.Duration) *watchAndRun {
	w := &watcher{
		directory: directoryToWatch,
		exclude:   runnable.Excludes,
		verbose:   false,
	}
	r := &runner{
		runnableTemplate: runnable,
		delay:            delay,
		ignoreChangesFor: ignoreChangesFor,
	}

	return &watchAndRun{w, r, false}
}

func (w *watchAndRun) WatchAndRun() error {
	// w.watcher.SetVerboseLogging(w.Verbose)
	// w.runner.SetVerboseLogging(w.Verbose)

	// Setup signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// Run
	c, err := w.watcher.Watch()
	if err != nil {
		return err
	}

	go w.runner.Run(c)

	// Make an initial run
	w.runner.run()

	<-sigs

	fmt.Println()

	// If it's running, stop it. To not leak processes
	w.runner.Stop()

	colors.Blue("keyboard interrupt detected, quiting...")

	return nil
}
