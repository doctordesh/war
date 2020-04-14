package war

type watchAndRun struct {
	watcher Watcher
	runner  Runner

	Verbose bool
}

func New(directory string, match []string, exclude []string, commandString string, env []string, delay int) *watchAndRun {
	w := NewWatcher(directory, match, exclude)
	r := NewRunner(commandString, env, delay)
	return &watchAndRun{w, r, false}
}

func (w *watchAndRun) WatchAndRun() error {
	w.watcher.SetVerboseLogging(w.Verbose)
	w.runner.SetVerboseLogging(w.Verbose)

	c, err := w.watcher.Watch()
	if err != nil {
		return err
	}

	w.runner.Run(c)

	return nil
}